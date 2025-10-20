package converter

import (
	"github.com/bahmN/rocket-factory/order/internal/model"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/go-faster/jx"
	"github.com/google/uuid"
)

func convertStringsToRaw(parts []string) []jx.Raw {
	raws := make([]jx.Raw, len(parts))
	for i, s := range parts {
		raws[i] = jx.Raw([]byte(`"` + s + `"`))
	}
	return raws
}

func CreateOrderToModel(order *orderV1.CreateOrderRequest) (model.CreateOrderReq, error) {
	partsUUIDs := make([]string, 0, len(order.PartsUUID))
	for _, item := range order.PartsUUID {
		id, err := uuid.Parse(string(item))
		if err != nil {
			return model.CreateOrderReq{}, err
		}
		partsUUIDs = append(partsUUIDs, id.String())
	}

	return model.CreateOrderReq{
		UserUUID:  order.UserUUID,
		PartsUUID: partsUUIDs,
	}, nil
}

func CreateOrderToAPI(order *model.CreateOrderResp) *orderV1.CreateOrderResponse {
	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}
}

func OrderToAPI(order model.OrderInfo) *orderV1.Order {
	var transactionUUID orderV1.OptString
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.OptString{Value: *order.TransactionUUID, Set: true}
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.TransactionUUID != nil {
		paymentMethod = orderV1.OptPaymentMethod{Value: paymentMethodToAPI(*order.PaymentMethod), Set: true}
	}
	return &orderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       convertStringsToRaw(order.PartUUIDs),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderV1.OrderStatus(order.Status),
	}
}

func paymentMethodToAPI(method string) orderV1.PaymentMethod {
	switch method {
	case "CARD":
		return orderV1.PaymentMethodCARD
	case "SBP":
		return orderV1.PaymentMethodSBP
	case "CREDIT_CARD":
		return orderV1.PaymentMethodCREDITCARD
	case "INVESTOR_MONEY":
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodCARD
	}
}
