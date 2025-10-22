package v1

import (
	"context"
	"log"

	"github.com/bahmN/rocket-factory/order/internal/converter"
	"github.com/bahmN/rocket-factory/order/internal/model"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/go-faster/errors"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrEmptyUUID) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.BadRequestError{
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

	return converter.OrderToAPI(order), nil
}
