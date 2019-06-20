package postgres

import (
	"context"

	"github.com/raedahgroup/dcrextdata/mempool"
)

func (pg PgDb) StoreMempool(context.Context, mempool.Mempool) error  {
	panic("not implemented")
}

func (pg PgDb) LastMempool(ctx context.Context) (mempool.Mempool, error)  {
	panic("not implemented")
}