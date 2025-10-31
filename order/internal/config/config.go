package config

import (
	"github.com/bahmN/rocket-factory/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	OrderHTTP     OrderHTTPConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		OrderHTTP:     orderCfg,
		Postgres:      postgresCfg,
		InventoryGRPC: inventoryCfg,
		PaymentGRPC:   paymentCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
