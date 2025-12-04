package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/bahmN/rocket-factory/iam/internal/converter"
	"github.com/bahmN/rocket-factory/iam/internal/model"
	authV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoaimRequest) (*authV1.WhoaimResponse, error) {
	if req.SessionUuid == "" {
		return &authV1.WhoaimResponse{}, model.ErrSessionUUIDIsMissing
	}

	sessionUUID, err := uuid.Parse(req.SessionUuid)
	if err != nil {
		return &authV1.WhoaimResponse{}, model.ErrInvalidSessionUUID
	}

	response, err := a.service.Whoami(ctx, sessionUUID)
	if err != nil {
		return &authV1.WhoaimResponse{}, err
	}

	return converter.WhoamiResponseToProto(response, sessionUUID.String()), nil
}
