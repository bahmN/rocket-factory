package service

import (
	"context"

	"github.com/bahmN/rocket-factory/assembly/internal/model"
)

type OrderProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
