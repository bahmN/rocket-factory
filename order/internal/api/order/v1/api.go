package v1

import (
	"github.com/bahmN/rocket-factory/order/internal/service"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
)

type api struct {
	orderV1.UnimplementedHandler

	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
