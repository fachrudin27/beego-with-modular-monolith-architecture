package repository

import (
	"context"
	"database/sql"
	"firstbeegoapi/internal/ordering/domain"
	"firstbeegoapi/internal/shared"
)

func (o *OrderingRepository) GetOrderByProductIdAct(ctx context.Context, payload int64) (domain.Order, error) {
	if o.db == nil {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("error", "Ordering Get Service Log", logCtx.Service, logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte("*sql.DB is nil"))
		}
		return domain.Order{}, sql.ErrConnDone
	}

	const query = `
	SELECT order_name
	FROM orders
	LIMIT 1
`

	var orderName string

	err := o.db.QueryRowContext(ctx, query).Scan(&orderName)
	if err != nil {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("error", "Ordering Get Service Log", logCtx.Service, logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte(err.Error()))
		}
		return domain.Order{}, err
	}

	return domain.Order{
		OrderName: orderName,
	}, nil
}
