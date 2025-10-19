package service

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req model.CreateOrderReq) (model.CreateOrderResp, error)
	Get(ctx context.Context, uuid string) (model.OrderInfo, error)
	Pay(ctx context.Context, uuid, method string) (string, error)
	Cancel(ctx context.Context, uuid string) (string, error)
}
