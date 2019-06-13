package web

import (
"fmt"
	"net/http"
)

func (s *Server) GetExchange(res http.ResponseWriter, req *http.Request) {
	

	data := map[string]interface{}{
		"exchange": 	 s.db.GetExchangeDataTick(),
	}


	s.render("exchange.html", data, res)
}

func (s *Server) GetSend(res http.ResponseWriter, req *http.Request) {

}

func (s *Server) PostSend(res http.ResponseWriter, req *http.Request) {

}

func (s *Server) GetReceive(res http.ResponseWriter, req *http.Request) {

}
