package order

import (
	"errors"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuit) TestPayOrderSuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = "card"
		transactionUUID = gofakeit.UUID()
		order           = model.OrderInfo{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.OrderInfo")).Return(nil)

	result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

	s.Require().NoError(err)
	s.Equal(transactionUUID, result)
}

func (s *ServiceSuit) TestPayOrderAlreadyPaidOrCanceled() {
	orderUUID := gofakeit.UUID()
	order := model.OrderInfo{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPAID,
	}

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)

	result, err := s.service.Pay(s.ctx, orderUUID, "card")
	s.Require().ErrorIs(err, model.ErrOrderPaidOrCanceled)
	s.Equal("", result)
}

func (s *ServiceSuit) TestPayOrderPaymentClientError() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	order := model.OrderInfo{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, "card").Return("", errors.New("payment failed"))

	result, err := s.service.Pay(s.ctx, orderUUID, "card")
	s.Require().Error(err)
	s.Equal("", result)
}
