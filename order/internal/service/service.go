package service

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req model.CreateOrderReq) (model.CreateOrderResp, error)
	Get(ctx context.Context, uuid string) (model.OrderInfo, error)
	Pay(ctx context.Context, uuid, method string) (string, error)
	Cancel(ctx context.Context, uuid string) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
