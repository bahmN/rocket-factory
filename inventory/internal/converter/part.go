package converter

import (
	"github.com/bahmN/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FilterToModel(filter *inventoryV1.PartsFilter) model.Filter {
	categories := make([]int32, len(filter.Categories))
	for i, c := range filter.Categories {
		categories[i] = int32(c)
	}

	return model.Filter{
		UUIDs:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func PartsToProto(parts []model.Part) []*inventoryV1.Part {
	list := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		list[i] = PartToProto(part)
	}

	return list
}

func PartToProto(part model.Part) *inventoryV1.Part {
	var createdAt *timestamppb.Timestamp
	if !part.CreatedAt.IsZero() {
		createdAt = timestamppb.New(part.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if !part.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func categoryToProto(category int32) inventoryV1.Category {
	switch category {
	case 1:
		return inventoryV1.Category_CATEGORY_ENGINE
	case 2:
		return inventoryV1.Category_CATEGORY_FUEL
	case 3:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case 4:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED
	}
}

func dimensionsToProto(dim *model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func manufacturerToProto(man *model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    man.Name,
		Cuntry:  man.Cuntry,
		Website: man.Website,
	}
}

func metadataToProto(metadata map[string]string) map[string]*structpb.Value {
	metadataProto := make(map[string]*structpb.Value)
	for k, v := range metadata {
		metadataProto[k] = structpb.NewStringValue(v)
	}

	return metadataProto
}
