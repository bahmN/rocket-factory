package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryV1API "github.com/bahmN/rocket-factory/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/bahmN/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/bahmN/rocket-factory/inventory/internal/service/inventory"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50053"

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load .env file: %v\n", err)
	}

	dbURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("failed to connect to database: %v\n", err)
	}
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("failed to disconnect: %v\n", err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping db: %v", err)
		return
	}

	db := client.Database("inventory-service")

	lis, err := net.Listen("tcp", "localhost"+grpcPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()

	repository := inventoryRepository.NewRepository(db)
	err = repository.InitTestData(ctx)
	if err != nil {
		log.Printf("failed to init test data: %v\n", err)
	}

	service := inventoryService.NewService(repository)
	api := inventoryV1API.NewApi(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %s\n", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
