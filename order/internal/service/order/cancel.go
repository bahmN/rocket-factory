package order

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	if uuid == "" {
		logger.Info(ctx, "no uuid provided")
		return model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			logger.Info(ctx, "order not found", zap.String("uuid", uuid))
			return model.ErrOrderNotFound
		}

		logger.Info(ctx, "error in getting order", zap.Error(err))
		return err
	}

	if order.Status == model.OrderStatusPAID || order.Status == model.OrderStatusCANCELLED {
		logger.Info(ctx, "order already cancelled", zap.String("uuid", uuid), zap.String("status", order.Status))
		return model.ErrOrderPaidOrCanceled
	}

	order.Status = model.OrderStatusCANCELLED
	err = s.repo.Update(ctx, order.OrderUUID, order)
	if err != nil {
		logger.Warn(ctx, "error in updating order", zap.Error(err))
		return err
	}

	logger.Info(ctx, "order cancelled", zap.String("uuid", uuid))
	return nil
}
