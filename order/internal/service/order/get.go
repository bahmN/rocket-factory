package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.OrderInfo, error) {
	if uuid == "" {
		return model.OrderInfo{}, model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return model.OrderInfo{}, err
	}

	return order, nil
}
