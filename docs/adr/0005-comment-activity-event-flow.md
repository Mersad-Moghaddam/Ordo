# ADR 0005: Comment Activity and Event Flow

## Status
Accepted

## Context
Phase 4 requires collaborative comment operations to be auditable and evented for downstream notifications and timelines.

## Decision
Implement comment commands with transactional writes across comments, activity logs, and outbox events.

Emit activity records for comment create, update, and delete operations with payloads that capture actor and comment fields.

Emit outbox events mirroring activity types so workers can project notifications and external integrations asynchronously.

## Consequences
Positive outcomes:
- Comment timeline remains queryable from relational activity logs.
- Evented side effects are decoupled from request latency.
- Consistency preserved between comment state and emitted event intent.

Tradeoffs:
- Additional write amplification for each comment command.
- Requires worker phase to consume and fan out collaboration events.
