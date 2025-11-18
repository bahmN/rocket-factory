package service

import (
	"context"

	"github.com/bahmN/rocket-factory/notification/internal/model"
)

type TelegramService interface {
	SendOrderPaidNotify(ctx context.Context, msg model.OrderPaidEvent) error
	SendOrderAssemblyNotify(ctx context.Context, msg model.OrderAssembledEvent) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}
