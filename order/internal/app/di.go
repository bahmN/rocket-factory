package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	orderV1API "github.com/bahmN/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/bahmN/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/bahmN/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/bahmN/rocket-factory/order/internal/config"
	kafkaConv "github.com/bahmN/rocket-factory/order/internal/converter/kafka"
	"github.com/bahmN/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/bahmN/rocket-factory/order/internal/repository"
	orderRepository "github.com/bahmN/rocket-factory/order/internal/repository/order"
	"github.com/bahmN/rocket-factory/order/internal/service"
	"github.com/bahmN/rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/bahmN/rocket-factory/order/internal/service/order"
	"github.com/bahmN/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/bahmN/rocket-factory/platform/pkg/closer"
	wrapperKafka "github.com/bahmN/rocket-factory/platform/pkg/kafka"
	wrapperKafkaConsumer "github.com/bahmN/rocket-factory/platform/pkg/kafka/consumer"
	wrapperKafkaProducer "github.com/bahmN/rocket-factory/platform/pkg/kafka/producer"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	middlewareGRPC "github.com/bahmN/rocket-factory/platform/pkg/middleware/grpc"
	HTTPMiddleware "github.com/bahmN/rocket-factory/platform/pkg/middleware/http"
	kafkaMiddleware "github.com/bahmN/rocket-factory/platform/pkg/middleware/kafka"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	authV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1API           orderV1.Handler
	orderService         service.OrderService
	orderRepository      repository.OrderRepository
	orderConsumerService service.ConsumerService
	orderProducerService service.OrderProducerService

	dbPool *pgxpool.Pool

	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient

	consumerGroup          sarama.ConsumerGroup
	orderAssembledConsumer wrapperKafka.Consumer

	orderAssembledDecoder kafkaConv.OrderAssembledDecoder
	syncProducer          sarama.SyncProducer
	orderPaidProducer     wrapperKafka.Producer

	iamClient HTTPMiddleware.IAMClient
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
			d.OrderProducerService(),
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

// OrderProducerService - Создает сервис kafka producer
func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = order_producer.NewService(d.OrderPaidProducer())
	}
	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = order_consumer.NewService(d.OrderAssembledConsumer(), d.OrderRepository(ctx), d.OrderAssembledDecoder())
	}
	return d.orderConsumerService
}

// ConsumerGroup - Создается consumer group на основе данных из конфигурации
func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания consumer group: %s\n", err.Error()))
		}

		// Добавляем закрытие ConsumerGroup
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

// OrderAssembledConsumer - Создается consumer с определенной consumer group и списком топиков для прослушивания
func (d *diContainer) OrderAssembledConsumer() wrapperKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrapperKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

// OrderAssembledDecoder - Создается декодер для входящих событий
func (d *diContainer) OrderAssembledDecoder() kafkaConv.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return d.orderAssembledDecoder
}

// SyncProducer - создает базового producer с указанными брокерами
func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания sync producer: %s\n", err.Error()))
		}

		// Добавляем закрытие producer
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error { return p.Close() })

		d.syncProducer = p
	}
	return d.syncProducer
}

// OrderPaidProducer - создает producer который отправляет в топик, заданный в конфигурации
func (d *diContainer) OrderPaidProducer() wrapperKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrapperKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.orderPaidProducer
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
