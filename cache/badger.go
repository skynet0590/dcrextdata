package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v2"
)

type Normalizer interface {
	Normalize() Lengther
}

func (charts ChartData) SaveAxis(data Normalizer, key string) error {
	rec := data.Normalize()
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