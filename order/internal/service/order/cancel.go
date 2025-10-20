package order

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) (string, error) {
	if uuid == "" {
		return "", model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}

		return "", err
	}

	if order.Status == model.OrderStatusPAID || order.Status == model.OrderStatusCANCELLED {
		return "", model.ErrOrderPaidOrCanceled
	}

	order.Status = model.OrderStatusCANCELLED
	err = s.repo.Update(ctx, order.OrderUUID, order)
	if err != nil {
		return "", err
	}

	return "order cancelled", nil
}
