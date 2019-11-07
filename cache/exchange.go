package cache

import "errors"

const(
	ExchangeCloseAxis axisType = "close"
	ExchangeHighAxis axisType = "high"
	ExchangeOpenAxis axisType = "open"
	ExchangeLowAxis axisType = "low"
)

type exchangeAxis struct {

}

type exchangeSet struct {
	exchanges map[string]exchangeTickSet
}

func (set *exchangeSet) Snip(length int) {
	for _, tickSet := range set.exchanges {
		tickSet.Snip(length)
	}
}

// exchangeTickSet is a set of exchange tick data
type exchangeTickSet struct {
	Time ChartUints
	Ticks ChartNullFloats
}

// Snip truncates the exchangeSet to a provided length.
func (set *exchangeTickSet) Snip(length int) {
	set.Time = set.Time.snip(length)
	set.Ticks = set.Ticks.snip(length)
}

func (exch *exchangeSet) ticks(exchange string, axis axisType, interval int) (time ChartUints, ticks ChartFloats, err error) {
	if tickSet, found := exch.exchanges[exchange]; found {
		return tickSet.Time, nil, nil
	} else {
		return nil, nil, errors.New("exchange not found")
	}

	return nil, nil, nil
}