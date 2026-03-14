## Summary
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
