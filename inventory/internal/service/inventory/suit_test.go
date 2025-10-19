package inventory

import (
	"context"
	"testing"

	"github.com/bahmN/rocket-factory/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuit struct {
	suite.Suite
	ctx                 context.Context
	inventoryRepository *mocks.InventoryRepository
	service             *service
}

func (s *ServiceSuit) SetupTest() {
	s.ctx = context.Background()
	s.inventoryRepository = &mocks.InventoryRepository{}
	s.service = NewService(s.inventoryRepository)
}

func (s *ServiceSuit) TearDownTest() {}

func TestNewService(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
