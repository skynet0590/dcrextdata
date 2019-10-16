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

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "block propagation chart",
		Fetcher:  pg.fetchBlockPropagationChart,
		Appender: appendBlockPropagationChart,
	})

	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "PoW chart",
		Fetcher:  pg.fetchPowChart,
		Appender: appendPowChart,
	})
}
