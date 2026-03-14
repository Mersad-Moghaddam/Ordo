-- name: InsertOutboxEvent :exec
INSERT INTO outbox_events (
    id,
    aggregate_type,
    aggregate_id,
    event_type,
    payload,
    status,
    attempts,
    next_retry_at,
    idempotency_key
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
