# ADR 0001: Clean Architecture Bootstrap and Event Backbone

## Status
Accepted

## Context
Ordo requires a maintainable backend that supports collaborative task workflows, strict quality gates, and progressive delivery across phased implementation. The baseline must enable transactional consistency for write operations and event-driven expansion without introducing heavy broker dependencies.

## Decision
Adopt a clean architecture layout with these top-level layers:
- `internal/delivery` for Fiber HTTP adapters.
- `internal/usecase` for application orchestration.
- `internal/domain` for enterprise entities and typed domain errors.
- `internal/repository` for abstraction contracts and transaction management boundaries.
- `internal/infrastructure` for concrete adapters (MySQL, Redis, logging, metrics).

Adopt SQL-first persistence with `sqlc` generation configured in `sqlc.yaml`.

Adopt Redis Streams as event transport with reserved names:
- stream: `ordo_events`
- consumer group: `ordo_workers`
- dead-letter queue: `ordo_events_dlq`

Adopt functional options for runtime configuration and logger initialization.

Enforce code quality through `revive` with strict rules, Makefile targets, and GitHub Actions checks.

## Consequences
Positive outcomes:
- Layer isolation reduces coupling and supports testability.
- `sqlc` ensures compile-time query typing for MySQL operations.
- Event naming standardization allows future outbox-to-stream workers with predictable contracts.
- Functional options provide extensible configuration without constructor breakage.

Tradeoffs:
- Additional boilerplate in early phases.
- Strict linting may increase development effort.

## Alternatives Considered
- Monolithic handler-service-repository layering with direct infra imports in domain was rejected due to boundary leakage.
- ORM-based data access was rejected to preserve SQL control and explicit performance tuning.
- External broker adoption was deferred to reduce operational complexity.
