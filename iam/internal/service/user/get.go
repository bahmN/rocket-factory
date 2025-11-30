package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/bahmN/rocket-factory/iam/internal/model"
)

func (s *service) Get(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	return s.userRepo.Get(ctx, userUUID)
}
