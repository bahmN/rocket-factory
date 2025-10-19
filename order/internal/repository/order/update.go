package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

func (r *repository) Update(ctx context.Context, uuid string, info model.OrderInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := r.data[uuid]

	if info.OrderUUID != "" {
		order.OrderUUID = info.OrderUUID
	}

	if info.UserUUID != "" {
		order.UserUUID = info.UserUUID
	}

	if info.PartUUIDs != nil {
		order.PartUUIDs = info.PartUUIDs
	}

	if info.TotalPrice > 0 {
		order.TotalPrice = info.TotalPrice
	}

	if info.TransactionUUID != nil {
		order.TransactionUUID = info.TransactionUUID
	}

	if info.PaymentMethod != nil {
		order.PaymentMethod = info.PaymentMethod
	}

	if info.Status != "" {
		order.Status = info.Status
	}

	r.data[uuid] = order

	return nil
}
