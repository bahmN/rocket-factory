package order

import (
	"github.com/bahmN/rocket-factory/order/internal/client/grpc"
	"github.com/bahmN/rocket-factory/order/internal/repository"
	def "github.com/bahmN/rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	repo repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderPaidProducer def.OrderProducerService
}

func NewService(
	repo repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	orderPaidProducer def.OrderProducerService,
) *service {
	return &service{
		repo:              repo,
		inventoryClient:   inventoryClient,
		paymentClient:     paymentClient,
		orderPaidProducer: orderPaidProducer,
	}
}
