package api

import (
	"strconv"

	"firstbeegoapi/internal/ordering/domain"
	"firstbeegoapi/internal/shared"
)

func validateCheckOrderByProductIDRequest(objectID string) (domain.CheckOrderByProductIdRequest, error) {
	if objectID == "" {
		return domain.CheckOrderByProductIdRequest{}, shared.NewValidationError("missing_object_id", "object id is required")
	}

	id, err := strconv.Atoi(objectID)
	if err != nil {
		return domain.CheckOrderByProductIdRequest{}, shared.NewValidationError("invalid_object_id", "object id must be a number")
	}
	if id <= 0 {
		return domain.CheckOrderByProductIdRequest{}, shared.NewValidationError("invalid_object_id", "object id must be greater than zero")
	}

	return domain.CheckOrderByProductIdRequest{
		ProductId: int64(id),
	}, nil
}
