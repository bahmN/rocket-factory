package v1

import (
	"context"
	"log"

	"github.com/bahmN/rocket-factory/inventory/internal/converter"
	"github.com/bahmN/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with %v uuid not found", req.Uuid)
		}

		log.Printf("unknown error: %v", err)

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
