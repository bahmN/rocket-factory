package order

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuit) TestGetOrderSuccess() {
	orderUUID := gofakeit.UUID()
	expectedOrder := model.OrderInfo{
		OrderUUID: orderUUID,
		UserUUID:  gofakeit.UUID(),
		Status:    model.OrderStatusPENDINGPAYMENT,
	}
	ctx := context.Background()

	s.orderRepository.On("Get", ctx, orderUUID).Return(expectedOrder, nil)

	result, err := s.service.Get(ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(expectedOrder, result)
}

func (s *ServiceSuit) TestGetOrderEmptyUUID() {
	ctx := context.Background()
	result, err := s.service.Get(ctx, "")
	s.Require().ErrorIs(err, model.ErrEmptyUUID)
	s.Equal(model.OrderInfo{}, result)
}

func (s *ServiceSuit) TestGetOrderRepoError() {
	orderUUID := gofakeit.UUID()
	ctx := context.Background()

	s.orderRepository.On("Get", ctx, orderUUID).Return(model.OrderInfo{}, errors.New("repo error"))

	result, err := s.service.Get(ctx, orderUUID)
	s.Require().Error(err)
	s.Equal(model.OrderInfo{}, result)
}
