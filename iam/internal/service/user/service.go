package user

import (
	"github.com/bahmN/rocket-factory/iam/internal/repository"
	def "github.com/bahmN/rocket-factory/iam/internal/service"
	"github.com/bahmN/rocket-factory/platform/pkg/hasher"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepo repository.UserRepository
	hasher   hasher.PasswordHasher
}

func NewService(userRepo repository.UserRepository, hasher hasher.PasswordHasher) *service {
	return &service{
		userRepo: userRepo,
		hasher:   hasher,
	}
}
