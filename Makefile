SERVER_BINARY_NAME=server
CLIENT_BINARY_NAME=client

SERVER_MAIN_FILE=./cmd/server/main.go
CLIENT_MAIN_FILE=./cmd/client/main.go

all: run-server

run-server: build-server
	./$(SERVER_BINARY_NAME)

client: build-client
	./$(CLIENT_BINARY_NAME)

build-server: deps
	go build -o $(SERVER_BINARY_NAME) $(SERVER_MAIN_FILE)

build-client: deps
	go build -o $(CLIENT_BINARY_NAME) $(CLIENT_MAIN_FILE)

deps:
	go mod tidy
	go mod download

clean:
	rm -f $(SERVER_BINARY_NAME)
	rm -f $(CLIENT_BINARY_NAME)

lint: 
	golangci-lint run

test:
	@go clean -testcache
	@go test -v -cover ./tests/integration ./internal/handlers

test-unit:
	@go clean -testcache
	@go test -v -cover ./internal/handlers

compose-up:
	cd build && make compose-up;

migrations-up:
	cd build && make migrations-up;

compose-stop:
	cd build && make compose-stop;

migrations-down:
	cd build && make migrations-down;

genproto:
	protoc -I=proto --go_out=./pkg --go_opt=paths=source_relative \
	--go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
	proto/api/api.proto proto/service/service.proto

.PHONY: all run build deps clean lint test genproto client