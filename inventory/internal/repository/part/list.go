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

	repoFilter := converter.FilterToRepoModel(filter)

	if len(repoFilter.UUIDs) > 0 {
		uuidSet := make(map[string]struct{}, len(repoFilter.UUIDs))
		for _, u := range repoFilter.UUIDs {
			uuidSet[u] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if _, ok := uuidSet[part.UUID]; ok {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(repoFilter.Names) > 0 {
		nameSet := make(map[string]struct{}, len(repoFilter.Names))
		for _, n := range repoFilter.Names {
			nameSet[n] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if _, ok := nameSet[part.Name]; ok {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(filter.Categories) > 0 {
		catSet := make(map[int32]struct{}, len(filter.Categories))
		for _, c := range filter.Categories {
			catSet[c] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if _, ok := catSet[part.Category]; ok {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(filter.ManufacturerCountries) > 0 {
		countrySet := make(map[string]struct{}, len(filter.ManufacturerCountries))
		for _, c := range filter.ManufacturerCountries {
			countrySet[c] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if part.Manufacturer != nil {
				if _, ok := countrySet[part.Manufacturer.Cuntry]; ok {
					tmp = append(tmp, part)
				}
			}
		}
		parts = tmp
	}

	if len(filter.Tags) > 0 {
		tagSet := make(map[string]struct{}, len(filter.Tags))
		for _, t := range filter.Tags {
			tagSet[t] = struct{}{}
		}
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
		parts = tmp
	}

	if len(parts) == 0 {
		return nil, model.ErrPartNotFound
	}

	return converter.PartsSliceToModel(parts), nil
}
