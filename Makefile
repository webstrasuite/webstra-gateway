BINARY=webstraGateway

build:
	@echo "Building binary..."
	@go build -o bin/${BINARY} cmd/main.go

run: build
	@echo "Running binary..."
	@bin/${BINARY}

test:
	@echo "Running tests..."
	@go test -v ./...

proto:
	@echo "Generating protobuf definitions..."
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pb/*.proto
	@echo "Done generating!"

.PHONY: proto