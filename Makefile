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

.PHONY: all run build deps clean lint test