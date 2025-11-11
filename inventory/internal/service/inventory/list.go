package inventory

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		logger.Info(ctx, "error in getting list part", zap.Error(err))
		return nil, err
	}

	logger.Info(ctx, "got list part", zap.Any("part", parts))
	return parts, nil
}
