package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryGRPCPort = "50052"
	paymentGRPCPort   = "50051"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) UpdateOrder(uuidOrder string, order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuidOrder] = order
}

func (s *OrderStorage) GetOrder(uuidOrder string) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuidOrder]
	if !ok {
		return nil
	}

	return order
}

type Handler struct {
	storage             *OrderStorage
	inventoryGRPCClient inventoryV1.InventoryServiceClient
	paymentGRPCClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, inventoryGRPCClient inventoryV1.InventoryServiceClient, paymentGRPCClient paymentV1.PaymentServiceClient) *Handler {
	return &Handler{
		storage:             storage,
		inventoryGRPCClient: inventoryGRPCClient,
		paymentGRPCClient:   paymentGRPCClient,
	}
}

func (h *Handler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	partsUUIDs := lo.Map(req.PartsUUID, func(item jx.Raw, _ int) string {
		id, _ := uuid.Parse(string(item))

		return id.String()
	})

	parts, err := h.inventoryGRPCClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partsUUIDs,
		},
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "one or more parts not found",
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("failed to list parts from inventory service: %v", err),
			}, nil
		}
	}

	orderUUID := uuid.NewString()

	totalPrice := lo.Reduce(parts.Parts, func(agg float64, item *inventoryV1.Part, _ int) float64 {
		return agg + item.Price
	}, 0)

	newOrder := &orderV1.Order{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartsUUID,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	h.storage.UpdateOrder(orderUUID, newOrder)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *Handler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "order not found",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID || order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.Conflict{
			Code:    409,
			Message: "order paid or cancelled",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED
	h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.CancelOrderResponse{
		Message: "order cancelled",
	}, nil
}

func (h *Handler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "order not found",
		}, nil
	}

	return order, nil
}

func (h *Handler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("empty order UUID"),
		}, nil
	}

	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "order not found",
		}, nil
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "order paid or cancelled",
		}, nil
	}

	payment, err := h.paymentGRPCClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Internal {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("payment service internal error: %v", err),
			}, nil
		}
	}

	order.PaymentMethod = orderV1.NewOptPaymentMethod(req.PaymentMethod)
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptString(payment.TransactionUuid)
	h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: payment.TransactionUuid,
	}, nil
}

func convertPaymentMethod(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}

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
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

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
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
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
