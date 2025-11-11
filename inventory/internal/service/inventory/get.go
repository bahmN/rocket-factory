package inventory

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	if uuid == "" {
		logger.Info(ctx, "no uuid provided")
		return model.Part{}, model.ErrEmptyUUID
	}

	part, err := s.repo.GetPart(ctx, uuid)
	if err != nil {
		logger.Info(ctx, "error in getting part", zap.Error(err))
		return model.Part{}, err
	}

	logger.Info(ctx, "got part", zap.Any("part", part))
	return part, nil
}
