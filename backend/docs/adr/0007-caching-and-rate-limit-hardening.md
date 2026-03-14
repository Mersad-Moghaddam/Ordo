# ADR 0007: Caching and Rate-Limit Hardening

## Status
Accepted

## Context
Phase 6 requires platform hardening to reduce backend load and protect command endpoints from abuse patterns.

## Decision
Introduce a read-through cache abstraction in the usecase layer with generic payload support.

Introduce a rate limiter abstraction that enforces subject-operation windows via atomic counter increments.

Use Redis-compatible store interfaces and provide in-memory adapter for deterministic unit testing.

Persist optional rate-limit policy metadata in MySQL for future runtime policy management.

## Consequences
Positive outcomes:
- Reduced repeated read pressure on persistence for cacheable query responses.
- Controlled write burst behavior with deterministic request caps.
- Clean architecture boundaries preserved through store interfaces.

Tradeoffs:
- Cache invalidation strategy requires discipline in later phases.
- Additional operational tuning for window sizes and limits.
