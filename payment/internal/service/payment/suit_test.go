package payment

import (
	"testing"

	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"github.com/stretchr/testify/suite"
)

type ServiceSuit struct {
	suite.Suite
	service *service
}

func (s *ServiceSuit) SetupTest() {
	logger.SetNopLogger()

	s.service = NewService()
}

func (s *ServiceSuit) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuit))
}
