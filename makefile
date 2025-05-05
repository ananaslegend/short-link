BINARY_NAME=app
MAIN_PATH=./cmd/api
CMD_NAME_ARG=api

REDIS_PORT=6379
REDIS_IMAGE=redis:7.2.5-alpine
REDIS_CONTAINER_NAME=sl-redis
REDIS_PASSWORD=pass

POSTGRES_PORT=5432
POSTGRES_IMAGE=postgres:16.3-alpine
POSTGRES_CONTAINER_NAME=sl-postgres
POSTGRES_USER=sl
POSTGRES_PASSWORD=sl_password
POSTGRES_DB=sl
POSTGRES_HOST=localhost
POSTGRES_MIGRATION_PATH=./migrations/postgres
POSTGRES_MIGRATION_DRIVER=postgres

CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=8123
CLICKHOUSE_USER=sl
CLICKHOUSE_PASSWORD=sl_password
CLICKHOUSE_DATABASE=statistics
CLICKHOUSE_IMAGE=clickhouse/clickhouse-server:24.4-alpine
CLICKHOUSE_CONTAINER_NAME := sl-clickhouse-server
CLICKHOUSE_MIGRATION_PATH=./migrations/clickhouse
CLICKHOUSE_MIGRATION_DRIVER=clickhouse

# Colors for output
GREEN = "\033[0;32m"
YELLOW = "\033[0;33m"
RED = "\033[0;31m"
BLUE = "\033[0;34m"
NO_COLOR = "\033[0m"

.PHONY: help
help:
	@echo $(BLUE)"Available commands:"$(NO_COLOR)
	@echo "  lint                   		- Run golangci-lint"
	@echo "  build                  		- Build the application"
	@echo "  run                    		- Run the application locally"
	@echo "  test                   		- Run tests"
	@echo "  clean                  		- Clean build artifacts"
	@echo "  tidy                   		- Tidy go.mod"
	@echo "  update-deps            		- Update dependencies"
	@echo "  docker-run             		- Run application in Docker using default .env file"
	@echo "  docker-run env={env}   		- Run application in Docker using .env.{env} file"
	@echo "                           			Example: make docker-run env=dev"
	@echo ""
	@echo "  get-migrate-tool       		- Install migration tool"
	@echo ""
	@echo "  redis-up               		- Start Redis container"
	@echo "  redis-down             		- Stop Redis container"
	@echo ""
	@echo "  postgres-migrate-up    		- Run PostgreSQL migrations up"
	@echo "  postgres-migrate-down  		- Run PostgreSQL migrations down"
	@echo "  postgres-up            		- Start PostgreSQL container"
	@echo "  postgres-down          		- Stop PostgreSQL container"
	@echo "  postgres-connection-string 	- Show PostgreSQL connection string"
	@echo ""
	@echo "  clickhouse-up          		- Start ClickHouse container"
	@echo "  clickhouse-down        		- Stop ClickHouse container"
	@echo "  clickhouse-migrate-up  		- Run ClickHouse migrations up"

.PHONY: lint
lint:
	@echo $(GREEN)"[GOLANGCI-LINT] running linter..."$(NO_COLOR)
	@golangci-lint run --config=.golangci.yml --timeout=2m

.PHONY: build
build:
	@echo $(GREEN)"[BUILD] building application..."$(NO_COLOR)
	@go build -v -o $(BINARY_NAME) $(MAIN_PATH)

.PHONY: run
run:
	@echo $(GREEN)"[RUN] running application..."$(NO_COLOR)
	@go run $(MAIN_PATH)

.PHONY: test
test:
	@echo $(GREEN)"[TEST] running tests..."$(NO_COLOR)
	@go test -v ./...

.PHONY: clean
clean:
	@echo $(YELLOW)"[CLEAN] cleaning build artifacts..."$(NO_COLOR)
	@rm -f $(BINARY_NAME)
	@rm -rf ./vendor

.PHONY: tidy
tidy:
	@echo $(GREEN)"[TIDY] tidying go.mod..."$(NO_COLOR)
	@go mod tidy

.PHONY: update-deps
update-deps:
	@echo $(GREEN)"[UPDATE-DEPS] updating dependencies..."$(NO_COLOR)
	@go get -u ./...
	@go mod tidy

.PHONY: docker-run
docker-run:
	@echo $(GREEN)"[DOCKER-RUN] running application in Docker..."$(NO_COLOR)

	@docker build -t $(BINARY_NAME) --build-arg CMD_NAME_ARG=${CMD_NAME_ARG} .

	@if [ -z "$(env)" ]; then \
		echo $(YELLOW)"Using default .env file"$(NO_COLOR); \
		docker run --rm --env-file .env -p 8080:8080 $(BINARY_NAME); \
	else \
		echo $(YELLOW)"Using .env.$(env) file"$(NO_COLOR); \
		docker run --rm --env-file .env.$(env) -p 8080:8080 $(BINARY_NAME); \
	fi

