# Ordo

Ordo is a split full-stack monorepo with:

- **`backend/`**: Go API service (Fiber) following clean architecture (delivery/usecase/repository/infrastructure), with in-memory adapters for local development and tests.
- **`frontend/`**: React + TypeScript + Vite client for running core product flows against the API.

## Repository Structure

```text
.
├── backend/
│   ├── cmd/api/                # API entrypoint
│   ├── internal/               # domain, usecases, delivery, infrastructure, repository ports
│   ├── migrations/             # SQL schema migrations
│   ├── api/                    # phased OpenAPI specs
│   ├── docs/adr/               # architecture decision records
│   └── Makefile                # backend dev/test/lint/benchmark commands
└── frontend/
    ├── src/                    # React app source
    ├── package.json            # frontend scripts and dependencies
    └── vite.config.ts          # Vite configuration
```


## Repository Hygiene

- Backend runtime/source files live **only** under `backend/`.
- Frontend runtime/source files live **only** under `frontend/`.
- The repository root is intentionally kept minimal (workspace-level docs only).

## Prerequisites

- **Go** 1.22+
- **Node.js** 20+
- **npm** 10+

## Quick Start

### 1) Run Backend

```bash
cd backend
go run ./cmd/api
```

Backend starts on `http://127.0.0.1:8080` by default.

### 2) Run Frontend

```bash
cd frontend
npm ci
npm run dev
```

Frontend starts on Vite's dev server (usually `http://127.0.0.1:5173`).

## Backend Commands

From `backend/`:

```bash
make run         # run API
make test        # go test ./...
make benchmark   # benchmark suite
make revive      # lint via revive
make tidy        # go mod tidy
```

## Frontend Commands

From `frontend/`:

```bash
npm run dev      # start dev server
npm run build    # production build
npm run preview  # preview built app
```

## API Highlights

- `POST /auth/register`
- `POST /auth/login`
- `GET /workspaces`
- `POST /tasks`
- `GET /health`
- `GET /metrics`

Use the frontend to exercise end-to-end flows (register → login → workspace/task actions).

## Quality and CI

Backend CI workflow runs tests, linting, and benchmarks. Local validation is recommended before PRs:

```bash
cd backend && go test ./...
cd frontend && npm run build
```

## Notes

- Generated/build artifacts are intentionally excluded from version control (`frontend/dist`, `frontend/node_modules`, `*.tsbuildinfo`).
- The repository is now cleanly split between backend and frontend concerns.
