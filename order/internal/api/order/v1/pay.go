package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bahmN/rocket-factory/order/internal/model"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, string(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrEmptyUUID) || errors.Is(err, model.ErrOrderPaidOrCanceled) {
			return &orderV1.NotFoundError{
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

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
