package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger"
	"github.com/friendsofgo/errors"
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

func (charts ChartData) NormalizeLength() error {
	ids := []string{
		Mempool, Propagation, PowChart, VSP, Exchange, Snapshot, Community,
	}
	var err error
	for _, chartID := range ids {
		if cerr := charts.normalizeLength(chartID); cerr != nil {
			err = errors.Wrap(cerr, "Normalize failed for " + chartID)
		}
	}
	return err
}

// length correction
func (charts ChartData) normalizeLength(chartID string) error {
	// TODO: use transaction
	switch chartID {
	case Mempool:
		return charts.normalizeMempoolLength()
		
	case Propagation:
		return charts.normalizePropagationLength()
		
	case PowChart:
		return charts.normalizePowChartLength()
		
	case VSP:
		return charts.normalizeVSPLength()
		
	case Exchange:
		return charts.normalizeExchangeLength()
		
	case Snapshot:
		return charts.normalizeSnapshotLength()
	case Community:
		return nil

	}

	return nil
}

func (charts ChartData) normalizeMempoolLength() error {
	var firstLen, shortest, longest int
	key := Mempool + "-" + string(TimeAxis)
	firstLen, err := charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen

	key = Mempool + "-" + string(MempoolFees)
	dLen, err := charts.chartFloatsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizeMempoolLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}

	key = Mempool + "-" + string(MempoolSize)
	dLen, err = charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizeMempoolLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}

	key = Mempool + "-" + string(MempoolTxCount)
	dLen, err = charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizeMempoolLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}

	if longest != shortest {
		return charts.snipMempool(shortest)
	}
	return nil
}

func (charts ChartData) normalizePropagationLength() error {
	var firstLen, shortest, longest int
	key := Propagation + "-" + string(HeightAxis)
	firstLen, err := charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen

	key = Propagation + "-" + string(BlockTimestamp)
	dLen, err := charts.chartFloatsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizePropagationLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}

	key = Propagation + "-" + string(VotesReceiveTime)
	dLen, err = charts.chartFloatsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizePropagationLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}

	for _, source := range charts.syncSource {
		key = Propagation + "-" + string(BlockPropagation) + "-" + source
		dLen, err = charts.chartFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizePropagationLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}
	}

	if longest != shortest {
		return charts.snipPropagationChart(shortest)
	}
	return nil
}

func (charts ChartData) normalizePowChartLength() error {
	var firstLen, shortest, longest int
	key := PowChart + "-" + string(TimeAxis)
	firstLen, err := charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen

	for _, source := range charts.PowSources {
		key = PowChart + "-" + string(WorkerAxis) + "-" + source
		dLen, err := charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizePowChartLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		key = PowChart + "-" + string(HashrateAxis) + "-" + source
		dLen, err = charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizePowChartLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}
	}

	if longest != shortest {
		return charts.snipPowChart(shortest)
	}
	return nil
}

func (charts ChartData) normalizeVSPLength() error {
	var firstLen, shortest, longest int
	key := VSP + "-" + string(TimeAxis)
	firstLen, err := charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen

	for _, source := range charts.VSPSources {
		// ImmatureAxis
		key = VSP + "-" + string(ImmatureAxis) + "-" + source
		dLen, err := charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// LiveAxis
		key = VSP + "-" + string(LiveAxis) + "-" + source
		dLen, err = charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// VotedAxis
		key = VSP + "-" + string(VotedAxis) + "-" + source
		dLen, err = charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// MissedAxis
		key = VSP + "-" + string(MissedAxis) + "-" + source
		dLen, err = charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// PoolFeesAxis
		key = VSP + "-" + string(PoolFeesAxis) + "-" + source
		dLen, err = charts.chartNullFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// ProportionLiveAxis
		key = VSP + "-" + string(ProportionLiveAxis) + "-" + source
		dLen, err = charts.chartNullFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// ProportionMissedAxis
		key = VSP + "-" + string(ProportionMissedAxis) + "-" + source
		dLen, err = charts.chartNullFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// UsersActiveAxis
		key = VSP + "-" + string(UsersActiveAxis) + "-" + source
		dLen, err = charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// UserCountAxis
		key = VSP + "-" + string(UserCountAxis) + "-" + source
		dLen, err = charts.chartNullUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeVSPLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

	}

	if longest != shortest {
		return charts.snipPowChart(shortest)
	}
	return nil
}

