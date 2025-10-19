package payment

import (
	"github.com/bahmN/rocket-factory/payment/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuit) TestPayOrderSuccess() {
	paymentInfo := model.PaymentInfo{
		OrderUUID:     gofakeit.UUID(),
		UserUUID:      gofakeit.UUID(),
		PaymentMethod: int32(gofakeit.Float32Range(0, 4)),
	}

	uuid, err := s.service.PayOrder(s.ctx, paymentInfo)

	s.Require().NoError(err)
	s.NotNil(uuid)
	s.IsType("", uuid, "uuid должен быть строкой")
}
