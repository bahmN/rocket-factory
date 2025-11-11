package integration

import (
	"context"
	"os"
	"time"

	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
)

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = defaultDatabaseName
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()

	parts := []*repoModel.Part{
		{
			UUID:          partUUID,
			Name:          "Двигатель",
			Description:   "Мощный ракетный двигатель",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      1, // ENGINE
			Dimensions: &repoModel.Dimensions{
				Length: 200,
				Width:  100,
				Height: 120,
				Weight: 300,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "RocketMotors",
				Country: "Russia",
				Website: "https://rocketmotors.example.com",
			},
			Tags:      []string{"основной", "мотор"},
			Metadata:  map[string]string{"серия": "X100"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	docs := make([]interface{}, len(parts))
	for i, part := range parts {
		docs[i] = part
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = defaultDatabaseName
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertMany(ctx, docs)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}
