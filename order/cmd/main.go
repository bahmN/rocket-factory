package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderV1API "github.com/bahmN/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/bahmN/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/bahmN/rocket-factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/bahmN/rocket-factory/order/internal/repository/order"
	orderService "github.com/bahmN/rocket-factory/order/internal/service/order"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryGRPCPort = "50053"
	paymentGRPCPort   = "50051"
)

func main() {
	inventoryConn, err := grpc.NewClient(
		"localhost:"+inventoryGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("inventory connection close error: %v", err)
		}
	}()
	iClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		"localhost:"+paymentGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("payment connection close error: %v", err)
		}
	}()
	pClient := paymentV1.NewPaymentServiceClient(paymentConn)

	repository := orderRepository.NewRepository()
	service := orderService.NewService(
		repository,
		inventoryClient.NewClient(iClient),
		paymentClient.NewClient(pClient))
	api := orderV1API.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		panic("Failed to create server: %v" + err.Error())
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Starting server on port %s", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown server: %v", err)
	}

	log.Println("Server gracefully stopped")
}
