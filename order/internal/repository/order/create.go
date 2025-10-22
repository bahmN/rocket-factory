package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, info model.OrderInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[info.OrderUUID] = converter.OrderInfoToRepo(info)

	return nil
}
