# Ordo (Backend)

Ordo is currently a backend-focused task management API project written in Go.

## Overview

This repository now contains only backend/runtime API infrastructure. All frontend assets and code have been removed.

The backend service provides:

- Authentication flows (register/login/refresh)
- Workspace and project management
- Task management and task status transitions
- Collaboration endpoints (comments and activities)
- Health and Prometheus metrics endpoints

## Core Backend Service

- Entrypoint: `cmd/api/main.go`
- HTTP server bootstrap: `internal/delivery/http/server.go`
- Config loader: `internal/infrastructure/config/config.go`
- Logging: `internal/infrastructure/logging/`
- Security helpers: `internal/infrastructure/security/`
- In-memory persistence adapter (default local runtime): `internal/infrastructure/persistence/memory/`

## Architecture

The backend follows clean architecture patterns under `internal/`:

- `domain/` — entities, value objects, and domain errors
- `usecase/` — application/business workflows
- `repository/` — data-access contracts
- `infrastructure/` — framework/adapter implementations
- `delivery/http/` — Fiber handlers and transport DTOs

## Repository Structure

```text
.
├── cmd/api/                          # API entrypoint
├── internal/                         # backend architecture layers
├── migrations/                       # database migration assets
├── api/                              # phased OpenAPI specs
├── docs/adr/                         # architecture decision records
├── Makefile                          # run/test/lint/benchmark commands
├── go.mod
└── README.md
```

## API Surface

Base path: `/api/v1`

### Auth
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`

### Workspace & projects
- `POST /workspaces`
- `POST /workspaces/:workspaceId/memberships`
- `POST /workspaces/:workspaceId/projects`
- `GET /users/:userId/workspaces`
- `GET /workspaces/:workspaceId/projects`

### Tasks
- `POST /tasks`
- `PATCH /tasks/:taskId/status`
- `GET /projects/:projectId/tasks`

### Collaboration
- `POST /comments`
- `PATCH /comments/:commentId`
- `DELETE /comments/:commentId`
- `GET /tasks/:taskId/comments`
- `GET /tasks/:taskId/activities`

### Platform
- `GET /health`
- `GET /metrics`

## Configuration

Supported environment variables:

- `ORDO_HTTP_PORT` (default: `8080`)
- `ORDO_MYSQL_DSN`
- `ORDO_REDIS_ADDRESS`

## Local Development

### Prerequisites

- Go 1.22+

### Run API

```bash
go run ./cmd/api
```

Default address: `http://127.0.0.1:8080`

## Development Commands

```bash
make run         # run backend API
make test        # go test ./... with coverage profile
make benchmark   # benchmark platform usecases
make revive      # lint via revive
make migrate     # run DB migrations (requires ORDO_MYSQL_MIGRATE_DSN)
make sqlc        # generate sqlc artifacts
make tidy        # go mod tidy
```

## Quality Checks

Before merging:

```bash
go test ./...
```

## Notes

- OpenAPI specs are in `api/`.
- ADRs are in `docs/adr/`.
- Migration files are in `migrations/`.
- No frontend application is included in this repository at this time.
