package order

import (
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

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.OrderInfo")).Return(nil)

	result, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal("order cancelled", result)
	s.orderRepository.AssertCalled(s.T(), "Update", s.ctx, orderUUID, mock.MatchedBy(func(o model.OrderInfo) bool {
		return o.Status == model.OrderStatusCANCELLED
	}))
}

func (s *ServiceSuit) TestCancelOrderEmptyUUID() {
	result, err := s.service.Cancel(s.ctx, "")
	s.Require().ErrorIs(err, model.ErrEmptyUUID)
	s.Equal("", result)
}

func (s *ServiceSuit) TestCancelOrderRepoError() {
	orderUUID := gofakeit.UUID()
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(model.OrderInfo{}, errors.New("repo error"))

	result, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Equal("", result)
}

func (s *ServiceSuit) TestCancelOrderAlreadyPaidOrCancelled() {
	orderUUID := gofakeit.UUID()
	order := model.OrderInfo{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPAID,
	}

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)

	result, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().ErrorIs(err, model.ErrOrderPaidOrCanceled)
	s.Equal("", result)
}
