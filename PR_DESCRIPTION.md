## Summary
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
