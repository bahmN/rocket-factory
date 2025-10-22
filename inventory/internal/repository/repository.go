package repository

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}
