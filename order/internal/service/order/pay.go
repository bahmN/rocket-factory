package order

import (
	"context"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	uuidgen "github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) Pay(ctx context.Context, uuid, method string) (string, error) {
	if uuid == "" {
		logger.Info(ctx, "uuid is empty")
		return "", model.ErrEmptyUUID
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		logger.Info(ctx, "error in getting order", zap.String("uuid", uuid))
		return "", err
	}

	if order.Status != model.OrderStatusPENDINGPAYMENT {
		logger.Info(ctx, "order is not pending payment", zap.String("uuid", uuid), zap.String("status", order.Status))
		return "", model.ErrOrderPaidOrCanceled
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, method)
	if err != nil {
		logger.Info(ctx, "error in paying order", zap.String("uuid", uuid), zap.Error(err))
		return "", err
	}

	order.PaymentMethod = &method
	order.Status = model.OrderStatusPAID
	order.TransactionUUID = &transactionUUID

	err = s.repo.Update(ctx, uuid, order)
	if err != nil {
		return "", err
	}

	err = s.orderPaidProducer.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuidgen.NewString(),
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   method,
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		return "", err
	}

	logger.Info(ctx, "order payed successfully", zap.String("transaction_uuid", transactionUUID), zap.String("method", method))
	return transactionUUID, nil
}
