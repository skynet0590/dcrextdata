package postgres

import (
	"strings"
	
	"github.com/raedahgroup/dcrextdata/pow"
)

func (pg *PgDb) LastPowEntryTime(source string) (time int64) {
	rows := pg.db.QueryRow(LastPowEntryTime, source)
	_ = rows.Scan(&time)
	return
}

//
func (pg *PgDb) AddPowData(data []pow.PowData) error {
	added := 0
	for _, d := range data {
		_, err := pg.db.Exec(InsertPowData, d.Time, d.NetworkHashrate, d.PoolHashrate, d.Workers, d.NetworkDifficulty, d.CoinPrice, d.BtcPrice, d.Source)
		if err != nil {
			if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
				return err
			}
		}
		added++
	}
	if len(data) == 1 {
		log.Infof("Added %d entry from %s (%s)", added, data[0].Source,
			UnixTimeToString(data[0].Time))
	} else {
		last := data[len(data)-1]
		log.Infof("Added %d entries from %s (%s to %s)", added, last.Source,
			UnixTimeToString(data[0].Time), UnixTimeToString(last.Time))
	}

	return nil
}
