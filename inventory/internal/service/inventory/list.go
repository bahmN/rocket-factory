package inventory

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
