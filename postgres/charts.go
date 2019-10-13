package postgres

import "github.com/raedahgroup/dcrextdata/cache"

func (pg *PgDb) RegisterCharts(charts *cache.ChartData) {
	charts.AddUpdater(cache.ChartUpdater{
		Tag:      "mempool chart",
		Fetcher:  pg.chartMempool,
		Appender: appendChartMempool,
	})
}
