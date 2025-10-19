package order

import (
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

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(expectedOrder, nil)

	result, err := s.service.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(expectedOrder, result)
}

func (s *ServiceSuit) TestGetOrderEmptyUUID() {
	result, err := s.service.Get(s.ctx, "")
	s.Require().ErrorIs(err, model.ErrEmptyUUID)
	s.Equal(model.OrderInfo{}, result)
}

func (s *ServiceSuit) TestGetOrderRepoError() {
	orderUUID := gofakeit.UUID()
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(model.OrderInfo{}, errors.New("repo error"))

	result, err := s.service.Get(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Equal(model.OrderInfo{}, result)
}
