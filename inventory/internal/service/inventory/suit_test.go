package inventory

import (
	"testing"

	"github.com/bahmN/rocket-factory/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuit struct {
	suite.Suite
	inventoryRepository *mocks.InventoryRepository
	service             *service
}

func (s *ServiceSuit) SetupTest() {
	s.inventoryRepository = &mocks.InventoryRepository{}
	s.service = NewService(s.inventoryRepository)
}

func (s *ServiceSuit) TearDownTest() {}

func TestNewService(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