func (charts ChartData) normalizeExchangeLength() error {

	var shortest, longest int
	for _, exchangeKeys := range charts.ExchangeKeys {
		key := exchangeKeys + "-" + string(TimeAxis)
		firstLen, err := charts.chartUintsLength(key)
		if err != nil {
			return err
		}
		shortest, longest = firstLen, firstLen

		// ExchangeOpenAxis
		key = exchangeKeys + "-" + string(ExchangeOpenAxis)
		dLen, err := charts.chartFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeExchangeLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// ExchangeCloseAxis
		key = exchangeKeys + "-" + string(ExchangeCloseAxis)
		dLen, err = charts.chartFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeExchangeLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// ExchangeHighAxis
		key = exchangeKeys + "-" + string(ExchangeHighAxis)
		dLen, err = charts.chartFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeExchangeLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}

		// ExchangeLowAxis
		key = exchangeKeys + "-" + string(ExchangeLowAxis)
		dLen, err = charts.chartFloatsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeExchangeLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}
	}

	if longest != shortest {
		return charts.snipPowChart(shortest)
	}
	return nil
}

func (charts ChartData) normalizeSnapshotLength() error {
	var firstLen, shortest, longest int
	key := Snapshot + "-" + string(TimeAxis)
	firstLen, err := charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen

	// SnapshotNodes
	key = Snapshot + "-" + string(SnapshotNodes)
	dLen, err := charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizeSnapshotLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}

	// SnapshotReachableNodes
	key = Snapshot + "-" + string(SnapshotReachableNodes)
	dLen, err = charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	if dLen != firstLen {
		log.Warnf("charts.normalizeSnapshotLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
		if dLen < shortest {
			shortest = dLen
		} else if dLen > longest {
			longest = dLen
		}
	}
	if longest != shortest {
		if err = charts.snipSnapshotChart(shortest, SnapshotNodes); err != nil {
			log.Warnf("SnapshotNodeVersions fail at %s, %s", SnapshotNodes, err.Error())
		}
	}

	// SnapshotLocations
	key = Snapshot + "-" + string(SnapshotLocations) + "-" + string(TimeAxis)
	firstLen, err = charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen
	for _, source := range charts.NodeLocations {
		key = Snapshot + "-" + string(SnapshotLocations) + "-" + source
		dLen, err := charts.chartUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeSnapshotLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}
	}
	if longest != shortest {
		if err = charts.snipSnapshotChart(shortest, SnapshotLocations); err != nil {
			log.Warnf("SnapshotNodeVersions fail at %s, %s", SnapshotLocations, err.Error())
		}
	}

	key = Snapshot + "-" + string(SnapshotNodeVersions) + "-" + string(TimeAxis)
	firstLen, err = charts.chartUintsLength(key)
	if err != nil {
		return err
	}
	shortest, longest = firstLen, firstLen
	for _, source := range charts.NodeVersion {
		// SnapshotNodeVersions
		key = Snapshot + "-" + string(SnapshotNodeVersions) + "-" + source
		dLen, err := charts.chartUintsLength(key)
		if err != nil {
			return err
		}
		if dLen != firstLen {
			log.Warnf("charts.normalizeSnapshotLength: dataset for %s axis has mismatched length %d != %d", key, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}
	}

	if longest != shortest {
		if err = charts.snipSnapshotChart(shortest, SnapshotNodeVersions); err != nil {
			log.Warnf("SnapshotNodeVersions fail at %s, %s", SnapshotNodeVersions, err.Error())
		}
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

func (charts ChartData) snipSnapshotChart(length int, axis axisType) error {
	var keys []string
	switch axis {
	case SnapshotNodes:
		keys = []string{
			Snapshot + "-" + string(TimeAxis),
			Snapshot + "-" + string(SnapshotNodes),
			Snapshot + "-" + string(SnapshotReachableNodes),
		}
		break
	case SnapshotLocations:
		keys = append(keys, Snapshot + "-" + string(SnapshotLocations) + "-" + string(TimeAxis))
		for _, country := range charts.NodeLocations {
			keys = append(keys, Snapshot+"-"+string(SnapshotLocations)+"-"+country)
		}
		break
	case SnapshotNodeVersions:
		keys = append(keys, Snapshot + "-" + string(SnapshotNodeVersions) + "-" + string(TimeAxis))
		for _, userAgent := range charts.NodeVersion {
			keys = append(keys, Snapshot+"-"+string(SnapshotNodeVersions)+"-"+userAgent)
		}
		break
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
