package v1

import (
	"context"

	clientConverter "github.com/bahmN/rocket-factory/order/internal/client/converter"
	generatedPaymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	payment, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: clientConverter.PaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		return "", err
	}

	return payment.TransactionUuid, nil
}
