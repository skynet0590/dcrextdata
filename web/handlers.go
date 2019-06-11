package web

import (
	"net/http"
)

func (s *Server) GetBalance(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	s.render("balance.html", data, res)
}

func (s *Server) GetSend(res http.ResponseWriter, req *http.Request) {

}

func (s *Server) PostSend(res http.ResponseWriter, req *http.Request) {

}

func (s *Server) GetReceive(res http.ResponseWriter, req *http.Request) {

}
