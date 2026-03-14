## Summary
Implemented Phase 1 authentication and user management foundations with JWT-like signed tokens, refresh rotation, RBAC authorization primitives, and persistence schema/query scaffolding.

## Key Changes
- Added auth domain models and typed errors for user/session/token workflows.
- Added auth repository interfaces for users and refresh sessions.
- Added security utilities for password hashing and signed token issue/verify operations with functional options.
- Added auth usecase service implementing register, login, refresh rotation, and role authorization.
- Added Fiber auth delivery components (handlers, DTOs, middleware) for register/login/refresh and role checks.
- Added migration `0002_auth_sessions` and SQLC queries for users and refresh sessions.
- Added ADR `0002-auth-token-rotation-and-rbac.md` and OpenAPI phase1 snippet.
- Added unit tests for auth usecase and security token/hash utilities.

## Validation
- `go test ./internal/usecase/auth ./internal/infrastructure/security ./internal/infrastructure/config -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
Bootstrap Phase 0 for Ordo with clean architecture scaffolding, strict revive linting, functional configuration patterns, baseline migrations, and CI checks.

## Key Changes
- Added Fiber API entrypoint with `/health` and `/metrics` endpoints.
- Added strict `revive.toml` rules and Makefile targets (`run`, `test`, `revive`, `migrate`, `sqlc`, `tidy`).
- Added idempotent MySQL migrations for `users` and `outbox_events`.
- Added SQLC configuration and initial outbox insert query.
- Added ADR 0001 documenting architectural choices.
- Added CI workflow for tests and linting.

## Validation
- `go test ./... -coverprofile=coverage.out`
- `go test ./internal/infrastructure/config ./internal/infrastructure/persistence/mysql -cover`
- `make revive`
