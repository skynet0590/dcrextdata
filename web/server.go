package web

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"context"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrextdata/postgres/models"
)

type DataQuery interface {
	AllExchange(ctx context.Context) (models.ExchangeSlice, error)
	FetchExchangeTicks(ctx context.Context, name string, offset int, limit int) (models.ExchangeTickSlice, error)
	FetchVSPs(ctx context.Context) (models.VSPSlice, error) 
	VSPTicks(ctx context.Context, vspName string, offset int, limit int) (models.VSPTickSlice, error)
}

type Server struct {
	templates    map[string]*template.Template
	lock         sync.RWMutex
	db 		DataQuery
}

func StartHttpServer(httpHost, httpPort string, db DataQuery) {
	server := &Server{
		templates:    map[string]*template.Template{},
		db: db,
	}

	// load templates
	server.loadTemplates()

	router := chi.NewRouter()
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/public")
	FileServer(router, "/static", http.Dir(filesDir))
	server.registerHandlers(router)

	address := net.JoinHostPort(httpHost, httpPort)

	fmt.Printf("starting http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting web server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (s *Server) loadTemplates() {
	layout := "web/views/layout.html"
	tpls := map[string]string{
		"exchange.html": "web/views/exchange.html",
	}

	for i, v := range tpls {
		tpl, err := template.New(i).ParseFiles(v, layout)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		s.lock.Lock()
		s.templates[i] = tpl
		s.lock.Unlock()
	}
}

func (s *Server) render(tplName string, data map[string]interface{}, res http.ResponseWriter) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if tpl, ok := s.templates[tplName]; ok {
		err := tpl.Execute(res, data)
		if err != nil {
			log.Fatalf("error executing template: %s", err.Error())
		}
		return
	}

	log.Fatalf("template %s is not registered", tplName)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/", s.GetExchangeTicks)
}
