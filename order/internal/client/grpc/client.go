package grpc

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}
