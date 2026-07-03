package app

import (
	"context"
	"database/sql"
	"errors"

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

	if s.repo == nil {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("warn", "Ordering Get Service Log", logCtx.Service, "/app/get_ordering", logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte("ordering repository is not configured"))
		}
		return domain.CheckOrderByProductIdResponse{
			ProductName: "repository not configured",
		}, nil
	}

	repo, err := s.repo.GetOrderByProductIdAct(ctx, request.ProductId)
	if err != nil {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("error", "Ordering Get Service Log", logCtx.Service, "/app/get_ordering", logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte(err.Error()))
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CheckOrderByProductIdResponse{}, shared.NewNotFoundError("order_not_found", "order not found")
		}
		return domain.CheckOrderByProductIdResponse{}, err
	}

	return domain.CheckOrderByProductIdResponse{
		ProductName: repo.OrderName,
	}, nil
}
