-- name: CreateTask :exec
INSERT INTO tasks (
    id,
    workspace_id,
    project_id,
    title,
    description,
    status,
    priority,
    assignee_user_id,
    created_by_user_id,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: FindTaskByID :one
SELECT id, workspace_id, project_id, title, description, status, priority, assignee_user_id, created_by_user_id, created_at, updated_at
FROM tasks
WHERE id = ?;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = ?, updated_at = ?
WHERE id = ?;

-- name: ListTasksByProjectID :many
SELECT id, workspace_id, project_id, title, description, status, priority, assignee_user_id, created_by_user_id, created_at, updated_at
FROM tasks
WHERE project_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountTasksByProjectID :one
SELECT COUNT(*)
FROM tasks
WHERE project_id = ?;

-- name: ListPendingOutboxEvents :many
SELECT id, aggregate_type, aggregate_id, event_type, payload, status, attempts, next_retry_at, idempotency_key, created_at, updated_at
FROM outbox_events
WHERE status = 'pending' OR (status = 'retry' AND (next_retry_at IS NULL OR next_retry_at <= NOW()))
ORDER BY created_at ASC
LIMIT ?;

-- name: MarkOutboxEventPublished :exec
UPDATE outbox_events
SET status = 'published', updated_at = NOW()
WHERE id = ?;

-- name: MarkOutboxEventRetry :exec
UPDATE outbox_events
SET status = 'retry', attempts = ?, next_retry_at = FROM_UNIXTIME(?), updated_at = NOW()
WHERE id = ?;
