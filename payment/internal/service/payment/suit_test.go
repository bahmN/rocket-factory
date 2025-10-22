package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuit struct {
	suite.Suite
	service *service
}

func (s *ServiceSuit) SetupTest() {
	s.service = NewService()
}

func (s *ServiceSuit) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
