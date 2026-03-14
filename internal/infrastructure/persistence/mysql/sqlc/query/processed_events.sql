-- name: FindProcessedEvent :one
SELECT idempotency_key, processed_at, expires_at
FROM processed_events
WHERE idempotency_key = ? AND expires_at > NOW();

-- name: UpsertProcessedEvent :exec
INSERT INTO processed_events (
    idempotency_key,
    processed_at,
    expires_at
) VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
    processed_at = VALUES(processed_at),
    expires_at = VALUES(expires_at);

-- name: DeleteExpiredProcessedEvents :exec
DELETE FROM processed_events
WHERE expires_at <= NOW();
