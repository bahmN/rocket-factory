package repository

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.OrderInfo) error
	Get(ctx context.Context, uuid string) (model.OrderInfo, error)
	Update(ctx context.Context, uuid string, order model.OrderInfo) error
}
