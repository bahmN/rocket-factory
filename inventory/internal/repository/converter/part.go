package converter

import (
	"github.com/bahmN/rocket-factory/inventory/internal/model"
	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
)

func PartToModel(part *repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Dimensions:    dimensionToModel(part.Dimensions),
		Manufacturer:  manufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func FilterToRepoModel(filter model.Filter) *repoModel.Filter {
	return &repoModel.Filter{
		UUIDs:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            filter.Categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func PartsSliceToModel(parts []*repoModel.Part) []model.Part {
	var partsModel []model.Part
	for _, part := range parts {
		partsModel = append(partsModel, PartToModel(part))
	}

	return partsModel
}

func dimensionToModel(dim *repoModel.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}
func manufacturerToModel(manu *repoModel.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    manu.Name,
		Cuntry:  manu.Cuntry,
		Website: manu.Website,
	}
}
