package payment

import (
	"context"
	"log"

	"github.com/bahmN/rocket-factory/payment/internal/model"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (s *service) PayOrder(ctx context.Context, info model.PaymentInfo) (string, error) {
	if lo.IsEmpty(info.OrderUUID) || lo.IsEmpty(info.UserUUID) || lo.IsEmpty(info.PaymentMethod) {
		return "", model.ErrRequestParamIsEmpty
	}

	transactionUUID := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction UUID: %s", transactionUUID)

	return transactionUUID, nil
}
