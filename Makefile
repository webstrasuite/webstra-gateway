build:
	@go build -o bin/gateway-service

run: build
	@bin/gateway-service

