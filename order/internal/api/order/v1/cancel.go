package v1

import (
	"context"

	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	res, err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "order not found",
		}, nil
	}

	return &orderV1.CancelOrderResponse{
		Message: res,
	}, nil
}
