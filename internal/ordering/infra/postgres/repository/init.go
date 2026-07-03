package repository

import "database/sql"

// type OrderingRepository interface {
// 	CheckOrderByProductIdAct(ctx context.Context, productId int) (domain.CheckOrderByProductIdResponse, error)
// }

type OrderingRepository struct {
	db *sql.DB
}

func NewOrderingRepository(db *sql.DB) *OrderingRepository {
	return &OrderingRepository{db: db}
}
