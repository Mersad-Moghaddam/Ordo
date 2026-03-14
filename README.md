# Ordo Monorepo

This repository is organized into two top-level applications:

- `backend/` — Go backend API (Fiber + MySQL/sqlc + Redis Streams architecture)
- `frontend/` — Modern React + TypeScript + Vite web client

## Quick Start

### Backend

```bash
cd backend
go run ./cmd/api
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Set backend URL from the UI (defaults to `http://127.0.0.1:8080`) and run the built-in full endpoint flow.
