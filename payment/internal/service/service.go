package service

import (
	"context"

	"github.com/bahmN/rocket-factory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, info model.PaymentInfo) (string, error)
}
