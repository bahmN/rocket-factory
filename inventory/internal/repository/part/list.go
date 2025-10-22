package part

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	"github.com/bahmN/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	r.mu.RLock()
	parts := make([]*repoModel.Part, 0, len(r.data))
	for _, part := range r.data {
		parts = append(parts, part)
	}
	r.mu.RUnlock()

	parts = applyUUIDFilter(parts, filter)
	parts = applyNameFilter(parts, filter)
	parts = applyCategoryFilter(parts, filter)
	parts = applyCountryFilter(parts, filter)
	parts = applyTagsFilter(parts, filter)

	if len(parts) == 0 {
		return nil, model.ErrPartNotFound
	}

	return converter.PartsSliceToModel(parts), nil
}

func applyUUIDFilter(parts []*repoModel.Part, filter model.Filter) []*repoModel.Part {
	if len(filter.UUIDs) == 0 {
		return parts
	}
	uuidSet := makeSetString(filter.UUIDs)
	tmp := parts[:0]
	for _, part := range parts {
		if _, ok := uuidSet[part.UUID]; ok {
			tmp = append(tmp, part)
		}
	}
	return tmp
}

func applyNameFilter(parts []*repoModel.Part, filter model.Filter) []*repoModel.Part {
	if len(filter.Names) == 0 {
		return parts
	}
	nameSet := makeSetString(filter.Names)
	tmp := parts[:0]
	for _, part := range parts {
		if _, ok := nameSet[part.Name]; ok {
			tmp = append(tmp, part)
		}
	}
	return tmp
}

func applyCategoryFilter(parts []*repoModel.Part, filter model.Filter) []*repoModel.Part {
	if len(filter.Categories) == 0 {
		return parts
	}
	catSet := makeSetInt32(filter.Categories)
	tmp := parts[:0]
	for _, part := range parts {
		if _, ok := catSet[part.Category]; ok {
			tmp = append(tmp, part)
		}
	}
	return tmp
}

func applyCountryFilter(parts []*repoModel.Part, filter model.Filter) []*repoModel.Part {
	if len(filter.ManufacturerCountries) == 0 {
		return parts
	}
	countrySet := makeSetString(filter.ManufacturerCountries)
	tmp := parts[:0]
	for _, part := range parts {
		if part.Manufacturer != nil {
			if _, ok := countrySet[part.Manufacturer.Cuntry]; ok {
				tmp = append(tmp, part)
			}
		}
	}
	return tmp
}

func applyTagsFilter(parts []*repoModel.Part, filter model.Filter) []*repoModel.Part {
	if len(filter.Tags) == 0 {
		return parts
	}
	tagSet := makeSetString(filter.Tags)
	tmp := parts[:0]
	for _, part := range parts {
		found := false
		for _, tag := range part.Tags {
			if _, ok := tagSet[tag]; ok {
				found = true
				break
			}
		}
		if found {
			tmp = append(tmp, part)
		}
	}
	return tmp
}

// вспомогательные функции для множеств
func makeSetString(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, s := range items {
		set[s] = struct{}{}
	}
	return set
}

func makeSetInt32(items []int32) map[int32]struct{} {
	set := make(map[int32]struct{}, len(items))
	for _, s := range items {
		set[s] = struct{}{}
	}
	return set
}
