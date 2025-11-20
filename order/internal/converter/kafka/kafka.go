package kafka

import "github.com/bahmN/rocket-factory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}
