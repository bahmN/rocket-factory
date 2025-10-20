package order

import (
	"testing"

	clientMocks "github.com/bahmN/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/bahmN/rocket-factory/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuit struct {
	suite.Suite
	orderRepository *mocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient
	service         *service
}

func (s *ServiceSuit) SetupTest() {
	s.orderRepository = &mocks.OrderRepository{}
	s.inventoryClient = &clientMocks.InventoryClient{}
	s.paymentClient = &clientMocks.PaymentClient{}

	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func (s *ServiceSuit) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
