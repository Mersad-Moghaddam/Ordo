-- name: UpsertRateLimitPolicy :exec
INSERT INTO rate_limit_policies (
    policy_key,
    request_limit,
    window_seconds
) VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
    request_limit = VALUES(request_limit),
    window_seconds = VALUES(window_seconds);

-- name: FindRateLimitPolicyByKey :one
SELECT policy_key, request_limit, window_seconds, created_at, updated_at
FROM rate_limit_policies
WHERE policy_key = ?;
