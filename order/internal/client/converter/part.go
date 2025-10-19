package converter

import (
	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func metadataToModel(metadata map[string]*structpb.Value) map[string]string {
	metadataProto := make(map[string]string)
	for k, v := range metadata {
		metadataProto[k] = v.String()
	}

	return metadataProto
}

func PartsFilterToProto(filter model.Filter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func PartsListToModel(list *inventoryV1.ListPartsResponse) []model.Part {
	parts := make([]model.Part, 0, len(list.Parts))

	for _, part := range list.Parts {
		parts = append(parts, model.Part{
			UUID:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      int32(part.Category),
			Dimensions: &model.Dimensions{
				Length: part.Dimensions.Length,
				Width:  part.Dimensions.Width,
				Height: part.Dimensions.Height,
				Weight: part.Dimensions.Weight,
			},
			Manufacturer: &model.Manufacturer{
				Name:    part.Manufacturer.Name,
				Cuntry:  part.Manufacturer.Cuntry,
				Website: part.Manufacturer.Website,
			},
			Tags:      part.Tags,
			Metadata:  metadataToModel(part.Metadata),
			CreatedAt: part.CreatedAt.AsTime(),
			UpdatedAt: part.UpdatedAt.AsTime(),
		})
	}

	return parts
}

func PaymentMethodToProto(paymentMethod string) paymentV1.PaymentMethod {
	switch paymentMethod {
	case "CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "INVESTOR_MONEY":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}
