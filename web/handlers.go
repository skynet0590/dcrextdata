package web

import (
	"net/http"
)

func (s *Server) GetExchange(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{
		// "exchange": 	 s.db.GetExchangeDataTick(),
	}

	s.render("exchange.html", data, res)
}

