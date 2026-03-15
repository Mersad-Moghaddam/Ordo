# Ordo

Ordo is a full-stack task management monorepo composed of a Go backend and a React + TypeScript frontend. The project is structured to keep backend and frontend concerns cleanly separated while supporting fast local development.

## Highlights

- Clean architecture backend (`delivery`, `usecase`, `repository`, `infrastructure`)
- Jira-style task workspace frontend (board + list modes)
- Monorepo-friendly workflows for build, test, lint, and benchmarking
- ADRs and phased OpenAPI specs for architecture and API evolution

## Technology

### Backend (`backend/`)
- Go
- Fiber
- SQL migrations
- OpenAPI phase specs

### Frontend (`frontend/`)
- React 18
- TypeScript
- Vite

## Repository Layout

```text
.
├── backend/
│   ├── cmd/api/                # API entrypoint
│   ├── internal/               # Domain/usecase/delivery/infrastructure
│   ├── migrations/             # Backend schema migrations
│   ├── api/                    # OpenAPI specs
│   ├── docs/adr/               # Architecture decision records
│   └── Makefile                # Backend development commands
├── frontend/
│   ├── src/                    # React application source
│   ├── package.json            # npm scripts + dependencies
│   └── vite.config.ts          # Vite config
└── README.md
```

## Prerequisites

- Go 1.22+
- Node.js 20+
- npm 10+

## Quick Start

### 1) Run backend

```bash
cd backend
go run ./cmd/api
```

Backend default: `http://127.0.0.1:8080`

### 2) Run frontend

```bash
cd frontend
npm ci
npm run dev
```

Frontend default: `http://127.0.0.1:5173`

## Development Commands

### Backend

```bash
cd backend
make run         # start API
make test        # go test ./...
make benchmark   # benchmark suite
make revive      # lint
make tidy        # go mod tidy
```

### Frontend

```bash
cd frontend
npm run dev      # start dev server
npm run build    # type-check + build
npm run preview  # preview production build
```

## Frontend Features

- Task creation, deletion, status updates
- Board and list views
- Filters by status and priority
- Search by title/description/assignee
- Local persistence for iterative prototyping

## API Highlights

- `POST /auth/register`
- `POST /auth/login`
- `GET /workspaces`
- `POST /tasks`
- `GET /health`
- `GET /metrics`

OpenAPI evolution lives under `backend/api/`.

## Quality Checks

Run before creating a PR:

```bash
cd backend && go test ./...
cd frontend && npm run build
```

## Contribution Guidelines

- Keep backend runtime/source files inside `backend/`.
- Keep frontend runtime/source files inside `frontend/`.
- Do not commit generated artifacts (`dist/`, `node_modules/`, `*.tsbuildinfo`).
- Prefer focused, atomic commits.

## License

No license file is currently defined. Add one before public distribution.
