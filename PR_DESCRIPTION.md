## Summary
Refined Phase 0 bootstrap by tightening runtime lifecycle behavior and test naming standards while preserving clean architecture boundaries.

## Key Changes
- Added explicit HTTP server shutdown support with context-aware termination.
- Updated API bootstrap to perform graceful shutdown through errgroup coordination and non-panic process exit.
- Normalized table-driven test function identifier names to comply with strict descriptive naming conventions.
- Preserved baseline architecture, migrations, sqlc config, lint rules, and CI flow from the initial bootstrap.

## Validation
- `go test ./internal/infrastructure/config -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
