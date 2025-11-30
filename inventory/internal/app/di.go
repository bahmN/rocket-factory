package app

import (
	"context"
	"fmt"

	inventoryV1API "github.com/bahmN/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/bahmN/rocket-factory/inventory/internal/config"
	"github.com/bahmN/rocket-factory/inventory/internal/repository"
	inventoryRepository "github.com/bahmN/rocket-factory/inventory/internal/repository/part"
	"github.com/bahmN/rocket-factory/inventory/internal/service"
	inventoryService "github.com/bahmN/rocket-factory/inventory/internal/service/inventory"
	"github.com/bahmN/rocket-factory/platform/pkg/closer"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	middlewareGRPC "github.com/bahmN/rocket-factory/platform/pkg/middleware/grpc"
	authV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	iamClient middlewareGRPC.IAMClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewApi(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IAMClient(ctx context.Context) middlewareGRPC.IAMClient {
	if d.iamClient == nil {
		grpcIAM := authV1.NewAuthServiceClient(d.IAMConn(ctx))
		d.iamClient = grpcIAM
	}
	return d.iamClient
}

func (d *diContainer) IAMConn(_ context.Context) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		config.AppConfig().IamGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к IAM Service: %v", err))
	}

	closer.AddNamed("IAM client", func(ctx context.Context) error {
		if err := conn.Close(); err != nil {
			logger.Error(ctx, "❌ Ошибка при закрытии подключения с IAM Service", zap.Error(err))
			return err
		}
		return nil
	})

	return conn
}
