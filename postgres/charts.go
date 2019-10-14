package postgres

import "github.com/raedahgroup/dcrextdata/cache"

func (pg *PgDb) RegisterCharts(charts *cache.ChartData) {
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
}
