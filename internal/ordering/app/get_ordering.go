package app

import (
	"context"

	"firstbeegoapi/internal/ordering/domain"
	"firstbeegoapi/internal/shared"
)

func (s *OrderingService) CheckOrderByProductIdAct(ctx context.Context, request domain.CheckOrderByProductIdRequest) (domain.CheckOrderByProductIdResponse, error) {
	if request.ProductId <= 0 {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("warn", "Ordering Get Service Log", logCtx.Service, "/app/get_ordering", logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte("product id must be greater than zero"))
		}
		return domain.CheckOrderByProductIdResponse{}, shared.NewValidationError("invalid_product_id", "product id must be greater than zero")
	}

	return domain.CheckOrderByProductIdResponse{
		ProductId: request.ProductId,
	}, nil
}
