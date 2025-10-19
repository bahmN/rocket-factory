package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (s *service) Create(ctx context.Context, req model.CreateOrderReq) (model.CreateOrderResp, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.Filter{
		UUIDs: req.PartsUUID,
	})
	if err != nil {
		return model.CreateOrderResp{}, model.ErrPartsNotFound
	}

	orderUUID := uuid.NewString()

	totalPrice := lo.Reduce(parts, func(agg float64, item model.Part, _ int) float64 {
		return agg + item.Price
	}, 0)

	newOrder := model.OrderInfo{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartsUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	err = s.repo.Update(ctx, orderUUID, newOrder)
	if err != nil {
		return model.CreateOrderResp{}, err
	}

	return model.CreateOrderResp{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
