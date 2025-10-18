package part

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	repoConverter "github.com/bahmN/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.data[uuid]
	if !exists {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(part), nil
}
