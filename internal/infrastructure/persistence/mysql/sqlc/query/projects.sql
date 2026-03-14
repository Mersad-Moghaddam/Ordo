-- name: CreateProject :exec
INSERT INTO projects (
    id,
    workspace_id,
    project_key,
    display_name,
    description,
    created_by_user_id,
    created_at,
    updated_at,
    archived_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: FindProjectByID :one
SELECT id, workspace_id, project_key, display_name, description, created_by_user_id, created_at, updated_at, archived_at
FROM projects
WHERE id = ?;

-- name: FindProjectByWorkspaceAndKey :one
SELECT id, workspace_id, project_key, display_name, description, created_by_user_id, created_at, updated_at, archived_at
FROM projects
WHERE workspace_id = ? AND project_key = ?;

-- name: ListProjectsByWorkspaceID :many
SELECT id, workspace_id, project_key, display_name, description, created_by_user_id, created_at, updated_at, archived_at
FROM projects
WHERE workspace_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountProjectsByWorkspaceID :one
SELECT COUNT(*)
FROM projects
WHERE workspace_id = ?;
