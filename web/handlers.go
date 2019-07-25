package web

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/exchanges/ticks"
	"github.com/raedahgroup/dcrextdata/vsp"
)

var (
	recordsPerPage  = 20
	defaultInterval = 5
)

// /home
func (s *Server) homePage(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	s.render("home.html", data, res)
}

// /exchange
func (s *Server) getExchangeTicks(res http.ResponseWriter, req *http.Request) {
	pageToLoad := 1
	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	allExhangeTicksSlice, totalCount, err := s.db.AllExchangeTicks(ctx, "", offset, defaultInterval, recordsPerPage)
	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	allExhangeSlice, err := s.db.AllExchange(ctx)
	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	currencyPairs, err := s.db.AllExchangeTicksCurrencyPair(ctx)
	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	intervals, err := s.db.AllExchangeTicksInterval(ctx)
	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"exData":         allExhangeTicksSlice,
		"allExData":      allExhangeSlice,
		"currencyPairs":  currencyPairs,
		"intervals":      intervals,
		"currentPage":    pageToLoad,
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
		"selectedFilter": "All",
		"selectedCpair":  "All",
		"selectedNum":    recordsPerPage,
	}

	totalTxLoaded := int(offset) + len(allExhangeTicksSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("exchange.html", data, res)
}

func (s *Server) getFilteredExchangeTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedFilter := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")
	selectedCurrencyPair := req.FormValue("selectedCurrencyPair")
	interval := req.FormValue("selectedInterval")

	data := map[string]interface{}{}

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}

	filterInterval, err := strconv.Atoi(interval)
	if err != nil || filterInterval <= 0 {
		filterInterval = defaultInterval
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	var allExhangeTicksSlice []ticks.TickDto
	var totalCount int64

	if selectedFilter == "All" && selectedCurrencyPair == "All" {
		allExhangeTicksSlice, totalCount, err = s.db.AllExchangeTicks(ctx, "", filterInterval, offset, pageSize)
		if err != nil {
			s.renderErrorJSON(err.Error(), res)
			return
		}

	} else if selectedFilter == "All" && selectedCurrencyPair != "All" {
		allExhangeTicksSlice, totalCount, err = s.db.AllExchangeTicks(ctx, selectedCurrencyPair, filterInterval, offset, pageSize)
		if err != nil {
			s.renderErrorJSON(err.Error(), res)
			return
		}
	} else if selectedFilter != "All" && selectedCurrencyPair == "All" {
		allExhangeTicksSlice, totalCount, err = s.db.FetchExchangeTicks(ctx, "", selectedFilter, filterInterval, offset, pageSize)
		if err != nil {
			s.renderErrorJSON(err.Error(), res)
			return
		}
	} else {
		allExhangeTicksSlice, totalCount, err = s.db.FetchExchangeTicks(ctx, selectedCurrencyPair, selectedFilter, filterInterval, offset, pageSize)
		if err != nil {
			s.renderErrorJSON(err.Error(), res)
			return
		}
	}

	var intervalRange string
	if len(allExhangeTicksSlice) == 0 {
		if filterInterval == 5 {
			intervalRange = "5m"
		}
		if filterInterval == 60 {
			intervalRange = "1h"
		}
		if filterInterval == 120 {
			intervalRange = "2h"
		}
		if filterInterval == 1440 {
			intervalRange = "1d"
		}

		data["message"] = fmt.Sprintf("%s does not have %s data.", strings.Title(selectedFilter), intervalRange)
		s.renderJSON(data, res)
		return
	}

	allExhangeSlice, err := s.db.AllExchange(ctx)
	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	data["exData"] = allExhangeTicksSlice
	data["allExData"] = allExhangeSlice
	data["selectedFilter"] = selectedFilter
	data["currentPage"] = pageToLoad
	data["previousPage"] = int(pageToLoad - 1)
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	defer s.renderJSON(data, res)

	totalTxLoaded := int(offset) + len(allExhangeTicksSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}
}

func (s *Server) getChartData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	selectedTick := req.FormValue("selectedTick")
	selectedCurrencyPair := req.FormValue("selectedCurrencyPair")
	selectedInterval := req.FormValue("selectedInterval")
	sources := req.FormValue("sources")

	data := map[string]interface{}{}

	ctx := context.Background()
	interval, err := strconv.Atoi(selectedInterval)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	chartData, err := s.db.ExchangeTicksChartData(ctx, selectedTick, selectedCurrencyPair, interval, sources)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}
	if len(chartData) == 0 {
		data["message"] = fmt.Sprintf("No data to generate %s chart.", sources)
		s.renderJSON(data, res)
		return
	}

	data["chartData"] = chartData

	defer s.renderJSON(data, res)
}

