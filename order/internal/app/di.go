package app

import (
	"context"
	"fmt"

	orderV1API "github.com/bahmN/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/bahmN/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/bahmN/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/bahmN/rocket-factory/order/internal/config"
	"github.com/bahmN/rocket-factory/order/internal/repository"
	orderRepository "github.com/bahmN/rocket-factory/order/internal/repository/order"
	"github.com/bahmN/rocket-factory/order/internal/service"
	orderService "github.com/bahmN/rocket-factory/order/internal/service/order"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService service.OrderService

	orderRepository repository.OrderRepository

	dbPool *pgxpool.Pool

	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			inventoryClient.NewClient(d.InventoryClient(ctx)),
			paymentClient.NewClient(d.PaymentClient(ctx)),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.InitDBPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) InitDBPool(ctx context.Context) *pgxpool.Pool {
	if d.dbPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to create pool: %v", err))
		}

		d.dbPool = pool
	}

	return d.dbPool
}

func (d *diContainer) InventoryClient(ctx context.Context) inventoryV1.InventoryServiceClient {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory: %v", err))
		}

		client := inventoryV1.NewInventoryServiceClient(conn)

		d.inventoryClient = client
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) paymentV1.PaymentServiceClient {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment: %v", err))
		}

		client := paymentV1.NewPaymentServiceClient(conn)

		d.paymentClient = client
	}

	return d.paymentClient
}
