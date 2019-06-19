package web

import (
	"context"
	"net/http"
)

func (s *Server) GetExchangeTicks(res http.ResponseWriter, req *http.Request) {
	allExhangeSlice, err := s.db.AllExchangeTicks(context.Background(), 0, 30)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"exData": allExhangeSlice,
	}

	s.render("exchange.html", data, res)
}

func (s *Server) GetVspTicks(res http.ResponseWriter, req *http.Request) {
	allVSPSlice, err := s.db.AllVSPTicks(context.Background(), 0, 30)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"vspData": allVSPSlice,
	}

	s.render("vsp.html", data, res)
}
