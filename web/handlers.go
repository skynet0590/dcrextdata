package web

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/raedahgroup/dcrextdata/exchanges/ticks"
	"github.com/raedahgroup/dcrextdata/vsp"
)

const (
	recordsPerPage = 20
)

func (s *Server) GetExchangeTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	filter := req.Form["exchange"]

	var selectedFilter string
	if len(filter) == 0 || filter[0] == "All" || filter[0] == "" {
		selectedFilter = "All"
	} else {
		selectedFilter = filter[0]
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	var recordsPerPage int = recordsPerPage
	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()
	var allExhangeTicksSlice []ticks.TickDto
	// var err error
	if selectedFilter == "All" {
		allExhangeTicksSlice, err = s.db.AllExchangeTicks(ctx, offset, recordsPerPage)
		if err != nil {
			panic(err)
		}
	} else {
		allExhangeTicksSlice, err = s.db.FetchExchangeTicks(ctx, selectedFilter, offset, recordsPerPage)
		if err != nil {
			panic(err)
		}
	}

	allExhangeSlice, err := s.db.AllExchange(ctx)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	totalCount, err := s.db.AllExchangeTicksCount(ctx)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"exData":         allExhangeTicksSlice,
		"allExData":      allExhangeSlice,
		"selectedFilter": selectedFilter,
		"currentPage":    int(pageToLoad),
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(txPerPage))),
	}

	totalTxLoaded := int(offset) + len(allExhangeTicksSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("exchange.html", data, res)
}

func (s *Server) GetVspTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	filter := req.Form["vsp"]

	var selectedFilter string
	if len(filter) == 0 || filter[0] == "All" || filter[0] == "" {
		selectedFilter = "All"
	} else {
		selectedFilter = filter[0]
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	var recordsPerPage int = recordsPerPage
	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	var allVSPSlice []vsp.VSPTickDto
	if selectedFilter == "All" {
		allVSPSlice, err = s.db.AllVSPTicks(ctx, offset, recordsPerPage)
		if err != nil {
			panic(err)
		}
	} else {
		allVSPSlice, err = s.db.VSPTicks(ctx, selectedFilter, offset, recordsPerPage)
		if err != nil {
			panic(err)
		}
	}

	allVspData, err := s.db.FetchVSPs(ctx)
	if err != nil {
		panic(err) // todo add appropraite error handler
	}

	totalCount, err := s.db.AllVSPTickCount(ctx)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"vspData":        allVSPSlice,
		"allVspData":     allVspData,
		"selectedFilter": selectedFilter,
		"currentPage":    int(pageToLoad),
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(txPerPage))),
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
