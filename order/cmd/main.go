package main

import (
	"sync"
	"time"

	orderV1 "github.com/bahmN/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

const (
	UNKNOWN        = 0
	CARD           = 1
	SBP            = 2
	CREDIT_CARD    = 3
	INVESTOR_MONEY = 4
)

type OrderStorage struct {
	mu    sync.RWMutex
	order map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		order: make(map[string]*orderV1.Order),
	}
}

type Handler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func main() {

}
