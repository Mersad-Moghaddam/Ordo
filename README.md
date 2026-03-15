# Ordo

Ordo is a full-stack task management platform in a monorepo, built with a Go API and a React + TypeScript frontend.

It is designed around clean architecture on the backend and a modern, Jira-style planning experience on the frontend.

---

## Table of Contents

- [1. What this project includes](#1-what-this-project-includes)
- [2. Core services and modules](#2-core-services-and-modules)
- [3. Architecture overview](#3-architecture-overview)
- [4. Repository structure](#4-repository-structure)
- [5. API domains](#5-api-domains)
- [6. Configuration](#6-configuration)
- [7. Local development](#7-local-development)
- [8. Testing, quality, and validation](#8-testing-quality-and-validation)
- [9. Data, migrations, and specs](#9-data-migrations-and-specs)
- [10. Frontend feature overview](#10-frontend-feature-overview)
- [11. Contribution guide](#11-contribution-guide)
- [12. Troubleshooting](#12-troubleshooting)
- [13. License](#13-license)

---

## 1. What this project includes

Ordo currently contains two main runtime services:

1. **Backend API service** (`backend/`)
   - Go + Fiber HTTP API
   - Health and Prometheus metrics endpoints
   - Auth, workspace/project, task, and collaboration domains

2. **Frontend web app** (`frontend/`)
   - React + TypeScript + Vite
   - Jira-style task board/list UI
   - Local-first task interactions for fast iteration and UX testing

---

## 2. Core services and modules

### Backend service

- **Entrypoint:** `backend/cmd/api/main.go`
- **HTTP server bootstrap:** `backend/internal/delivery/http/server.go`
- **Configuration:** `backend/internal/infrastructure/config/config.go`
- **Logging:** Zap logger in `backend/internal/infrastructure/logging/`
- **Persistence (default runtime):** in-memory store (`backend/internal/infrastructure/persistence/memory/`)
- **Security primitives:** token + password services (`backend/internal/infrastructure/security/`)

### Frontend app

- **Entrypoint:** `frontend/src/main.tsx`
- **Main UI container:** `frontend/src/App.tsx`
- **Task state module:** `frontend/src/hooks/useTasks.ts`
- **Shared types:** `frontend/src/types.ts`
- **Styling:** `frontend/src/styles.css`

---

## 3. Architecture overview

### Backend architecture (clean architecture)

`backend/internal/` is organized by responsibility:

- `domain/` — business entities and domain errors
- `usecase/` — application/business workflows
- `repository/` — data access contracts
- `infrastructure/` — concrete adapters (config, security, persistence, broker/cache)
- `delivery/http/` — HTTP transport handlers and DTO mapping

This design keeps business logic independent from framework and storage details.

### Frontend architecture

- **UI layer** in `App.tsx` handles rendering and interactions.
- **State + behavior layer** in `useTasks.ts` centralizes CRUD, filtering, grouping, persistence, and metrics.
- **Type layer** in `types.ts` keeps task model contracts explicit.

---

## 4. Repository structure

```text
.
├── backend/
│   ├── cmd/api/                      # Backend entrypoint
│   ├── internal/                     # Domain/usecase/repository/infrastructure/delivery
│   ├── migrations/                   # DB migrations (backend-owned)
│   ├── api/                          # Phased OpenAPI specs
│   ├── docs/adr/                     # Architecture Decision Records
│   ├── Makefile                      # Backend run/test/lint/benchmark
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── App.tsx
│   │   ├── hooks/useTasks.ts
│   │   ├── styles.css
│   │   └── types.ts
│   ├── package.json
│   └── vite.config.ts
├── api/                              # Additional/root-level API specs
├── docs/adr/                         # Additional/root-level ADRs
├── migrations/                       # Additional/root-level migrations
└── README.md
```

---

## 5. API domains

Base API prefix: `/api/v1`

### Auth
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`

### Workspace & project management
- `POST /workspaces`
- `POST /workspaces/:workspaceId/memberships`
- `POST /workspaces/:workspaceId/projects`
- `GET /users/:userId/workspaces`
- `GET /workspaces/:workspaceId/projects`

### Task management
- Task creation/update/list and status transitions are exposed under task routes in the task delivery module.

### Collaboration
- `POST /comments`
- `PATCH /comments/:commentId`
- `DELETE /comments/:commentId`
- `GET /tasks/:taskId/comments`
- `GET /tasks/:taskId/activities`

### Platform endpoints
- `GET /health` — liveness check
- `GET /metrics` — Prometheus metrics

For contract details, see phased specs in `backend/api/` (and root `api/` where applicable).

---

## 6. Configuration

Environment variables currently supported by backend configuration:

- `ORDO_HTTP_PORT` (default: `8080`)
- `ORDO_MYSQL_DSN` (default set in code)
- `ORDO_REDIS_ADDRESS` (default: `localhost:6379`)

The default runtime path uses in-memory adapters, but project structure supports MySQL/Redis-oriented infrastructure and migration/spec assets.

---

## 7. Local development

## Prerequisites

- Go 1.22+
- Node.js 20+
- npm 10+

### Start backend

```bash
cd backend
go run ./cmd/api
```

Backend default URL: `http://127.0.0.1:8080`

### Start frontend

```bash
cd frontend
npm ci
npm run dev
```

Frontend default URL: `http://127.0.0.1:5173`

---

## 8. Testing, quality, and validation

### Backend

```bash
cd backend
make test
make benchmark
make revive
```

### Frontend

```bash
cd frontend
npm run build
```

Recommended pre-PR baseline:

```bash
cd backend && go test ./...
cd frontend && npm run build
```

---

## 9. Data, migrations, and specs

### Migrations

- Backend-owned migrations: `backend/migrations/`
- Additional/root migration assets: `migrations/`

### OpenAPI specs

- Backend phased specs: `backend/api/openapi.phase*.yaml`
- Additional/root specs: `api/openapi.phase*.yaml`

### ADRs

- Backend ADRs: `backend/docs/adr/`
- Additional/root ADRs: `docs/adr/`

---

## 10. Frontend feature overview

Current frontend UX includes:

- Jira-style board and list toggles
- Task creation form (title, description, assignee, due date, priority, status)
- Search and filters (query, status, priority)
- Status movement controls and task deletion
- Sprint completion KPI and summary cards
- LocalStorage-backed persistence for rapid local iteration

---

## 11. Contribution guide

- Keep backend runtime code under `backend/`.
- Keep frontend runtime code under `frontend/`.
- Prefer small, focused commits.
- Validate changes with local build/test commands before PR.
- Avoid committing generated artifacts (`dist/`, `node_modules/`, `*.tsbuildinfo`, etc.).

---

## 12. Troubleshooting

### Frontend build fails

- Ensure Node/npm versions meet prerequisites.
- Reinstall dependencies:

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

### Backend startup issues

- Verify Go toolchain version.
- Check port availability for `ORDO_HTTP_PORT`.
- Confirm environment variable values are valid (especially integer parsing for port).

### No data persistence in backend

- Current default runtime uses in-memory storage for fast local execution.
- For DB-backed flows, use migration/spec artifacts and infrastructure components as a foundation.

---

## 13. License

No license file is currently defined in this repository. Add one before public distribution.
