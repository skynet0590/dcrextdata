package web

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrextdata/exchanges/ticks"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/raedahgroup/dcrextdata/vsp"
	"github.com/raedahgroup/dcrextdata/pow"
)

type DataQuery interface {
	AllExchangeTicks(ctx context.Context, offset int, limit int) ([]ticks.TickDto, error)
	AllExchangeTicksCount(ctx context.Context) (int64, error)
	AllExchange(ctx context.Context) (models.ExchangeSlice, error)	
	FetchExchangeTicks(ctx context.Context, name string, offset int, limit int) ([]ticks.TickDto, error)
	
	FetchVSPs(ctx context.Context) (models.VSPSlice, error)
	VSPTicks(ctx context.Context, vspName string, offset int, limit int) ([]vsp.VSPTickDto, error)
	AllVSPTicks(ctx context.Context, offset int, limit int) ([]vsp.VSPTickDto, error)
	AllVSPTickCount(ctx context.Context) (int64, error)

	FetchPowData(ctx context.Context, offset int, limit int) ([]pow.PowDataDto, error)
	CountPowData(ctx context.Context) (int64, error)
	FetchPowDataBySource(ctx context.Context, source string, offset int, limit int) ([]pow.PowDataDto, error)
	CountPowDataBySource(ctx context.Context, source string) (int64, error)
}

type Server struct {
	templates map[string]*template.Template
	lock      sync.RWMutex
	db        DataQuery
}

func StartHttpServer(httpHost, httpPort string, db DataQuery) {
	server := &Server{
		templates: map[string]*template.Template{},
		db:        db,
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
		"vsp.html": "web/views/vsp.html",
		"pow.html": "web/views/pow.html",
	}

	for i, v := range tpls {
		tpl, err := template.New(i).Funcs(templateFuncMap()).ParseFiles(v, layout)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		s.lock.Lock()
		s.templates[i] = tpl
		s.lock.Unlock()
	}
}

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"incByOne": func(number int) int {
			return number + 1
		},
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
	r.Get("/vspticks", s.GetVspTicks)
	r.Get("/pow", s.GetPowData)
}
