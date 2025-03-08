SERVICE := postgres
DB_NAME := orderdb
USER := postgres
PASSWORD := postgres

HOST := localhost
PORT := 5432

GOOSE_DRIVER := postgres
DB_STRING := "postgres://$(USER):$(PASSWORD)@$(HOST):$(PORT)/$(DB_NAME)?sslmode=disable"
MIGRATIONS := ./migrations
ENV_NAME = "order-manager"

all: compose-up goose-up

# docker compose
compose-up:
	@docker-compose -p ${ENV_NAME} up -d $(SERVICE)

compose-down:
	@docker-compose -p ${ENV_NAME} down $(SERVICE)

compose-stop:
	@docker-compose -p ${ENV_NAME} stop $(SERVICE)

compose-start:
	@docker-compose -p ${ENV_NAME} start $(SERVICE)

compose-ps:
	@docker-compose -p ${ENV_NAME} ps $(SERVICE)

compose-rm:
	@docker-compose -p ${ENV_NAME} rm $(SERVICE)


# migrations
goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	@goose -dir $(MIGRATIONS) $(GOOSE_DRIVER) $(DB_STRING) create rename_me sql

goose-up:
	@goose -dir $(MIGRATIONS) $(GOOSE_DRIVER) $(DB_STRING) up

goose-down:
	@goose -dir $(MIGRATIONS) $(GOOSE_DRIVER) $(DB_STRING) down

goose-down-all:
	@goose -dir $(MIGRATIONS) $(GOOSE_DRIVER) $(DB_STRING) down-to 0

goose-restart: goose-down-all goose-up

