## Simple server with gRPC

### .env variables
Set it to you needed

    SERVER_ADDRESS_PORT=localhost:50051

### Protobuf generate

    protoc --go_out=. --go-grpc_out=. ./pkg/proto/hello/hello.proto
