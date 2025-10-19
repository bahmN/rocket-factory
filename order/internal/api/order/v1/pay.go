package v1

import (
	"context"
	"fmt"

	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, string(req.PaymentMethod))
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "order not found",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
