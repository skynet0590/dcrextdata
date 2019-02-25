package main

import (
	"net/http"
	"time"
)

type ExchangeCollector struct {
	exchanges []Exchange
	period    int64
}

func NewExchangeCollector(exchangeLasts map[string]int64, period int64) (*ExchangeCollector, error) {
	exchanges := make([]Exchange, 0, len(exchangeLasts))

	for exchange, last := range exchangeLasts {
		if contructor, ok := ExchangeConstructors[exchange]; ok {
			ex, err := contructor(&http.Client{Timeout: 300 * time.Second}, last, period) // Consider if sharing a single client is better
			if err != nil {
				return nil, err
			}
			exchanges = append(exchanges, ex)
		}
	}

	return &ExchangeCollector{
		exchanges: exchanges,
		period:    period,
	}, nil
}

func (ec *ExchangeCollector) HistoricSyncRequired() bool {
	now := time.Now().Unix()
	for _, ex := range ec.exchanges {
		if now-ex.LastUpdateTime() > ec.period {
			return true
		}
	}
	return false
}

func (ec *ExchangeCollector) HistoricSync(data chan []DataTick) bool {
	errs := make(chan error)
	for _, ex := range ec.exchanges {
		go func(ex Exchange, dataChan chan []DataTick) {
			errs <- ex.Historic(dataChan)
		}(ex, data)
	}

	failed := false
	for err := range errs {
		if err != nil {
			failed = true
			excLog.Error(err)
		}
	}
	return !failed
}

func (ec *ExchangeCollector) Collect(data chan []DataTick, quit chan struct{}) {
	ticker := time.NewTicker(time.Duration(ec.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			excLog.Trace("Triggering exchange collectors")
			for _, ex := range ec.exchanges {
				go ex.Collect(data)
			}
		case <-quit:
			return
		}

	}
}
