package v1

import (
	"github.com/bahmN/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewApi(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
