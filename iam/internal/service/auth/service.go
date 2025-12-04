package auth

import (
	"github.com/bahmN/rocket-factory/iam/internal/repository"
	"github.com/bahmN/rocket-factory/platform/pkg/hasher"
)

type service struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	hasher      hasher.PasswordHasher
}

func NewService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, hasher hasher.PasswordHasher) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		hasher:      hasher,
	}
}
