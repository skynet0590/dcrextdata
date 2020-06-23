package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	ExchangeCloseAxis axisType = "close"
	ExchangeHighAxis  axisType = "high"
	ExchangeOpenAxis  axisType = "open"
	ExchangeLowAxis   axisType = "low"
)

type exchangeSet struct {
	// holds a set of exchange tick where the key is exchange name, currency pair and interval joined by -
	Ticks map[string]exchangeTick
}

type exchangeTick struct {
	Time    ChartUints
	Open    ChartFloats
	Close   ChartFloats
	High    ChartFloats
	Low     ChartFloats
	cacheID uint64
}

func (tickSet *exchangeTick) Snip(length int) {
	tickSet.Time = tickSet.Time.snip(length)
	tickSet.Open = tickSet.Open.snip(length)
	tickSet.Close = tickSet.Close.snip(length)
	tickSet.High = tickSet.High.snip(length)
	tickSet.Low = tickSet.Low.snip(length)
}

func (set *exchangeSet) Append(charts *ChartData, key string, time ChartUints, open ChartFloats, close ChartFloats, high ChartFloats, low ChartFloats) {
	if len(time) == 0 {
		return
	}
	if err := charts.AppendChartUintsAxis(key+"-"+string(TimeAxis), time); err != nil {
		log.Errorf("Error in append exchange time, %s", err.Error())
	}
	if err := charts.AppendChartFloatsAxis(key+"-"+string(ExchangeOpenAxis), open); err != nil {
		log.Errorf("Error in append exchange open axis, %s", err.Error())
	}
	if err := charts.AppendChartFloatsAxis(key+"-"+string(ExchangeCloseAxis), close); err != nil {
		log.Errorf("Error in append exchange close axis, %s", err.Error())
	}
	if err := charts.AppendChartFloatsAxis(key+"-"+string(ExchangeHighAxis), high); err != nil {
		log.Errorf("Error in append exchange high axis, %s", err.Error())
	}
	if err := charts.AppendChartFloatsAxis(key+"-"+string(ExchangeLowAxis), low); err != nil {
		log.Errorf("Error in append exchange low axis, %s", err.Error())
	}
}

// BuildExchangeKey returns exchange name, currency pair and interval joined by -
func BuildExchangeKey(exchangeName string, currencyPair string, interval int) string {
	return fmt.Sprintf("%s-%s-%d", exchangeName, currencyPair, interval)
}

func ExtractExchangeKey(setKey string) (exchangeName string, currencyPair string, interval int) {
	keys := strings.Split(setKey, "-")
	if len(keys) > 0 {
		exchangeName = keys[0]
	}

	if len(keys) > 1 {
		currencyPair = keys[1]
	}

	if len(keys) > 2 {
		interval, _ = strconv.Atoi(keys[2])
	}
	return
}

func (charts *ChartData) ExchangeSetTime(key string) uint64 {
	var dates ChartUints
	if err := charts.ReadAxis(key + "-" + string(TimeAxis), &dates); err != nil {
		return 0
	}
	if len(dates) < 1 {
		return 0
	}
	return dates[len(dates)-1]
}

func makeExchangeChart(ctx context.Context, charts *ChartData, axis axisType, key ...string) ([]byte, error) {
	if len(key) < 1 {
		return nil, errors.New("exchange set key is required for exchange chart")
	}
	var dates ChartUints
	if err := charts.ReadAxis(key[0] + "-" + string(TimeAxis), &dates); err != nil {
		return nil, err
	}

	var yAxis ChartFloats
	if err := charts.ReadAxis(key[0] + "-" + string(axis), &yAxis); err != nil {
		log.Errorf("Cannot create exchange chart, %s", err.Error())
		return nil, errors.New("no record found for the selected exchange")
	}

	return charts.Encode(nil, dates, yAxis)

}
