package v1

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/converter"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "order not found",
		}, nil
	}

	return converter.OrderToAPI(order), nil
}
