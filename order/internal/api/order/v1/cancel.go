package v1

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	res, err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrEmptyUUID) || errors.Is(err, model.ErrOrderPaidOrCanceled) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "order not found",
			}, nil
		}
	}

	return &orderV1.CancelOrderResponse{
		Message: res,
	}, nil
}
