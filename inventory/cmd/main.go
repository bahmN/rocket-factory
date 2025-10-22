package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryV1API "github.com/bahmN/rocket-factory/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/bahmN/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/bahmN/rocket-factory/inventory/internal/service/inventory"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50053"

func main() {
	lis, err := net.Listen("tcp", "localhost"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	repository := inventoryRepository.NewRepository()
	repository.InitTestData()
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
