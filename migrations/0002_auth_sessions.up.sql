CREATE TABLE IF NOT EXISTS refresh_sessions (
    session_id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    refresh_token_version BIGINT NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    replacement_session_id CHAR(36) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX index_refresh_sessions_user_id (user_id),
    INDEX index_refresh_sessions_expires_at (expires_at),
    CONSTRAINT fk_refresh_sessions_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
