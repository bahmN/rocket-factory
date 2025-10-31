package config

import (
	"github.com/bahmN/rocket-factory/payment/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
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

	paymentCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		PaymentGRPC: paymentCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
