package converter

import (
	"github.com/bahmN/rocket-factory/payment/internal/model"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(info *paymentV1.PayOrderRequest) model.PaymentInfo {
	return model.PaymentInfo{
		OrderUUID:     info.OrderUuid,
		UserUUID:      info.UserUuid,
		PaymentMethod: int32(info.PaymentMethod),
	}
}
