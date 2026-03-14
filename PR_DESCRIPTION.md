## Summary
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
