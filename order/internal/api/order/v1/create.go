package v1

import (
	"context"
	"fmt"

	"github.com/bahmN/rocket-factory/order/internal/converter"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	orderResponse, err := a.orderService.Create(ctx, converter.CreateOrderToModel(req))
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("internal service error: %v", err),
		}, nil
	}

	return converter.CreateOrderToAPI(&orderResponse), nil
}
