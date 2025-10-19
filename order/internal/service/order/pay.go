package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
)

func (s *service) Pay(ctx context.Context, uuid, method string) (string, error) {
	if uuid == "" {
		return "", model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return "", err
	}

	if order.Status != model.OrderStatusPENDINGPAYMENT {
		return "", model.ErrOrderPaidOrCanceled
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, method)
	if err != nil {
		return "", err
	}

	order.PaymentMethod = &method
	order.Status = model.OrderStatusPAID
	order.TransactionUUID = &transactionUUID

	err = s.repo.Update(ctx, uuid, order)
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}
