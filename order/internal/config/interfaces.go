package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type OrderHTTPConfig interface {
	Address() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsDir() string
}