.PHONY: redis-up
redis-up:
	@echo $(GREEN)"[REDIS] starting local container $(REDIS_CONTAINER_NAME)..." $(NO_COLOR)
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^$(REDIS_CONTAINER_NAME)$$"; then \
  		echo $(RED)"[REDIS] container already exists, stopping it..."$(NO_COLOR); \
		make redis-down; \
		echo $(GREEN)"[REDIS] starting local container $(REDIS_CONTAINER_NAME)..." $(NO_COLOR); \
	fi
	@docker run -d --rm --name $(REDIS_CONTAINER_NAME) -p $(REDIS_PORT):6379 $(REDIS_IMAGE) --requirepass $(REDIS_PASSWORD)

.PHONY: redis-down
redis-down:
	@echo $(YELLOW)"[REDIS] stopping local container...$(REDIS_CONTAINER_NAME)..."$(NO_COLOR)
	@docker stop $(REDIS_CONTAINER_NAME)

.PHONY: get-migrate-tool
get-migrate-tool:
	@echo $(GREEN)"[MIGRATE] getting migrate tool..."$(NO_COLOR)
	@go install -tags "$(POSTGRES_MIGRATION_DRIVER) $(CLICKHOUSE_MIGRATION_DRIVER)" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: postgres-migrate-up
postgres-migrate-up:
	@echo $(GREEN)"[POSTGRES] migrate up..."$(NO_COLOR)
	@migrate -path $(POSTGRES_MIGRATION_PATH) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

.PHONY: postgres-migrate-down
postgres-migrate-down:
	@echo $(GREEN)"[POSTGRES] migrate down..."$(NO_COLOR)
	@migrate -path $(POSTGRES_MIGRATION_PATH) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down

.PHONY: postgres-up
postgres-up:
	@echo $(GREEN)"[POSTGRES] starting local container..."$(NO_COLOR)
	@docker run -d --rm --name $(POSTGRES_CONTAINER_NAME) -p $(POSTGRES_PORT):5432 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) $(POSTGRES_IMAGE)

.PHONY: postgres-down
postgres-down:
	@echo $(YELLOW)"[POSTGRES] stopping local container... $(POSTGRES_CONTAINER_NAME)..."$(NO_COLOR)
	@docker stop $(POSTGRES_CONTAINER_NAME)

.PHONY: postgres-connection-string
postgres-connection-string:
	@echo $(YELLOW)"[POSTGRES] connection string: postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"$(NO_COLOR)

.PHONY: clickhouse-up
clickhouse-up:
	@echo $(GREEN)"[CLICKHOUSE] starting local container..."$(NO_COLOR)
	@docker run -d --rm --name $(CLICKHOUSE_CONTAINER_NAME) -p $(CLICKHOUSE_PORT):8123 -p 9000:9000 -p 9999:9999 -e CLICKHOUSE_USER=$(CLICKHOUSE_USER) -e CLICKHOUSE_PASSWORD=$(CLICKHOUSE_PASSWORD) -e CLICKHOUSE_DB=$(CLICKHOUSE_DATABASE) $(CLICKHOUSE_IMAGE)

.PHONY: clickhouse-down
clickhouse-down:
	@echo $(YELLOW)"[CLICKHOUSE] stopping local container..."$(NO_COLOR)
	@docker stop $(CLICKHOUSE_CONTAINER_NAME)

.PHONY: clickhouse-migrate-up
clickhouse-migrate-up:
	@echo $(GREEN)"[CLICKHOUSE] migrate up..."$(NO_COLOR)
	@migrate -path ./migrations/clickhouse -database "clickhouse://$(CLICKHOUSE_HOST):9000?username=$(CLICKHOUSE_USER)&password=$(CLICKHOUSE_PASSWORD)&database=$(CLICKHOUSE_DATABASE)&x-multi-statement=true" up

.PHONY: local-env-up
local-env-up:
	@echo $(GREEN)"[ENVIRONMENT] setting up environment..."$(NO_COLOR)
	@make get-migrate-tool
	@make redis-up
	@make postgres-up
	@sleep 2
	@make postgres-migrate-up
	@make clickhouse-up
	@sleep 5
	@make clickhouse-migrate-up

.PHONY: local-env-down
local-env-down:
	@make redis-down
	@make postgres-down
	@make clickhouse-down
