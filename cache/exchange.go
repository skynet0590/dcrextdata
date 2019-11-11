package cache

import (
	"errors"
	"fmt"
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

func (set *exchangeSet) Snip(length int) {
	for _, tickSet := range set.Ticks {
		tickSet.Snip(length)
	}
}

func newExchangeSet() *exchangeSet {
	return &exchangeSet{Ticks: map[string]exchangeTick{}}
}

func (set *exchangeSet) Append(key string, time ChartUints, open ChartFloats, close ChartFloats, high ChartFloats, low ChartFloats) {
	if existingTick, found := set.Ticks[key]; found {

		existingTick.Time = append(existingTick.Time, time...)
		existingTick.Open = append(existingTick.Open, open...)
		existingTick.Close = append(existingTick.Close, close...)
		existingTick.High = append(existingTick.High, high...)
		existingTick.Low = append(existingTick.Low, low...)
	} else {
		set.Ticks[key] = exchangeTick{
			Time:  time,
			Open:  open,
			Close: close,
			High:  high,
			Low:   low,
		}
	}
}

// BuildExchangeKey returns exchange name, currency pair and interval joined by -
func BuildExchangeKey(exchangeName string, currencyPair string, interval int) string {
	return fmt.Sprintf("%s-%s-%d", exchangeName, currencyPair, interval)
}

func (charts *ChartData) ExchangeSetTime(key string) uint64 {
	if tick, found := charts.Exchange.Ticks[key]; found && len(tick.Time) > 0 {
		return tick.Time[len(tick.Time)-1]
	}
	return 0
}

func makeExchangeChart(charts *ChartData, axis axisType, setKey ...string) ([]byte, error) {
	if len(setKey) < 1 {
		return nil, errors.New("exchange set key is required for exchange chart")
	}

	if tick, found := charts.Exchange.Ticks[setKey[0]]; found {
		var yAxis ChartFloats
		switch axis {
		case ExchangeOpenAxis:
			yAxis = tick.Open
			break
		case ExchangeCloseAxis:
			yAxis = tick.Close
			break
		case ExchangeLowAxis:
			yAxis = tick.Low
			break
		case ExchangeHighAxis:
			yAxis = tick.High
			break
		default:
			return nil, errors.New("invalid exchange chart axis")
		}

		return charts.encode(nil, tick.Time, yAxis)
	}
	return nil, errors.New("no record found for the selected exchange")
}
