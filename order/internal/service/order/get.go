package order

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.OrderInfo, error) {
	if uuid == "" {
		return model.OrderInfo{}, model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.OrderInfo{}, model.ErrOrderNotFound
		}

		return model.OrderInfo{}, err
	}

	return order, nil
}
