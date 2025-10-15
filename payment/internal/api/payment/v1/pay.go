package v1

import (
	"context"

	"github.com/bahmN/rocket-factory/payment/internal/converter"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	uuid, err := a.paymentService.PayOrder(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: uuid,
	}, nil
}