// /vsps
func (s *Server) getVspTicks(res http.ResponseWriter, req *http.Request) {
	pageToLoad := 1

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	allVSPSlice, err := s.db.AllVSPTicks(ctx, offset, recordsPerPage)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	allVspData, err := s.db.FetchVSPs(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	totalCount, err := s.db.AllVSPTickCount(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"vspData":        allVSPSlice,
		"allVspData":     allVspData,
		"selectedFilter": "All",
		"currentPage":    pageToLoad,
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allVSPSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("vsp.html", data, res)
}

func (s *Server) getFilteredVspTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedFilter := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	var allVSPSlice []vsp.VSPTickDto
	var totalCount int64
	if selectedFilter == "All" || selectedFilter == "" {
		allVSPSlice, err = s.db.AllVSPTicks(ctx, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}

		totalCount, err = s.db.AllVSPTickCount(ctx)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	} else {
		allVSPSlice, err = s.db.FiltredVSPTicks(ctx, selectedFilter, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}

		totalCount, err = s.db.FiltredVSPTicksCount(ctx, selectedFilter)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	}

	allVspData, err := s.db.FetchVSPs(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"vspData":        allVSPSlice,
		"allVspData":     allVspData,
		"selectedFilter": selectedFilter,
		"currentPage":    pageToLoad,
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	defer s.renderJSON(data, res)

	totalTxLoaded := int(offset) + len(allVSPSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}
}

// vspchartdata
func (s *Server) vspChartData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sources := req.FormValue("sources")
	selectedAttribute := req.FormValue("selectedAttribute")

	vsps := strings.Split(sources, "|")

	ctx := context.Background()
	dates, err := s.db.GetVspTickDistinctDates(ctx, vsps)
	if err != nil {
		s.renderErrorJSON(fmt.Sprintf("Error is getting dates from VSP table, %s", err.Error()), res)
		return
	}

	var resultMap = map[time.Time][]interface{}{}
	for _, date := range dates {
		resultMap[date] = []interface{}{date}
	}

	for _, source := range vsps {
		points, err := s.db.FetchChartData(ctx, selectedAttribute, source)
		if err != nil {
			s.renderErrorJSON(fmt.Sprintf("Error in fetching %s records for %s: %s", selectedAttribute, source, err.Error()), res)
			return
		}

		var vspPointMap = map[time.Time]interface{}{}
		for _, point := range points {
			vspPointMap[point.Date] = point.Record
		}

		//var previousDate time.Time
		for date, _ := range resultMap {
			if record, found := vspPointMap[date]; found {
				resultMap[date] = append(resultMap[date], record)
				//previousDate = date
			} else {
				resultMap[date] = append(resultMap[date], 0)
				/*if record, found := vspPointMap[previousDate]; found {
					resultMap[date] = append(resultMap[date], record)
				} else {
					fmt.Println(fmt.Sprintf("not found and %s not seen before", previousDate.String()))
					resultMap[date] = append(resultMap[date], 0)
				}*/
			}
		}
	}

	/*
		var chartData [][]interface{}
		for _, points := range resultMap {
			chartData = append(chartData, points)
		}*/

	s.renderJSON(resultMap, res)
}

// /PoW
func (s *Server) getPowData(res http.ResponseWriter, req *http.Request) {
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
		s.renderError(err.Error(), res)
		return
	}

	powSource, err := s.db.FetchPowSourceData(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	totalCount, err := s.db.CountPowData(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"powData":      allPowDataSlice,
		"powSource":    powSource,
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

func (s *Server) getFilteredPowData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedFilter := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")

	data := map[string]interface{}{}
	defer s.renderJSON(data, res)

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	var totalCount int64
	var totalTxLoaded int
	if selectedFilter == "All" || selectedFilter == "" {
		allPowDataSlice, err := s.db.FetchPowData(ctx, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}

		data["powData"] = allPowDataSlice

		totalCount, err = s.db.CountPowData(ctx)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}

		totalTxLoaded = int(offset) + len(allPowDataSlice)
	} else {
		allPowDataSlice, err := s.db.FetchPowDataBySource(ctx, selectedFilter, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}

		data["powData"] = allPowDataSlice

		totalCount, err = s.db.CountPowDataBySource(ctx, selectedFilter)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}

		totalTxLoaded = int(offset) + len(allPowDataSlice)
	}

	powSource, err := s.db.FetchPowSourceData(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data["powSource"] = powSource
	data["currentPage"] = int(pageToLoad)
	data["previousPage"] = int(pageToLoad - 1)
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}
}

