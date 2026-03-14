# ADR 0006: Worker Idempotency, Retry, and DLQ Strategy

## Status
Accepted

## Context
Phase 5 requires stream consumption that can tolerate temporary downstream failures while preventing duplicate side effects.

## Decision
Implement worker polling against Redis Streams group `ordo_workers` on stream `ordo_events`.

Use idempotency keys for de-duplication before dispatching notifications.

On processing failure, republish with incremented attempts and exponential backoff metadata.

When attempts exceed threshold, route event to `ordo_events_dlq` and acknowledge original message.

Persist processed idempotency keys in `processed_events` with TTL-based expiration.

## Consequences
Positive outcomes:
- Duplicate-safe processing for retried and redelivered events.
- Explicit dead-letter handling for poison messages.
- Predictable retry strategy with bounded retries.

Tradeoffs:
- Additional storage and cleanup for idempotency records.
- Increased worker logic complexity around retry and DLQ routing.
