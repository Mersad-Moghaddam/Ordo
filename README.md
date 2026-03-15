# Ordo

Ordo is a full-stack task management monorepo with a Go backend and a modern React frontend. It is structured for maintainability, iterative delivery, and clear separation between API and UI concerns.

## Why Ordo

- **Clean architecture backend** for long-term maintainability.
- **Modern TypeScript frontend** with an interactive Jira-style task board/list experience.
- **Monorepo layout** with isolated backend/frontend workflows and tooling.
- **Developer-friendly local setup** with focused commands for run, test, lint, and build.

## Tech Stack

### Backend (`backend/`)
- Go
- Fiber (HTTP server)
- Clean architecture layers (`delivery`, `usecase`, `repository`, `infrastructure`)
- SQL migrations and phased OpenAPI specs

### Frontend (`frontend/`)
- React 18
- TypeScript
- Vite
- Local-first task workspace UI (board + list views)

## Repository Structure

```text
.
├── backend/
│   ├── cmd/api/                # API entrypoint
│   ├── internal/               # Domain + use cases + delivery + infrastructure
│   ├── migrations/             # SQL schema migrations
│   ├── api/                    # OpenAPI specs by phase
│   ├── docs/adr/               # Architecture Decision Records
│   └── Makefile                # Backend workflows (run/test/lint/bench)
├── frontend/
│   ├── src/                    # React app source
│   ├── package.json            # Frontend scripts/dependencies
│   └── vite.config.ts          # Vite configuration
└── README.md
```

## Prerequisites

- **Go** 1.22+
- **Node.js** 20+
- **npm** 10+

## Quick Start

### 1) Start the Backend

```bash
cd backend
go run ./cmd/api
```

Backend defaults to `http://127.0.0.1:8080`.

### 2) Start the Frontend

```bash
cd frontend
npm ci
npm run dev
```

Frontend defaults to `http://127.0.0.1:5173`.

---

## Development Workflows

### Backend Commands

Run from `backend/`:

```bash
make run         # Run API server
make test        # Execute go test ./...
make benchmark   # Run benchmark suite
make revive      # Lint using revive
make tidy        # go mod tidy
```

### Frontend Commands

Run from `frontend/`:

```bash
npm run dev      # Start Vite dev server
npm run build    # Type-check + production build
npm run preview  # Preview production build
```

## Frontend Capabilities

The current frontend provides a polished, user-friendly task workspace with:

- Board and list view modes
- Task creation and deletion
- Status updates and quick move actions
- Search and filter controls
- Priority and due-date visual indicators
- Local persistence for rapid prototyping

## API Surface (Highlights)

Representative backend endpoints include:

- `POST /auth/register`
- `POST /auth/login`
- `GET /workspaces`
- `POST /tasks`
- `GET /health`
- `GET /metrics`

See phased specs under `backend/api/` for endpoint evolution and details.

## Quality Expectations

Before opening a PR, run:

```bash
cd backend && go test ./...
cd frontend && npm run build
```

## Contribution Notes

- Keep backend runtime/source files inside `backend/`.
- Keep frontend runtime/source files inside `frontend/`.
- Avoid committing generated artifacts (`dist/`, `node_modules/`, `*.tsbuildinfo`).
- Prefer small, focused commits with clear messages.

## License

No license file is currently defined in this repository. Add one before public distribution.
