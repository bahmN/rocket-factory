package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50051"

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer

	mu sync.RWMutex
}

func (s *PaymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//TODO: add validating req

	transactionUUID := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	service := &PaymentService{
		mu: sync.RWMutex{},
	}

	paymentV1.RegisterPaymentServiceServer(s, service)

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
