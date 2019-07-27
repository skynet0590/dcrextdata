package web

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/exchanges/ticks"
	"github.com/raedahgroup/dcrextdata/mempool"
	"github.com/raedahgroup/dcrextdata/vsp"
	"github.com/raedahgroup/dcrextdata/pow"
)

var (
	recordsPerPage        = 20
	defaultInterval       = -1 // All
	exchangeTickIntervals = map[int]string{
		-1:   "All",
		5:    "5m",
		60:   "1h",
		120:  "2h",
		1440: "1d",
	}

	pageSizeSelector = map[int]int{
		20:   20,
		30:    30,
		50:   50,
		100:  100,
		150: 150,
	}
)

// /home
func (s *Server) homePage(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	s.render("home.html", data, res)
}

// /exchange
func (s *Server) getExchangeTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedExchange := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")
	selectedCurrencyPair := req.FormValue("selectedCurrencyPair")
	interval := req.FormValue("selectedInterval")
	fmt.Println(page, selectedExchange,numberOfRows,selectedCurrencyPair,interval)

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
	}

	filterInterval, err := strconv.Atoi(interval)
	if err != nil || filterInterval <= 0 {
		filterInterval = defaultInterval
	}

	if _, found := exchangeTickIntervals[filterInterval]; !found {
		filterInterval = defaultInterval
	}

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	if selectedExchange == "" {
		selectedExchange = "All"
	}

	if selectedCurrencyPair == "" {
		selectedCurrencyPair = "All"
	}

	offset := (int(pageToLoad) - 1) * pageSize

	ctx := context.Background()

	allExhangeTicksSlice, totalCount, err := s.db.FetchExchangeTicks(ctx, selectedCurrencyPair, selectedExchange, filterInterval, offset, pageSize)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	if len(allExhangeTicksSlice) == 0 {
		message := map[string]interface{}{
			"message": fmt.Sprintf("%s does not have %s data.", strings.Title(selectedExchange), exchangeTickIntervals[filterInterval]),
		}
		s.render("exchange.html", message, res)
		return
	}

	allExhangeSlice, err := s.db.AllExchange(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	currencyPairs, err := s.db.AllExchangeTicksCurrencyPair(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"exData":         allExhangeTicksSlice,
		"allExData":      allExhangeSlice,
		"currencyPairs":  currencyPairs,
		"intervals":      exchangeTickIntervals,
		"pageSizeSelector": pageSizeSelector,
		"currentPage":    pageToLoad,
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(pageSize))),
		"selectedFilter": selectedExchange,
		"selectedCurrencyPair":  selectedCurrencyPair,
		"selectedNum":    pageSize,
		"selectedInterval": filterInterval,
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
	selectedExchange := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")
	selectedCurrencyPair := req.FormValue("selectedCurrencyPair")
	interval := req.FormValue("selectedInterval")

	fmt.Println(page, selectedExchange,numberOfRows,selectedCurrencyPair,interval)
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

	if _, found := exchangeTickIntervals[filterInterval]; !found {
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

	allExhangeTicksSlice, totalCount, err = s.db.FetchExchangeTicks(ctx, selectedCurrencyPair, selectedExchange, filterInterval, offset, pageSize)
	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	if len(allExhangeTicksSlice) == 0 {
		data["message"] = fmt.Sprintf("%s does not have %s data.", strings.Title(selectedExchange), exchangeTickIntervals[filterInterval])
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
	data["selectedFilter"] = selectedExchange
	data["currentPage"] = pageToLoad
	data["previousPage"] = int(pageToLoad - 1)
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
		data["selectedCurrencyPair"]=  selectedCurrencyPair
		data["selectedNum"]=    pageSize
		data["selectedInterval"] = filterInterval

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
	req.ParseForm()
	page := req.FormValue("page")
	selectedVsp := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")

	fmt.Println(page, selectedVsp,numberOfRows)

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}
	
	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	if selectedVsp == "" {
		selectedVsp = "All"
	} 

	offset := (pageToLoad - 1) * pageSize

	ctx := context.Background()

	var allVSPSlice []vsp.VSPTickDto
	var totalCount int64
	if selectedVsp == "All" || selectedVsp == "" {
		allVSPSlice, totalCount, err = s.db.AllVSPTicks(ctx, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	} else {
		allVSPSlice, totalCount, err = s.db.FiltredVSPTicks(ctx, selectedVsp, offset, pageSize)
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
		"pageSizeSelector": pageSizeSelector,
		"selectedFilter": selectedVsp,
		"currentPage":    pageToLoad,
		"selectedNum": pageSize,
		"previousPage":   pageToLoad - 1,
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allVSPSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	s.render("vsp.html", data, res)
}

func (s *Server) getFilteredVspTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedVsp := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")

	fmt.Println(page, selectedVsp,numberOfRows)

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	if selectedVsp == "" {
		selectedVsp = "All"
	} 

	offset := (pageToLoad - 1) * pageSize

	ctx := context.Background()

	var allVSPSlice []vsp.VSPTickDto
	var totalCount int64
	if selectedVsp == "All" || selectedVsp == "" {
		allVSPSlice, totalCount, err = s.db.AllVSPTicks(ctx, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	} else {
		allVSPSlice, totalCount, err = s.db.FiltredVSPTicks(ctx, selectedVsp, offset, pageSize)
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
		"selectedFilter": selectedVsp,
		"selectedNum": numberOfRows,
		"currentPage":    pageToLoad,
		"previousPage":   pageToLoad - 1,
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	defer s.renderJSON(data, res)

	totalTxLoaded := int(offset) + len(allVSPSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}
}

// vspchartdata
func (s *Server) vspChartData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sources := req.FormValue("vsps")
	selectedAttribute := req.FormValue("selectedAttribute")

	vsps := strings.Split(sources, "|")

	ctx := req.Context()
	dates, err := s.db.GetVspTickDistinctDates(ctx, vsps)
	if err != nil {
		s.renderErrorJSON(fmt.Sprintf("Error is getting dates from VSP table, %s", err.Error()), res)
		return
	}

	var vspChartData = struct {
		CSV     string    `json:"csv"`
		MinDate time.Time `json:"min_date"`
		MaxDate time.Time `json:"max_date"`
	}{
		CSV: "Date," + strings.Join(vsps, ",") + "\n",
	}

	var resultMap = map[time.Time][]string{}
	for _, date := range dates {
		if vspChartData.MinDate.IsZero() || date.Before(vspChartData.MinDate) {
			vspChartData.MinDate = date
		}
		if vspChartData.MaxDate.IsZero() || date.After(vspChartData.MaxDate) {
			vspChartData.MaxDate = date
		}
		resultMap[date] = []string{date.String()}
	}

	for _, source := range vsps {
		points, err := s.db.FetchChartData(ctx, selectedAttribute, source)
		if err != nil {
			s.renderErrorJSON(fmt.Sprintf("Error in fetching %s records for %s: %s", selectedAttribute, source, err.Error()), res)
			return
		}

		var vspPointMap = map[time.Time]string{}
		var vspDates []time.Time
		for _, point := range points {
			vspPointMap[point.Date] = point.Record
			vspDates = append(vspDates, point.Date)
		}

		sort.Slice(vspDates, func(i, j int) bool {
			return vspDates[i].Before(vspDates[j])
		})

		for date, _ := range resultMap {
			if date.Year() == 1970 || date.IsZero() {
				continue
			}
			if record, found := vspPointMap[date]; found {
				skip := false
				if record == "0" || record == "" {
					skip = true
					for _, vspDate := range vspDates {
						if vspDate.Before(date) && vspPointMap[vspDate] != "" && vspPointMap[vspDate] != "0" {
							skip = false
						}
					}
				}
				if !skip {
					resultMap[date] = append(resultMap[date], record)
				} else {
					resultMap[date] = append(resultMap[date], "Nan")
				}
			} else {
				// if they have not been any record for this vsp, give a gap (Nan) else use space
				padding := "Nan"
				for _, vspDate := range vspDates {
					if vspDate.Before(date) && vspPointMap[vspDate] != "" && vspPointMap[vspDate] != "0" {
						padding = ""
					}
				}
				resultMap[date] = append(resultMap[date], padding)
			}
		}
	}

	for _, date := range dates {
		if date.Year() == 1970 || date.IsZero() {
			continue
		}

		points := resultMap[date]
		hasAtleastOneRecord := false
		for index, point := range points {
			// the first index is the date
			if index == 0 {
				continue
			}
			if point != "" && point != "Nan" {
				hasAtleastOneRecord = true
			}
		}

		if !hasAtleastOneRecord {
			continue
		}

		vspChartData.CSV += fmt.Sprintf("%s\n", strings.Join(points, ","))
	}

	s.renderJSON(vspChartData, res)
}

// /PoW
func (s *Server) getPowData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedPow := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")

	fmt.Println(page, selectedPow,numberOfRows)

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}
	
	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	if selectedPow == "" {
		selectedPow = "All"
	} 

	offset := (pageToLoad - 1) * recordsPerPage

	ctx := context.Background()

	var totalCount int64
	var allPowDataSlice []pow.PowDataDto
	if selectedPow == "All" || selectedPow == "" {
		allPowDataSlice, totalCount, err = s.db.FetchPowData(ctx, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	} else {
		allPowDataSlice, totalCount, err = s.db.FetchPowDataBySource(ctx, selectedPow, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	}

	powSource, err := s.db.FetchPowSourceData(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"powData":      allPowDataSlice,
		"selectedFilter": selectedPow,
		"pageSizeSelector": pageSizeSelector,
		"selectedNum": pageSize,
		"powSource":    powSource,
		"currentPage":  pageToLoad,
		"previousPage": pageToLoad - 1,
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allPowDataSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	s.render("pow.html", data, res)
}

func (s *Server) getFilteredPowData(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedPow := req.FormValue("filter")
	numberOfRows := req.FormValue("recordsPerPage")

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = recordsPerPage
	} else {
		pageSize = numRows
	}

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	if selectedPow == "" {
		selectedPow = "All"
	}

	offset := (pageToLoad - 1) * pageSize

	ctx := context.Background()

	var totalCount int64
	var allPowDataSlice []pow.PowDataDto
	if selectedPow == "All" || selectedPow == "" {
		allPowDataSlice, totalCount, err = s.db.FetchPowData(ctx, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	} else {
		allPowDataSlice, totalCount, err = s.db.FetchPowDataBySource(ctx, selectedPow, offset, pageSize)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	}

	powSource, err := s.db.FetchPowSourceData(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"powData":      allPowDataSlice,
		"selectedFilter": selectedPow,
		"selectedNum": pageSize,
		"powSource":    powSource,
		"currentPage":  pageToLoad,
		"previousPage": pageToLoad - 1,
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	defer s.renderJSON(data, res)

	totalTxLoaded := int(offset) + len(allPowDataSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}
}

func (s *Server) getPowChartDate(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sources := req.FormValue("pools")
	dataType := req.FormValue("datatype")

	pools := strings.Split(sources, "|")

	ctx := req.Context()
	dates, err := s.db.GetPowDistinctDates(ctx, pools)
	if err != nil {
		s.renderErrorJSON(fmt.Sprintf("Error is getting dates from PoW table, %s", err.Error()), res)
		return
	}

	var powChartData = struct {
		CSV     string    `json:"csv"`
		MinDate time.Time `json:"min_date"`
		MaxDate time.Time `json:"max_date"`
	}{
		CSV: "Date," + strings.Join(pools, ",") + "\n",
	}

	var resultMap = map[time.Time][]string{}
	for _, date := range dates {
		if powChartData.MinDate.IsZero() || date.Before(powChartData.MinDate) {
			powChartData.MinDate = date
		}
		if powChartData.MaxDate.IsZero() || date.After(powChartData.MaxDate) {
			powChartData.MaxDate = date
		}
		resultMap[date] = []string{date.String()}
	}

	for _, source := range pools {
		points, err := s.db.FetchPowChartData(ctx, source, dataType)
		if err != nil {
			s.renderErrorJSON(fmt.Sprintf("Error in fetching %s records for %s: %s", dataType, source, err.Error()), res)
			return
		}

		var pointMaps = map[time.Time]string{}
		var powDates []time.Time
		for _, point := range points {
			pointMaps[point.Date] = point.Record
			powDates = append(powDates, point.Date)
		}

		sort.Slice(powDates, func(i, j int) bool {
			return powDates[i].Before(powDates[j])
		})

		for date, _ := range resultMap {
			if date.Year() == 1970 || date.IsZero() {
				continue
			}
			if record, found := pointMaps[date]; found {
				skip := false
				if record == "0" || record == "" {
					skip = true
					for _, powDate := range powDates {
						if powDate.Before(date) && pointMaps[powDate] != "" && pointMaps[powDate] != "0" {
							skip = false
						}
					}
				}
				if !skip {
					resultMap[date] = append(resultMap[date], record)
				} else {
					resultMap[date] = append(resultMap[date], "Nan")
				}
			} else {
				// if they have not been any record for this vsp, give a gap (Nan) else use space
				padding := "Nan"
				for _, powDate := range powDates {
					if powDate.Before(date) && pointMaps[powDate] != "" && pointMaps[powDate] != "0" {
						padding = ""
					}
				}
				resultMap[date] = append(resultMap[date], padding)
			}
		}
	}

	for _, date := range dates {
		if date.Year() == 1970 || date.IsZero() {
			continue
		}

		points := resultMap[date]
		hasAtleastOneRecord := false
		for index, point := range points {
			// the first index is the date
			if index == 0 {
				continue
			}
			if point != "" && point != "Nan" {
				hasAtleastOneRecord = true
			}
		}

		if !hasAtleastOneRecord {
			continue
		}

		powChartData.CSV += fmt.Sprintf("%s\n", strings.Join(points, ","))
	}

	s.renderJSON(powChartData, res)

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

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
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
		"pageSizeSelector": pageSizeSelector,
		"selectedNum": pageSize,
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

// /propagationchartdata
func (s *Server) propagationChartData(res http.ResponseWriter, req *http.Request) {
	requestedRecordSet := req.FormValue("recordset")
	var data []mempool.PropagationChartData
	var err error

	if requestedRecordSet == "votes" {
		data, err = s.db.PropagationVoteChartData(req.Context())
	} else {
		data, err = s.db.PropagationBlockChartData(req.Context())
	}

	if err != nil {
		s.renderErrorJSON(err.Error(), res)
		return
	}

	var avgTimeForHeight = map[int64]float64{}
	var heightArr []int64
	for _, record := range data {
		if existingTime, found := avgTimeForHeight[record.BlockHeight]; found {
			avgTimeForHeight[record.BlockHeight] = (record.TimeDifference + existingTime)/2
		} else {
			avgTimeForHeight[record.BlockHeight] = record.TimeDifference
			heightArr = append(heightArr, record.BlockHeight)
		}
	}

	var yLabel = "Delay (s)"
	if requestedRecordSet == "votes" {
		yLabel = "Time Difference (s)"
	}

	var csv = fmt.Sprintf("Height,%s\n", yLabel)
	for _, height := range heightArr {
		timeDifference := fmt.Sprintf("%04.2f", avgTimeForHeight[height])
		csv += fmt.Sprintf("%d, %s\n", height, timeDifference)
	}

	s.renderJSON(csv, res)
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

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
	}

	offset := (pageToLoad - 1) * pageSize

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
		"pageSizeSelector": pageSizeSelector,
		"selectedNum": pageSize,
		"url": "/propagation",
		"previousPage": pageToLoad - 1,
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	totalTxLoaded := int(offset) + len(blockSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
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

func (s *Server) getBlockData(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	block, err := s.fetchBlockData(req)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data["propagation"] = block
	defer s.render("propagation.html", data, res)
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

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
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
		"pageSizeSelector": pageSizeSelector,
		"selectedNum": pageSize,
		"url": "/getBlockData",
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

	if _, foundPageSize := pageSizeSelector[pageSize]; !foundPageSize {
		pageSize = recordsPerPage
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
		"pageSizeSelector": pageSizeSelector,
		"selectedNum": pageSize,
		"url": "/getvotes",
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	totalTxLoaded := int(offset) + len(voteSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	return data, nil
}
