module github.com/bahmN/rocket-factory/inventory

go 1.25.0

replace github.com/bahmN/rocket-factory/shared => ../shared

require (
	github.com/bahmN/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/go-faster/errors v0.7.1
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
