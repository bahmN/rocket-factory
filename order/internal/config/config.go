package config

import (
	"os"

	"github.com/bahmN/rocket-factory/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	OrderHTTP              OrderHTTPConfig
	Postgres               PostgresConfig
	InventoryGRPC          InventoryGRPCConfig
	PaymentGRPC            PaymentGRPCConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderPaidProducer      OrderPaidProducerConfig
	Logger                 LoggerConfig
	IamGRPC                IAMConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	orderHTTPConfig, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerConfig, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidProducerConfig, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	iamGRPCConfig, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		OrderHTTP:              orderHTTPConfig,
		Postgres:               postgresConfig,
		InventoryGRPC:          inventoryGRPCConfig,
		PaymentGRPC:            paymentGRPCConfig,
		OrderAssembledConsumer: orderAssembledConsumerConfig,
		OrderPaidProducer:      orderPaidProducerConfig,
		Kafka:                  kafkaConfig,
		Logger:                 loggerConfig,
		IamGRPC:                iamGRPCConfig,
	}
	return nil
}

func AppConfig() *config { return appConfig }
