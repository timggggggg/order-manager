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

compose-up:
	cd build && make compose-up;

goose-up:
	cd build && make goose-up;

compose-stop:
	cd build && make compose-stop;

goose-stop:
	cd build && make goose-stop;

.PHONY: all run build deps clean lint test