package v1

import (
	"context"
	"errors"
	"log"

	"github.com/bahmN/rocket-factory/payment/internal/converter"
	"github.com/bahmN/rocket-factory/payment/internal/model"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	uuid, err := a.paymentService.PayOrder(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrRequestParamIsEmpty) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		log.Printf("unknown error: %v", err)

		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: uuid,
	}, nil
}
