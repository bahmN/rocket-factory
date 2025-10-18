package v1

import (
	def "github.com/bahmN/rocket-factory/order/internal/client/grpc"
	generatedPaymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient generatedPaymentV1.PaymentServiceClient
}

func NewClient(generatedClient generatedPaymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
