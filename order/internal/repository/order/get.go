package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.OrderInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.OrderInfo{}, model.ErrOrderNotFound
	}

	return converter.OrderToModel(*order), nil
}
