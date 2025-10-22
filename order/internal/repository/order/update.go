package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/samber/lo"
)

func (r *repository) Update(ctx context.Context, uuid string, info model.OrderInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := r.data[uuid]

	if lo.IsNil(order) {
		return model.ErrOrderNotFound
	}

	if lo.IsNotEmpty(info.OrderUUID) {
		order.OrderUUID = info.OrderUUID
	}

	if lo.IsNotEmpty(info.UserUUID) {
		order.UserUUID = info.UserUUID
	}

	if lo.IsNotNil(info.PartUUIDs) {
		order.PartUUIDs = info.PartUUIDs
	}

	if info.TotalPrice > 0 {
		order.TotalPrice = info.TotalPrice
	}

	if lo.IsNotNil(info.TransactionUUID) {
		order.TransactionUUID = info.TransactionUUID
	}

	if lo.IsNotNil(info.PaymentMethod) {
		order.PaymentMethod = info.PaymentMethod
	}

	if lo.IsNotEmpty(info.Status) {
		order.Status = info.Status
	}

	r.data[uuid] = order

	return nil
}
