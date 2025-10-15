package payment

import (
	"context"
	"log"

	"github.com/bahmN/rocket-factory/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) PayOrder(ctx context.Context, info model.PaymentInfo) (string, error) {
	transactionUUID := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction UUID: %s", transactionUUID)

	return transactionUUID, nil
}
