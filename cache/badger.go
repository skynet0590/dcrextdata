package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v2"
)

func (charts ChartData) SaveAxis(rec Lengther, key string) error {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(rec); err != nil {
		return err
	}
	err := charts.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), b.Bytes())
		return err
	})
	return err
}

func (charts ChartData) ClearVLog() {
	again:
		verr := charts.db.RunValueLogGC(0.7)
		if verr == nil {
			goto again
		}
}

func (charts ChartData) ReadAxis(key string, result Lengther) error {
	return charts.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			d := gob.NewDecoder(bytes.NewReader(val))
			if err := d.Decode(result); err != nil {
				return err
			}
			return nil
		})
	})
}

// Appenders

func (charts ChartData) AppendChartUintsAxis(key string, set ChartUints) error {
	var data ChartUints
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = append(data, set...)
	return charts.SaveAxis(data, key)
}

func (charts ChartData) AppendChartNullUintsAxis(key string, set ChartNullUints) error {
	var data chartNullIntsPointer
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = data.Append(set)
	return charts.SaveAxis(data, key)
}

func (charts ChartData) AppendChartFloatsAxis(key string, set ChartFloats) error {
	var data ChartFloats
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = append(data, set...)
	return charts.SaveAxis(data, key)
}

func (charts ChartData) AppendChartNullFloatsAxis(key string, set ChartNullFloats) error {
	var data chartNullFloatsPointer
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = data.Append(set)
	return charts.SaveAxis(data, key)
}

// Snip
func (charts ChartData) Snip(chartID string, length int, axis ...axisType) error {
	switch chartID {
	case Mempool:
		if len(axis) == 0 {
			axis = append(axis, TimeAxis, MempoolSize, MempoolTxCount)
		}
		for _, a := range axis {
			key := Mempool + "-" + string(a)
			if err := charts.snipChartUintsAxis(key, length); err != nil {
				if err != badger.ErrKeyNotFound {
					return err
				}
			}
		}
		key := Mempool + "-" + string(MempoolFees)
		if err := charts.snipChartFloatsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
		return nil
	case Propagation:
		if len(axis) == 0 {
			axis = append(axis, TimeAxis, MempoolFees, MempoolSize, MempoolTxCount)
		}
		for _, a := range axis {
			key := Mempool + "-" + string(a)
			if err := charts.snipChartUintsAxis(key, length); err != nil {
				if err != badger.ErrKeyNotFound {
					return err
				}
			}
		}
		return nil
	}
	return nil
}

func (charts ChartData) snipChartUintsAxis(key string, length int) error {
	var data ChartUints
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = data.snip(length)
	return charts.SaveAxis(data, key)
}

func (charts ChartData) snipChartNullUintsAxis(key string, length int) error {
	var data chartNullIntsPointer
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = data.snip(length)
	return charts.SaveAxis(data, key)
}

func (charts ChartData) snipChartNullFloatsAxis(key string, length int) error {
	var data chartNullFloatsPointer
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = data.snip(length)
	return charts.SaveAxis(data, key)
}

func (charts ChartData) snipChartFloatsAxis(key string, length int) error {
	var data ChartFloats
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	data = data.snip(length)
	return charts.SaveAxis(data, key)
}

// Length

func (charts ChartData) MempoolTimeTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(Mempool + "-" + string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length() - 1]
}

func (charts ChartData) PropagationHeightTip() uint64 {
	var heights ChartUints
	err := charts.ReadAxis(Propagation + "-" + string(HeightAxis), &heights)
	if err != nil {
		return 0
	}
	if len(heights) == 0 {
		return 0
	}
	return heights[heights.Length() - 1]
}

func (charts ChartData) PowTimeTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(PowChart + "-" + string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length() - 1]
}

func (charts ChartData) VSPTimeTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(VSP + "-" + string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length() - 1]
}

func (charts ChartData) SnapshotTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(Snapshot + "-" + string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length() - 1]
}
