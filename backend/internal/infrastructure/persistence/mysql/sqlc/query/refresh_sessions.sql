-- name: CreateRefreshSession :exec
INSERT INTO refresh_sessions (
    session_id,
    user_id,
    refresh_token_hash,
    refresh_token_version,
    issued_at,
    expires_at,
    revoked_at,
    replacement_session_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: FindRefreshSessionByID :one
SELECT session_id, user_id, refresh_token_hash, refresh_token_version, issued_at, expires_at, revoked_at, replacement_session_id
FROM refresh_sessions
WHERE session_id = ?;

-- name: RevokeRefreshSession :exec
UPDATE refresh_sessions
SET revoked_at = ?, replacement_session_id = ?
WHERE session_id = ?;
