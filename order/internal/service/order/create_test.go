package order

import (
	"context"
	"fmt"

	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuit) TestCreateOrderSuccess() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUUIDs = []string{partUUID1, partUUID2}
		parts     = []model.Part{
			{UUID: partUUID1, Price: gofakeit.Price(1000, 5000)},
			{UUID: partUUID2, Price: gofakeit.Price(1000, 5000)},
		}
	)
	ctx := context.Background()

	s.inventoryClient.On("ListParts", ctx, model.Filter{UUIDs: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("Create", ctx, mock.AnythingOfType("model.OrderInfo")).Return(nil)

	result, err := s.service.Create(ctx, model.CreateOrderReq{
		UserUUID:  userUUID,
		PartsUUID: partUUIDs,
	})
	s.Require().NoError(err)
	s.IsType(0.0, result.TotalPrice)
	s.IsType("", result.OrderUUID)
}

func (s *ServiceSuit) TestCreateOrderListPartsError() {
	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}
	ctx := context.Background()

	s.inventoryClient.On("ListParts", ctx, model.Filter{UUIDs: partUUIDs}).
		Return(nil, fmt.Errorf("inventory service error"))

	result, err := s.service.Create(ctx, model.CreateOrderReq{
		UserUUID:  userUUID,
		PartsUUID: partUUIDs,
	})

	s.Error(err) // тест должен зафейлиться
	s.Contains(err.Error(), "inventory service error")
	s.Equal("", result.OrderUUID)
}

func (s *ServiceSuit) TestCreateOrderEmptyParts() {
	userUUID := gofakeit.UUID()
	partUUIDs := []string{}
	ctx := context.Background()

	s.inventoryClient.On("ListParts", ctx, model.Filter{UUIDs: partUUIDs}).Return([]model.Part{}, nil)
	s.orderRepository.On("Create", ctx, mock.AnythingOfType("model.OrderInfo")).Return(nil)

	result, err := s.service.Create(ctx, model.CreateOrderReq{
		UserUUID:  userUUID,
		PartsUUID: partUUIDs,
	})

	s.Error(err)
	s.Contains(err.Error(), "parts not found")
	s.Equal("", result.OrderUUID)
}
