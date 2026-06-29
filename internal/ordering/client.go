package ordering

import (
	"context"

	"firstbeegoapi/internal/ordering/app"
	"firstbeegoapi/internal/ordering/domain"
)

type OrderingContract interface {
	CheckOrderByProductIdAct(ctx context.Context, productId int) (domain.CheckOrderByProductIdResponse, error)
}

type orderingApi struct {
	app.OrderingService
}

func NewModuleApi() OrderingContract {
	return &orderingApi{}
}

func (o *orderingApi) CheckOrderByProductIdAct(ctx context.Context, productId int) (domain.CheckOrderByProductIdResponse, error) {
	res, err := o.OrderingService.CheckOrderByProductIdAct(ctx, domain.CheckOrderByProductIdRequest{
		ProductId: int64(productId),
	})
	if err != nil {
		return domain.CheckOrderByProductIdResponse{}, err
	}
	return res, nil
}
