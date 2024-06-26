REDIS_PORT=6379
REDIS_IMAGE=redis:7.2.5-alpine
REDIS_CONTAINER_NAME=sl-redis

DB_DRIVER=postgres
MIGRATE_PATH=./migrations

POSTGRES_PORT=5432
POSTGRES_IMAGE=postgres:16.3-alpine
POSTGRES_CONTAINER_NAME=sl-postgres
POSTGRES_USER=sl
POSTGRES_PASSWORD=sl_password
POSTGRES_DB=sl

GREEN = "\033[0;32m"
YELLOW = "\033[0;33m"
RED = "\033[0;31m"
NO_COLOR = "\033[0m"

.PHONY: redis_up
redis_up:
	@echo $(GREEN)"[REDIS] starting local container $(REDIS_CONTAINER_NAME)..." $(NO_COLOR)
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^$(REDIS_CONTAINER_NAME)$$"; then \
  		echo $(RED)"[REDIS] container already exists, stopping it..."$(NO_COLOR); \
		make local-redis_down; \
		echo $(GREEN)"[REDIS] starting local container $(REDIS_CONTAINER_NAME)..." $(NO_COLOR); \
	fi
	@docker run -d --rm --name $(REDIS_CONTAINER_NAME) -p $(REDIS_PORT):6379 $(REDIS_IMAGE)

.PHONY: redis_down
redis_down:
	@echo $(YELLOW)"[REDIS] stopping local container $(REDIS_CONTAINER_NAME)..."$(NO_COLOR)
	@docker stop $(REDIS_CONTAINER_NAME)

.PHONY: get-migrate
get-migrate:
	@echo $(GREEN)"[MIGRATE] getting migrate tool..."$(NO_COLOR)
	@go install -tags $(DB_DRIVER) github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: postgres-migrate-up
postgres-migrate-up:
	@echo $(GREEN)"[POSTGRES] migrate up..."$(NO_COLOR)
	@migrate -path ./migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

.PHONY: postgres-migrate-down
postgres-migrate-down:
	@echo $(GREEN)"[POSTGRES] migrate down..."$(NO_COLOR)
	@migrate -path ./migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down

.PHONY: postgres_up
postgres_up:
	@echo $(GREEN)"[POSTGRES] starting local container..."$(NO_COLOR)
	@docker run -d --rm --name $(POSTGRES_CONTAINER_NAME) -p $(POSTGRES_PORT):5432 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) $(POSTGRES_IMAGE)

.PHONY: postgres_down
postgres_down:
	@echo $(YELLOW)"[POSTGRES] stopping local container $(POSTGRES_CONTAINER_NAME)..."$(NO_COLOR)
	@docker stop $(POSTGRES_CONTAINER_NAME)

.PHONY: postgres_connection_string
postgres_connection_string:
	@echo $(YELLOW)"[POSTGRES] connection string: postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"$(NO_COLOR)