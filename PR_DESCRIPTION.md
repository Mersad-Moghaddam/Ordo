## Summary
Implemented Phase 6 platform hardening foundations for caching and rate limiting.

## Key Changes
- Added generic read-through cache usecase utility with typed payload support.
- Added rate limiter usecase utility with subject-operation window enforcement.
- Added in-memory Redis-compatible cache and rate window store adapter.
- Added migration `0007_rate_limit_policies` and SQLC query file for policy persistence metadata.
- Added ADR `0007-caching-and-rate-limit-hardening.md` and OpenAPI phase6 snippet.
- Added unit tests for cache hit/miss behavior, fetch failures, and rate-limit enforcement.

## Validation
- `go test ./internal/usecase/platform ./internal/infrastructure/cache/redis -cover`
- `go test ./internal/usecase/auth ./internal/usecase/workspace ./internal/usecase/task ./internal/usecase/collab ./internal/usecase/worker ./internal/usecase/platform ./internal/infrastructure/security ./internal/infrastructure/config ./internal/infrastructure/broker/redis/producer ./internal/infrastructure/broker/redis/consumer ./internal/infrastructure/cache/redis -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
