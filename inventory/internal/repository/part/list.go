package part

import (
	"context"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	"github.com/bahmN/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (r *repository) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	bsonFilter := buildMongoFilter(filter)

	cursor, err := r.coll.Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			logger.Warn(ctx, "error closing mongo cursor", zap.Error(err))
		}
	}(cursor, ctx)

	var repoParts []*repoModel.Part
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, err
	}

	if len(repoParts) == 0 {
		return nil, model.ErrPartNotFound
	}

	return converter.PartsSliceToModel(repoParts), nil
}

func buildMongoFilter(filter model.Filter) bson.M {
	query := bson.M{}
	andConditions := make([]bson.M, 0)

	if len(filter.UUIDs) > 0 {
		andConditions = append(andConditions, bson.M{"uuid": bson.M{"$in": filter.UUIDs}})
	}

	if len(filter.Names) > 0 {
		andConditions = append(andConditions, bson.M{"name": bson.M{"$in": filter.Names}})
	}

	if len(filter.Categories) > 0 {
		andConditions = append(andConditions, bson.M{"category": bson.M{"$in": filter.Categories}})
	}

	if len(filter.ManufacturerCountries) > 0 {
		andConditions = append(andConditions, bson.M{"manufacturer.country": bson.M{"$in": filter.ManufacturerCountries}})
	}

	if len(filter.Tags) > 0 {
		andConditions = append(andConditions, bson.M{"tags": bson.M{"$in": filter.Tags}})
	}

	if len(andConditions) > 0 {
		query["$and"] = andConditions
	}

	return query
}
