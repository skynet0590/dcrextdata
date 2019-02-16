package db

import (
	"strings"

	"github.com/raedahgroup/dcrextdata/collection/exchanges"
	"github.com/raedahgroup/dcrextdata/db/internal"
)

func (pg *PgDb) AddExchangeData(data []exchanges.DataTick) (int, error) {
	added := 0
	for _, v := range data {
		_, err := pg.db.Exec(internal.InsertExchangeDataTick, v.High, v.Low, v.Open, v.Close, v.Time, v.Exchange)
		if err != nil {
			if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
				return added, err
			}
		}
		added++
	}
	return added, nil
}

func (pg *PgDb) LastExchangeEntryTime(exchange string) (time int64) {
	rows := pg.db.QueryRow(internal.LastExchangeEntryTime, exchange)
	_ = rows.Scan(&time)

	// if err != nil {
	// 	if strings.Contains(err.Error(), "no rows") {

	// 	} else {
	// 		log.Error("Could not retrieve last entry time: ", err)
	// 		return err
	// 	}
	// }

	return
}
