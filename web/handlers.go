package web

import (
	"net/http"
	"context"
)

func (s *Server) GetExchangeTicks(res http.ResponseWriter, req *http.Request) {
	exhangeSlice, err := s.db.FetchExchangeTicks(context.Background(), "bittrex", 0, 30)
	if err != nil {
		panic(err)
	}
	
	data := map[string]interface{}{
		"exData" : exhangeSlice,
	}

	s.render("exchange.html", data, res)
}

