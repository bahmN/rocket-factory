package testcontainers

// MongoDB constants
const (
	// Переопределяем grpc host т.к. из вне по localhost не достучаться
	GrpcHostKey = "GRPC_HOST"
	GrpcHost    = "0.0.0.0"

	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)
