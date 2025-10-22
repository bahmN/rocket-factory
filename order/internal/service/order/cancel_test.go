package order

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuit) TestCancelOrderSuccess() {
	orderUUID := gofakeit.UUID()
	order := model.OrderInfo{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	ctx := context.Background()

	s.orderRepository.On("Get", ctx, orderUUID).Return(order, nil)
	s.orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.OrderInfo")).Return(nil)

	err := s.service.Cancel(ctx, orderUUID)
	s.Require().NoError(err)
	s.orderRepository.AssertCalled(s.T(), "Update", ctx, orderUUID, mock.MatchedBy(func(o model.OrderInfo) bool {
		return o.Status == model.OrderStatusCANCELLED
	}))
}

func (s *ServiceSuit) TestCancelOrderEmptyUUID() {
	ctx := context.Background()
	err := s.service.Cancel(ctx, "")

	s.Require().ErrorIs(err, model.ErrEmptyUUID)
}

func (s *ServiceSuit) TestCancelOrderRepoError() {
	orderUUID := gofakeit.UUID()
	ctx := context.Background()

	s.orderRepository.On("Get", ctx, orderUUID).Return(model.OrderInfo{}, errors.New("repo error"))

	err := s.service.Cancel(ctx, orderUUID)
	s.Require().Error(err)
}

func (s *ServiceSuit) TestCancelOrderAlreadyPaidOrCancelled() {
	orderUUID := gofakeit.UUID()
	order := model.OrderInfo{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPAID,
	}
	ctx := context.Background()

	s.orderRepository.On("Get", ctx, orderUUID).Return(order, nil)

	err := s.service.Cancel(ctx, orderUUID)
	s.Require().ErrorIs(err, model.ErrOrderPaidOrCanceled)
}
