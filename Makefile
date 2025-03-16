BINARY_NAME=app

MAIN_FILE=./cmd/server/main.go

all: run

run: build
	./$(BINARY_NAME)

build: deps lint
	go build -o $(BINARY_NAME) $(MAIN_FILE)

deps:
	go mod tidy
	go mod download

clean:
	rm -f $(BINARY_NAME)

lint: 
	golangci-lint run

test:
	@go clean -testcache
	@go test -cover ./internal/handlers

test-unit:
	@go clean -testcache
	@go test -v -cover ./internal/handlers/unit_test.go ./internal/handlers/accept_order.go ./internal/handlers/issue_order.go ./internal/handlers/list_history.go ./internal/handlers/list_order.go ./internal/handlers/list_return.go ./internal/handlers/return_order.go ./internal/handlers/withdraw_order.go

compose-up:
	cd build && make compose-up;

goose-up:
	cd build && make goose-up;

compose-stop:
	cd build && make compose-stop;

goose-stop:
	cd build && make goose-stop;

.PHONY: all run build deps clean lint test