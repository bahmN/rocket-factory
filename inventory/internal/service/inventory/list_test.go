package inventory

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuit) TestListPartsSuccess() {
	filter := model.Filter{UUIDs: []string{gofakeit.UUID(), gofakeit.UUID()}}
	expectedParts := []model.Part{
		{UUID: filter.UUIDs[0], Price: gofakeit.Price(1000, 5000)},
		{UUID: filter.UUIDs[1], Price: gofakeit.Price(1000, 5000)},
	}

	ctx := context.Background()
	s.inventoryRepository.On("ListParts", ctx, filter).Return(expectedParts, nil)

	result, err := s.service.ListParts(ctx, filter)
	s.Require().NoError(err)
	s.Equal(expectedParts, result)
}

func (s *ServiceSuit) TestListPartsRepoError() {
	filter := model.Filter{UUIDs: []string{gofakeit.UUID()}}
	ctx := context.Background()

	s.inventoryRepository.On("ListParts", ctx, filter).Return(nil, errors.New("repo error"))

	result, err := s.service.ListParts(ctx, filter)
	s.Require().Error(err)
	s.Nil(result)
}