// /mempool
func (s *Server) mempoolPage(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	mempoolData, err := s.fetchMempoolData(req)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data["mempool"] = mempoolData

	s.render("mempool.html", data, res)
}

// /getmempool
func (s *Server) getMempool(res http.ResponseWriter, req *http.Request) {
	data, err := s.fetchMempoolData(req)
	defer s.renderJSON(data, res)

	if err != nil {
		data = map[string]interface{}{
			"error": err.Error(),
		}
		return
	}
}

func (s *Server) fetchMempoolData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("recordsPerPage")

	var pageSize = recordsPerPage

	if numberOfRows != "" {
		numRows, err := strconv.Atoi(numberOfRows)
		if err != nil || numRows <= 0 {
			pageSize = pageSize
		} else {
			pageSize = numRows
		}
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (pageToLoad - 1) * pageSize

	ctx := context.Background()

	mempoolSlice, err := s.db.Mempools(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.db.MempoolCount(ctx)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"mempoolData":  mempoolSlice,
		"currentPage":  pageToLoad,
		"previousPage": pageToLoad - 1,
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	totalTxLoaded := int(offset) + len(mempoolSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	return data, nil
}

func (s *Server) getMempoolChartData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	chartFilter := req.FormValue("chartFilter")
	ctx := context.Background()

	mempoolDataSlice, err := s.db.MempoolsChartData(ctx, chartFilter)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"mempoolchartData": mempoolDataSlice,
		"chartFilter":      chartFilter,
	}

	defer s.renderJSON(data, res)
}

// /propagation
func (s *Server) propagation(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	block, err := s.fetchPropagationData(req)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data["propagation"] = block

	s.render("propagation.html", data, res)
}

// /getPropagationData
func (s *Server) getPropagationData(res http.ResponseWriter, req *http.Request) {
	data, err := s.fetchPropagationData(req)
	defer s.renderJSON(data, res)

	if err != nil {
		data = map[string]interface{}{
			"error": err.Error(),
		}
		return
	}
}

func (s *Server) fetchPropagationData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("recordsPerPage")

	var pageSize = recordsPerPage

	if numberOfRows != "" {
		numRows, err := strconv.Atoi(numberOfRows)
		if err != nil || numRows <= 0 {
			pageSize = pageSize
		} else {
			pageSize = numRows
		}
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	blockSlice, err := s.db.Blocks(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.db.BlockCount(ctx)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"records":      blockSlice,
		"currentPage":  pageToLoad,
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	totalTxLoaded := int(offset) + len(blockSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	return data, nil
}

// /getblocks
func (s *Server) getBlocks(res http.ResponseWriter, req *http.Request) {
	data, err := s.fetchBlockData(req)
	defer s.renderJSON(data, res)

	if err != nil {
		data = map[string]interface{}{
			"error": err.Error(),
		}
		return
	}
}

func (s *Server) fetchBlockData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("recordsPerPage")

	var pageSize = recordsPerPage

	if numberOfRows != "" {
		numRows, err := strconv.Atoi(numberOfRows)
		if err != nil || numRows <= 0 {
			pageSize = pageSize
		} else {
			pageSize = numRows
		}
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	voteSlice, err := s.db.BlocksWithoutVotes(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.db.BlockCount(ctx)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"records":      voteSlice,
		"currentPage":  pageToLoad,
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	totalTxLoaded := int(offset) + len(voteSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	return data, nil
}

// /getvotes
func (s *Server) getVotes(res http.ResponseWriter, req *http.Request) {
	data, err := s.fetchVoteData(req)
	defer s.renderJSON(data, res)

	if err != nil {
		data = map[string]interface{}{
			"error": err.Error(),
		}
		return
	}
}

func (s *Server) fetchVoteData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("recordsPerPage")

	var pageSize = recordsPerPage

	if numberOfRows != "" {
		numRows, err := strconv.Atoi(numberOfRows)
		if err != nil || numRows <= 0 {
			pageSize = pageSize
		} else {
			pageSize = numRows
		}
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	voteSlice, err := s.db.Votes(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.db.VotesCount(ctx)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"records":      voteSlice,
		"currentPage":  pageToLoad,
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	totalTxLoaded := int(offset) + len(voteSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	return data, nil
}
