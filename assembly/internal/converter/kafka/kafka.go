package kafka

import "github.com/bahmN/rocket-factory/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
