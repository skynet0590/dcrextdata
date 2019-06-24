package web

import (
	"context"
	"math"
	"net/http"
	"strconv"
)

const (
	recordsPerPage = 20
)

func (s *Server) GetExchangeTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	var recordsPerPage int = recordsPerPage
	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()
	allExhangeSlice, err := s.db.AllExchangeTicks(ctx, offset, recordsPerPage)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	totalCount, err := s.db.AllExchangeTicksCount(ctx)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	data := map[string]interface{}{
		"exData":       allExhangeSlice,
		"currentPage":  int(pageToLoad),
		"previousPage": int(pageToLoad - 1),
<<<<<<< HEAD
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
=======
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(txPerPage))),
>>>>>>> mempool
	}

	totalTxLoaded := int(offset) + len(allExhangeSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("exchange.html", data, res)
}

func (s *Server) GetVspTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	var recordsPerPage int = recordsPerPage
	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	allVSPSlice, err := s.db.AllVSPTicks(ctx, offset, recordsPerPage)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	totalCount, err := s.db.AllVSPTickCount(ctx)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"vspData":      allVSPSlice,
		"currentPage":  int(pageToLoad),
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allVSPSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("vsp.html", data, res)
}

func (s *Server) GetPowData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	allPowDataSlice, err := s.db.FetchPowData(ctx, offset, recordsPerPage)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	totalCount, err := s.db.CountPowData(ctx)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	data := map[string]interface{}{
		"powData":      allPowDataSlice,
		"currentPage":  int(pageToLoad),
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allPowDataSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("pow.html", data, res)
}
