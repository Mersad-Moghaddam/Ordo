## Summary
Implemented Phase 5 worker foundations for Redis Streams consumption with idempotency, retry backoff, and DLQ routing.

## Key Changes
- Added worker usecase service for stream polling, idempotency checks, notification dispatch, retry scheduling, and DLQ publishing.
- Added in-memory Redis stream consumer adapter for local testing and wiring.
- Added migration `0006_worker_idempotency` and SQLC query file for processed event keys.
- Added ADR `0006-worker-idempotency-retry-dlq.md` and OpenAPI phase5 snippet.
- Added unit tests for worker success/retry/DLQ paths and stream consumer behavior.

## Validation
- `go test ./internal/usecase/worker ./internal/infrastructure/broker/redis/consumer -cover`
- `go test ./internal/usecase/auth ./internal/usecase/workspace ./internal/usecase/task ./internal/usecase/collab ./internal/usecase/worker ./internal/infrastructure/security ./internal/infrastructure/config ./internal/infrastructure/broker/redis/producer ./internal/infrastructure/broker/redis/consumer -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
Implemented Phase 4 comments and activity logging with transactional event persistence to outbox.

## Key Changes
- Added collab domain models and typed errors for comments and activity logs.
- Added collab repository interfaces for comments and activity log persistence contracts.
- Added collab usecase service implementing create/update/delete comment flows with transactional activity and outbox persistence.
- Added Fiber delivery DTO/handler for comments and task activity endpoints.
- Added migration `0005_comments_activity` and SQLC query file for comments and activity logs.
- Added ADR `0005-comment-activity-event-flow.md` and OpenAPI phase4 snippet.
- Added unit tests for comment lifecycle, authorization checks, and listing behavior.

## Validation
- `go test ./internal/usecase/collab -cover`
- `go test ./internal/usecase/auth ./internal/usecase/workspace ./internal/usecase/task ./internal/usecase/collab ./internal/infrastructure/security ./internal/infrastructure/config ./internal/infrastructure/broker/redis/producer -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
Implemented Phase 3 task core with transactional outbox writes and Redis Streams producer scaffolding.

## Key Changes
- Added task domain models and typed task/outbox errors.
- Added task and outbox repository interfaces.
- Added task usecase service for create, status update, and paginated list operations with transactional outbox persistence.
- Added Redis Streams outbox producer with retry bookkeeping and exponential backoff strategy.
- Added Fiber delivery DTO/handler for task endpoints.
- Added migration `0004_tasks_core` and SQLC task/outbox query file.
- Added ADR `0004-task-outbox-and-stream-producer.md` and OpenAPI phase3 snippet.
- Added unit tests for task usecase and outbox producer behavior.

## Validation
- `go test ./internal/usecase/task ./internal/infrastructure/broker/redis/producer -cover`
- `go test ./internal/usecase/auth ./internal/usecase/workspace ./internal/usecase/task ./internal/infrastructure/security ./internal/infrastructure/config ./internal/infrastructure/broker/redis/producer -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
Implemented Phase 2 workspace, membership, and project foundations with clean architecture boundaries, role-aware orchestration, pagination, migrations, and SQL query scaffolding.

## Key Changes
- Added workspace domain entities and typed domain errors.
- Added repository interfaces for workspace, membership, and project persistence contracts.
- Added workspace usecase service implementing workspace creation, owner bootstrap membership, membership management, project creation, role checks, and paginated listings.
- Added Fiber delivery handler/DTOs for workspace and project endpoints.
- Added migration `0003_workspaces_projects` and SQLC query files for workspaces, memberships, and projects.
- Added ADR `0003-workspace-membership-project-authorization.md` and OpenAPI phase2 snippet.
- Added table-driven tests for workspace lifecycle and authorization behavior.

## Validation
- `go test ./internal/usecase/workspace -cover`
- `go test ./internal/usecase/auth ./internal/infrastructure/security ./internal/infrastructure/config -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
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
