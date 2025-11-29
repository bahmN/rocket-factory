package v1

import (
	"github.com/bahmN/rocket-factory/iam/internal/service"
	userV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	service service.UserService
}

func NewUserAPI(service service.UserService) *api {
	return &api{service: service}
}
