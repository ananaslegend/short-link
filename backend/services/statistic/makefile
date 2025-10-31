BINARY_NAME=statistic
MAIN_PATH=./cmd/api
CMD_NAME_ARG=api

CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=8123
CLICKHOUSE_USER=sl
CLICKHOUSE_PASSWORD=sl_password
CLICKHOUSE_DATABASE=statistics
CLICKHOUSE_IMAGE=clickhouse/clickhouse-server:24.4-alpine
CLICKHOUSE_CONTAINER_NAME := clickhouse-server
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
	@echo "  lint                	- Run golangci-lint"
	@echo "  build               	- Build the application"
	@echo "  run                 	- Run the application locally"
	@echo "  test                	- Run tests"
	@echo "  clean               	- Clean build artifacts"
	@echo "  tidy                	- Tidy go.mod"
	@echo "  update-deps         	- Update dependencies"
	@echo "  docker-run          	- Run application in Docker using default .env file"
	@echo "  docker-run env={env}	- Run application in Docker using .env.{env} file"
	@echo "                          Example: make docker-run env=dev"
	@echo "  migrate-up-clickhouse 	- Run ClickHouse migrations up"
	@echo "  migrate-down-clickhouse	- Run ClickHouse migrations down"

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
	@go run $(MAIN_PATH) > tmp/app.log 2>&1

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
		docker run --rm --env-file .env -p 8081:8081 $(BINARY_NAME); \
	else \
		echo $(YELLOW)"Using .env.$(env) file"$(NO_COLOR); \
		docker run --rm --env-file .env.$(env) -p 8081:8081 $(BINARY_NAME); \
	fi

.PHONY: migrate-up-clickhouse
migrate-up-clickhouse:
	@echo $(GREEN)"[MIGRATE] running ClickHouse migrations up..."$(NO_COLOR)
	@migrate -database "$(CLICKHOUSE_MIGRATION_DRIVER)://$(CLICKHOUSE_USER):$(CLICKHOUSE_PASSWORD)@$(CLICKHOUSE_HOST):$(CLICKHOUSE_PORT)/$(CLICKHOUSE_DATABASE)" -path $(CLICKHOUSE_MIGRATION_PATH) up

.PHONY: migrate-down-clickhouse
migrate-down-clickhouse:
	@echo $(RED)"[MIGRATE] running ClickHouse migrations down..."$(NO_COLOR)
	@migrate -database "$(CLICKHOUSE_MIGRATION_DRIVER)://$(CLICKHOUSE_USER):$(CLICKHOUSE_PASSWORD)@$(CLICKHOUSE_HOST):$(CLICKHOUSE_PORT)/$(CLICKHOUSE_DATABASE)" -path $(CLICKHOUSE_MIGRATION_PATH) down