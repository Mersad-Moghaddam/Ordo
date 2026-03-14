CREATE TABLE IF NOT EXISTS processed_events (
    idempotency_key VARCHAR(255) PRIMARY KEY,
    processed_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX index_processed_events_expires_at (expires_at)
);
