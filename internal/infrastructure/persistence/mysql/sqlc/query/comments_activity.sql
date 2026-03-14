-- name: CreateComment :exec
INSERT INTO comments (
    id,
    workspace_id,
    project_id,
    task_id,
    author_user_id,
    body,
    created_at,
    updated_at,
    deleted_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: FindCommentByID :one
SELECT id, workspace_id, project_id, task_id, author_user_id, body, created_at, updated_at, deleted_at
FROM comments
WHERE id = ?;

-- name: UpdateCommentBody :exec
UPDATE comments
SET body = ?, updated_at = NOW()
WHERE id = ?;

-- name: SoftDeleteComment :exec
UPDATE comments
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = ?;

-- name: ListCommentsByTaskID :many
SELECT id, workspace_id, project_id, task_id, author_user_id, body, created_at, updated_at, deleted_at
FROM comments
WHERE task_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: CountCommentsByTaskID :one
SELECT COUNT(*)
FROM comments
WHERE task_id = ?;

-- name: CreateActivityLog :exec
INSERT INTO activity_logs (
    id,
    workspace_id,
    project_id,
    task_id,
    actor_user_id,
    activity_type,
    payload,
    created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListActivityLogsByTaskID :many
SELECT id, workspace_id, project_id, task_id, actor_user_id, activity_type, payload, created_at
FROM activity_logs
WHERE task_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: CountActivityLogsByTaskID :one
SELECT COUNT(*)
FROM activity_logs
WHERE task_id = ?;
