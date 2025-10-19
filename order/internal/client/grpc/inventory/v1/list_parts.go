package v1

import (
	"context"

	clientConverter "github.com/bahmN/rocket-factory/order/internal/client/converter"
	"github.com/bahmN/rocket-factory/order/internal/model"
	generatedInventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return clientConverter.PartsListToModel(parts), nil
}
