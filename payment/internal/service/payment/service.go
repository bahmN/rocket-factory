package payment

import def "github.com/bahmN/rocket-factory/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
