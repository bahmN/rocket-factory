package v1

import (
	"context"
	"fmt"
	"log"

	"github.com/bahmN/rocket-factory/order/internal/converter"
	"github.com/bahmN/rocket-factory/order/internal/model"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/go-faster/errors"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	order, err := converter.CreateOrderToModel(req)
	if err != nil {
		log.Printf("error converting data: %v", err)

		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	orderResponse, err := a.orderService.Create(ctx, order)
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
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

	return converter.CreateOrderToAPI(&orderResponse), nil
}
