package postgres

import "github.com/raedahgroup/dcrextdata/cache"

func (pg *PgDb) RegisterCharts(charts *cache.ChartData, syncSources []string, syncSourceDbProvider func(source string) (*PgDb, error)) {
	pg.syncSourceDbProvider = syncSourceDbProvider
	pg.syncSources = syncSources

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "mempool chart",
		Fetcher:  pg.retrieveChartMempool,
		Appender: appendChartMempool,
	})
	charts.AddRetriever(cache.Mempool, pg.fetchEncodeMempoolChart)

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "block propagation chart",
		Fetcher:  pg.fetchBlockPropagationChart,
		Appender: appendBlockPropagationChart,
	})
	charts.AddRetriever(cache.Propagation, pg.fetchEncodePropagationChart)

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "PoW chart",
		Fetcher:  pg.fetchCachePowChart,
		Appender: appendPowChart,
	})
	charts.AddRetriever(cache.PowChart, pg.fetchEncodePowChart)

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "VSP chart",
		Fetcher:  pg.fetchCacheVspChart,
		Appender: appendVspChart,
	})
	charts.AddRetriever(cache.VSP, pg.fetchEncodeVspChart)

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "Exchange chart",
		Fetcher:  pg.fetchExchangeChart,
		Appender: appendExchangeChart,
	})
	charts.AddRetriever(cache.Exchange, pg.fetchEncodeExchangeChart)

	charts.AddUpdater(cache.ChartUpdater{
		Tag: "Snapshot chart",
		Fetcher: pg.fetchNetworkSnapshotChart,
		Appender: appendSnapshotChart,
	})
	charts.AddRetriever(cache.Snapshot, pg.fetchEncodeSnapshotChart)
}