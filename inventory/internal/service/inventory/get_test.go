package inventory

import (
	"errors"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuit) TestGetPartSuccess() {
	partUUID := gofakeit.UUID()
	expectedPart := model.Part{
		UUID:  partUUID,
		Price: gofakeit.Price(1000, 5000),
	}

	s.inventoryRepository.On("GetPart", s.ctx, partUUID).Return(expectedPart, nil)

	result, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().NoError(err)
	s.Equal(expectedPart, result)
}

func (s *ServiceSuit) TestGetPartRepoError() {
	partUUID := gofakeit.UUID()
	s.inventoryRepository.On("GetPart", s.ctx, partUUID).Return(model.Part{}, errors.New("repo error"))

	result, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().Error(err)
	s.Equal(model.Part{}, result)
}
