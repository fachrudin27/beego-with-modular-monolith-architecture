package app

import (
	"context"
	"firstbeegoapi/internal/ordering/domain"
)

type OrderingRepository interface {
	GetOrderByProductIdAct(ctx context.Context, payload int64) (domain.Order, error)
}
