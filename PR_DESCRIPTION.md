## Summary
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
