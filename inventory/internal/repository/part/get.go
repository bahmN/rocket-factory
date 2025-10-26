package part

import (
	"context"
	"errors"

	"github.com/bahmN/rocket-factory/inventory/internal/model"
	repoConverter "github.com/bahmN/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	var part repoModel.Part
	err := r.coll.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return repoConverter.PartToModel(&part), nil
}
