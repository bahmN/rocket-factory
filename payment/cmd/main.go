package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiService "github.com/bahmN/rocket-factory/payment/internal/api/payment/v1"
	"github.com/bahmN/rocket-factory/payment/internal/interceptor"
	paymentService "github.com/bahmN/rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50051"

func main() {
	lis, err := net.Listen("tcp", "localhost"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.ValidatorInterceptor()),
		),
	)

	service := paymentService.NewService()
	api := apiService.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

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
