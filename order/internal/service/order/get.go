package order

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) Get(ctx context.Context, uuid string) (model.OrderInfo, error) {
	if uuid == "" {
		logger.Info(ctx, "uuid is empty")
		return model.OrderInfo{}, model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			logger.Info(ctx, "order not found", zap.String("uuid", uuid))
			return model.OrderInfo{}, model.ErrOrderNotFound

		}

		logger.Info(ctx, "get order failed", zap.String("uuid", uuid), zap.Error(err))
		return model.OrderInfo{}, err
	}

	return order, nil
}
