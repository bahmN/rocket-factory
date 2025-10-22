package inventory

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	if uuid == "" {
		return model.Part{}, model.ErrEmptyUUID
	}

	part, err := s.repo.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
