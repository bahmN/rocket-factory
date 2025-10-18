package v1

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/converter"
	"github.com/bahmN/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	if req.GetFilter() == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	parts, err := a.inventoryService.ListParts(ctx, converter.FilterToModel(req.GetFilter()))
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "parts not found")
		}

		return nil, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
