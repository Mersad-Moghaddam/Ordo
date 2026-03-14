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
