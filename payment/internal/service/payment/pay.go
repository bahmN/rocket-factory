package payment

import (
	"context"

	"github.com/bahmN/rocket-factory/payment/internal/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *service) PayOrder(ctx context.Context, info model.PaymentInfo) (string, error) {
	if lo.IsEmpty(info.OrderUUID) || lo.IsEmpty(info.UserUUID) {
		logger.Info(ctx, "required params is empty")
		return "", model.ErrRequestParamIsEmpty
	}

	transactionUUID := uuid.NewString()

	logger.Info(ctx, "Payment successfully", zap.String("uuid", transactionUUID))
	return transactionUUID, nil
}
