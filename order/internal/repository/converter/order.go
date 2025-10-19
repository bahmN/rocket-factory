package converter

import (
	"github.com/bahmN/rocket-factory/order/internal/model"
	repoModel "github.com/bahmN/rocket-factory/order/internal/repository/model"
)

func OrderToModel(order repoModel.OrderInfo) model.OrderInfo {
	return model.OrderInfo{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}
