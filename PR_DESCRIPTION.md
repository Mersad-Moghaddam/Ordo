## Summary
Implemented Phase 7 final polishing with CI hardening, benchmark coverage, and release-readiness artifacts.

## Key Changes
- Added benchmark suite for platform cache and rate-limiter paths.
- Added `make benchmark` target.
- Hardened GitHub Actions workflow to install revive on PATH and run tests, linting, and benchmarks.
- Added ADR `0008-final-polish-observability-ci-benchmarks.md`.
- Added OpenAPI phase7 snippet consolidating API and worker health/metrics surfaces.

## Validation
- `go test ./internal/usecase/platform -run Test -cover`
- `go test ./internal/usecase/platform -bench=. -run=^$ -benchmem`
- `go test ./internal/usecase/auth ./internal/usecase/workspace ./internal/usecase/task ./internal/usecase/collab ./internal/usecase/worker ./internal/usecase/platform ./internal/infrastructure/security ./internal/infrastructure/config ./internal/infrastructure/broker/redis/producer ./internal/infrastructure/broker/redis/consumer ./internal/infrastructure/cache/redis -cover`
- `go test ./...` (blocked by environment dependency download restrictions)
- `make revive` (blocked by missing revive binary in environment)
