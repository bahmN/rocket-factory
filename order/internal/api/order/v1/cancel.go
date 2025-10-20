package v1

import (
	"context"
	"errors"
	"log"

	"github.com/bahmN/rocket-factory/order/internal/model"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
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
				Message: err.Error(),
			}, nil
		}

		log.Printf("unknown error: %v", err)

		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	return &orderV1.CancelOrderResponse{
		Message: "order cancelled",
	}, nil
}
