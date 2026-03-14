# ADR 0004: Task Outbox and Redis Streams Producer

## Status
Accepted

## Context
Phase 3 requires task lifecycle operations to emit reliable domain events without coupling task transaction success to stream broker availability.

## Decision
Persist task write events into `outbox_events` within the same transaction as task state changes.

Introduce a producer component that reads pending outbox rows and publishes them to Redis Streams `ordo_events`, then marks rows as published.

On publish failure, persist retry metadata with exponential backoff and keep event idempotency keys in stream values for downstream consumer deduplication.

## Consequences
Positive outcomes:
- Task writes remain durable even when stream publication is temporarily unavailable.
- Backoff retries and idempotency keys improve robustness and duplicate safety.
- Producer remains infrastructure-only and detached from domain logic.

Tradeoffs:
- Additional persistence operations for every state-changing task command.
- Requires operational scheduling for producer polling cadence.
