package inventory

import (
	"github.com/bahmN/rocket-factory/inventory/internal/repository"
	def "github.com/bahmN/rocket-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	repo repository.InventoryRepository
}

func NewService(repo repository.InventoryRepository) *service {
	return &service{
		repo: repo,
	}
}
