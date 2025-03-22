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

.PHONY: all run build deps clean lint test