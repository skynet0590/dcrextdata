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

// /
func (s *Server) getExchangeTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()
	var allExhangeTicksSlice []ticks.TickDto
	
	// var err error
	allExhangeTicksSlice, err = s.db.AllExchangeTicks(ctx, offset, recordsPerPage)
	if err != nil {
		panic(err) // todo add appropraite error handler
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
		"currentPage":    pageToLoad,
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allExhangeTicksSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("exchange.html", data, res)
}

func (s *Server) GetFilteredExchangeTicks(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	page := req.FormValue("page")
	selectedFilter := req.FormValue("filter")

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()
	var allExhangeTicksSlice []ticks.TickDto
	// var err error
	if selectedFilter == "All" {
		allExhangeTicksSlice, err = s.db.AllExchangeTicks(ctx, offset, recordsPerPage)
		if err != nil {
			s.renderError(err.Error(), res)
			return
	} else {
		allExhangeTicksSlice, err = s.db.FetchExchangeTicks(ctx, selectedFilter, offset, recordsPerPage)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	}

	allExhangeSlice, err := s.db.AllExchange(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	totalCount, err := s.db.AllExchangeTicksCount(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data := map[string]interface{}{
		"exData":         allExhangeTicksSlice,
		"allExData":      allExhangeSlice,
		"selectedFilter": selectedFilter,
		"currentPage":    pageToLoad,
		"previousPage":   int(pageToLoad - 1),
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	defer RenderJSON(data, res)

	totalTxLoaded := int(offset) + len(allExhangeTicksSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}
}

// /vsps
func (s *Server) getVspTicks(res http.ResponseWriter, req *http.Request) {
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

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	var allVSPSlice []vsp.VSPTickDto
	if selectedFilter == "All" {
		allVSPSlice, err = s.db.AllVSPTicks(ctx, offset, recordsPerPage)
		if err != nil {
			s.renderError(err.Error(), res)
			return
		}
	} else {
		allVSPSlice, err = s.db.VSPTicks(ctx, selectedFilter, offset, recordsPerPage)
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

	totalCount, err := s.db.AllVSPTickCount(ctx)
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
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(allVSPSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	s.render("vsp.html", data, res)
}

// /pow
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

	totalCount, err := s.db.CountPowData(ctx)
	if err != nil {
		s.renderError(err.Error(), res)
		return
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

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	mempoolSlice, err := s.db.Mempools(ctx, offset, recordsPerPage)
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
		"previousPage": int(pageToLoad - 1),
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(mempoolSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	return data, nil
}

// /propagation
func (s *Server) propagation(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	block, err := s.fetchBlockData(req)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data["blocks"] = block

	votes, err := s.fetchVoteData(req)
	if err != nil {
		s.renderError(err.Error(), res)
		return
	}

	data["votes"] = votes

	s.render("propagation.html", data, res)
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

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	blockSlice, err := s.db.Blocks(ctx, offset, recordsPerPage)
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
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(blockSlice)
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

	pageToLoad, err := strconv.ParseInt(page, 10, 32)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (int(pageToLoad) - 1) * recordsPerPage

	ctx := context.Background()

	voteSlice, err := s.db.Votes(ctx, offset, recordsPerPage)
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
		"totalPages":   int(math.Ceil(float64(totalCount) / float64(recordsPerPage))),
	}

	totalTxLoaded := int(offset) + len(voteSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = int(pageToLoad + 1)
	}

	return data, nil
}
