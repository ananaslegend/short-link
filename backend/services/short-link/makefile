BINARY_NAME=app
MAIN_PATH=./cmd/api
CMD_NAME_ARG=api

REDIS_PORT=6379
REDIS_IMAGE=redis:7.2.5-alpine
REDIS_CONTAINER_NAME=redis
REDIS_PASSWORD=pass

POSTGRES_PORT=5432
POSTGRES_IMAGE=postgres:16.3-alpine
POSTGRES_CONTAINER_NAME=postgres
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
		docker run --rm --env-file .env -p 8080:8080 $(BINARY_NAME); \
	else \
		echo $(YELLOW)"Using .env.$(env) file"$(NO_COLOR); \
		docker run --rm --env-file .env.$(env) -p 8080:8080 $(BINARY_NAME); \
	fi
