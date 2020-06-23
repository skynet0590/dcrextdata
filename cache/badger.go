package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger"
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

// length correction
func (charts ChartData) normalizeLength(chartID string) error {
	switch chartID {
	case Mempool:
		key := Mempool + "-" + string(TimeAxis)
		timeLen, err := charts.chartUintsLength(key)
		if err != nil {
			return err
		}
		timeLen++
	}

	return nil
}

func (charts ChartData) chartUintsLength(key string) (int, error) {
	var data ChartUints
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return 0, err
		}
	}
	return data.Length(), nil
}

func (charts ChartData) chartFloatsLength(key string) (int, error) {
	var data ChartFloats
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return 0, err
		}
	}
	return data.Length(), nil
}

func (charts ChartData) chartNullUintsLength(key string) (int, error) {
	var data chartNullIntsPointer
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return 0, err
		}
	}
	return data.Length(), nil
}

func (charts ChartData) chartNullFloatsLength(key string) (int, error) {
	var data chartNullFloatsPointer
	err := charts.ReadAxis(key, &data)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return 0, err
		}
	}
	return data.Length(), nil
}

// Snip
func (charts ChartData) Snip(chartID string, length int) error {
	switch chartID {
	case Mempool:
		return charts.snipMempool(length)
	case Propagation:
		return charts.snipPropagationChart(length)
	case PowChart:
		return charts.snipPowChart(length)
	case VSP:
		return charts.snipVspChart(length)
	case Exchange:
		return charts.snipExchangeChart(length)
	case Snapshot:
		return charts.snipSnapshotChart(length)
	case Community:
		return nil
	}
	return nil
}

func (charts ChartData) snipMempool(length int) error {
	axis := []axisType{
		TimeAxis, MempoolSize, MempoolTxCount,
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
}

func (charts ChartData) snipPropagationChart(length int) error {
	key := Propagation + "-" + string(HeightAxis)
	if err := charts.snipChartUintsAxis(key, length); err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	keys := []string{
		Propagation + "-" + string(BlockTimestamp),
		Propagation + "-" + string(VotesReceiveTime),
	}
	for _, source := range charts.syncSource {
		keys = append(keys, Propagation+"-"+string(BlockPropagation)+"-"+source)
	}
	for _, key := range keys {
		if err := charts.snipChartFloatsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
	}
	return nil
}

func (charts ChartData) snipPowChart(length int) error {
	key := PowChart + "-" + string(TimeAxis)
	if err := charts.snipChartUintsAxis(key, length); err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	var keys []string
	for _, pool := range charts.PowSources {
		keys = append(keys, PowChart+"-"+string(WorkerAxis)+"-"+pool)
		keys = append(keys, PowChart+"-"+string(HashrateAxis)+"-"+pool)
	}

	for _, key := range keys {
		if err := charts.snipChartNullUintsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
	}
	return nil
}

func (charts ChartData) snipVspChart(length int) error {
	key := VSP + "-" + string(TimeAxis)
	if err := charts.snipChartUintsAxis(key, length); err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	}
	var keys []string
	for _, vspSource := range charts.VSPSources {
		keys = append(keys,
			VSP+"-"+string(ImmatureAxis)+"-"+vspSource,
			VSP+"-"+string(LiveAxis)+"-"+vspSource,
			VSP+"-"+string(VotedAxis)+"-"+vspSource,
			VSP+"-"+string(MissedAxis)+"-"+vspSource,
			VSP+"-"+string(UsersActiveAxis)+"-"+vspSource,
			VSP+"-"+string(UserCountAxis)+"-"+vspSource,
		)
	}

	for _, key := range keys {
		if err := charts.snipChartNullUintsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
	}

	keys = []string{}
	for _, vspSource := range charts.VSPSources {
		keys = append(keys,
			VSP+"-"+string(PoolFeesAxis)+"-"+vspSource,
			VSP+"-"+string(ProportionLiveAxis)+"-"+vspSource,
			VSP+"-"+string(ProportionMissedAxis)+"-"+vspSource,
		)
	}

	for _, key := range keys {
		if err := charts.snipChartNullFloatsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
	}

	return nil
}

func (charts ChartData) snipExchangeChart(length int) error {
	for _, exchangeKey := range charts.ExchangeKeys {
		key := exchangeKey + "-" + string(TimeAxis)
		if err := charts.snipChartUintsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
		var keys = []string{
			key + "-" + string(ExchangeOpenAxis),
			key + "-" + string(ExchangeCloseAxis),
			key + "-" + string(ExchangeHighAxis),
			key + "-" + string(ExchangeLowAxis),
		}
		for _, key := range keys {
			if err := charts.snipChartFloatsAxis(key, length); err != nil {
				if err != badger.ErrKeyNotFound {
					return err
				}
			}
		}
	}
	return nil
}

func (charts ChartData) snipSnapshotChart(length int) error {
	keys := []string{
		Snapshot + "-" + string(TimeAxis),
		Snapshot + "-" + string(SnapshotNodes),
		Snapshot + "-" + string(SnapshotReachableNodes),
	}
	for _, country := range charts.NodeLocations {
		keys = append(keys, Snapshot+"-"+string(SnapshotLocations)+"-"+country)
	}
	for _, userAgent := range charts.NodeVersion {
		keys = append(keys, Snapshot+"-"+string(SnapshotNodeVersions)+"-"+userAgent)
	}
	for _, key := range keys {
		if err := charts.snipChartUintsAxis(key, length); err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
	}
	return nil
}

func (charts ChartData) snipCommunityChart(length int) error {
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

func (charts ChartData) MempoolTimeTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(Mempool+"-"+string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length()-1]
}

func (charts ChartData) PropagationHeightTip() uint64 {
	var heights ChartUints
	err := charts.ReadAxis(Propagation+"-"+string(HeightAxis), &heights)
	if err != nil {
		return 0
	}
	if len(heights) == 0 {
		return 0
	}
	return heights[heights.Length()-1]
}

func (charts ChartData) PowTimeTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(PowChart+"-"+string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length()-1]
}

func (charts ChartData) VSPTimeTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(VSP+"-"+string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length()-1]
}

func (charts ChartData) SnapshotTip() uint64 {
	var dates ChartUints
	err := charts.ReadAxis(Snapshot+"-"+string(TimeAxis), &dates)
	if err != nil {
		return 0
	}
	if len(dates) == 0 {
		return 0
	}
	return dates[dates.Length()-1]
}
