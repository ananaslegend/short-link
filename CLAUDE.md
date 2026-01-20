# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a microservices-based short link platform written in Go using a monorepo structure. The platform consists of multiple services:

- **Short-Link Service** (`backend/services/short-link/`) - URL shortening functionality
- **Statistic Service** (`backend/services/statistic/`) - Analytics and usage statistics

**Tech Stack:**
- **Backend:** Go with Echo web framework
- **Architecture:** Uber-FX dependency injection, modular design
- **Databases:** PostgreSQL (links), Redis (caching), ClickHouse (analytics)
- **Observability:** OpenTelemetry, Prometheus, Grafana, Loki, Alloy, Tempo
- **Containerization:** Docker, Docker Compose with Helm charts

## Common Development Commands

### Short-Link Service
Commands should be run from `backend/services/short-link/` directory:

```bash
# Build and run
make build
make run

# Run with Docker Compose (full stack)
make docker-compose-up
make docker-compose-down

# Database operations
make migrate-up              # PostgreSQL migrations
make migrate-up-clickhouse   # ClickHouse migrations

# Code quality
make lint                    # Run golangci-lint
make test                    # Run tests
make tidy                    # Tidy dependencies

# Docker operations
make docker-run              # Run in Docker with default .env
make docker-run env=dev      # Run with specific environment file
```

### Statistic Service
Commands should be run from `backend/services/statistic/` directory:

```bash
# Build and run
make build
make run

# Database operations
make migrate-up-clickhouse   # ClickHouse migrations only
make migrate-down-clickhouse # Rollback ClickHouse migrations

# Code quality
make lint                    # Run golangci-lint
make test                    # Run tests
make tidy                    # Tidy dependencies

# Docker operations
make docker-run              # Run in Docker (port 8081)
make docker-run env=dev      # Run with specific environment file
```

## Project Architecture

### Directory Structure

#### Short-Link Service (`backend/services/short-link/`)
- `cmd/api/` - Application entry point
- `internal/app/` - Application framework (Echo wrapper, Redis wrapper, middleware)
- `internal/link/` - Core link shortening business logic and repository
- `internal/alias_generator/` - URL alias generation logic
- `migrations/` - Database migrations (PostgreSQL and ClickHouse)

#### Statistic Service (`backend/services/statistic/`)
- `cmd/api/` - Application entry point
- `internal/app/` - Application framework (Echo wrapper, ClickHouse wrapper, middleware)
- `internal/statistic/` - Analytics and statistics logic and repository
  - `domain/` - Domain models for statistics
  - `repository/` - ClickHouse data access layer
  - `service/` - Business logic for analytics
- `migrations/clickhouse/` - ClickHouse migrations for analytics

### Key Architectural Patterns
- **Microservices:** Independent services with clear separation of concerns
- **Dependency Injection:** Uses uber-go/fx for modular architecture and dependency management
- **Repository Pattern:** Data access abstracted through repository interfaces
- **Middleware:** Request logging and observability middleware via Echo
- **Configuration:** Viper-based configuration with environment variable support
- **Observability:** Comprehensive OpenTelemetry instrumentation for tracing and metrics

### Service Communication
- **Short-Link Service:** Runs on port 8080, handles URL shortening and redirects
- **Statistic Service:** Runs on port 8081, handles analytics and usage statistics
- Services can communicate via nats messaging 

### Configuration
Each service has its own configuration:
- **Short-Link Service:** PostgreSQL, Redis, ClickHouse, full observability stack
- **Statistic Service:** ClickHouse only, focused observability for analytics

Environment files: `.env` (default), `.env.example` (template) in each service directory

### Database Schema
- **PostgreSQL:** (Short-Link Service) Stores link mappings and metadata
- **ClickHouse:** (Both Services) Stores analytics and usage statistics
- **Redis:** (Short-Link Service) Caching layer for frequently accessed links

### Observability Stack
Full observability setup via Docker Compose:
- **OpenTelemetry Collector:** Collects traces and metrics
- **Prometheus:** Metrics storage and querying
- **Grafana:** Visualization dashboards
- **Loki:** Log aggregation
- **Alloy:** Log collection and forwarding
- **Tempo:** Distributed tracing storage

### Linting Configuration
Both services use golangci-lint with custom configuration in `.golangci.yml`:
- Comprehensive linter set with specific exclusions
- Auto-formatting enabled (gci, gofmt, gofumpt, goimports, golines)
- Test file exemptions for certain linters

### Development Notes
- **Service Independence:** Each service can be developed, built, and deployed independently
- **Shared Dependencies:** Both services use similar tech stack but with service-specific dependencies
- **Port Allocation:** Short-Link (8080), Statistic (8081)
- **Handler Implementation:** Statistic service handlers are not yet implemented (placeholder API module)
- migrations and postgres data schema for short-link db stored in @backend/services/short-link/migrations/postgres/
- migrations and clickhouse data schema for statistic db stored in @backend/services/short-link/migrations/clickhouse/
- If you create new files, add them to git (except build binaries)