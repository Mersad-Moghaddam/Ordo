-- name: CreateWorkspace :exec
INSERT INTO workspaces (
    id,
    workspace_key,
    display_name,
    created_by_user_id,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?);

-- name: FindWorkspaceByID :one
SELECT id, workspace_key, display_name, created_by_user_id, created_at, updated_at
FROM workspaces
WHERE id = ?;

-- name: FindWorkspaceByKey :one
SELECT id, workspace_key, display_name, created_by_user_id, created_at, updated_at
FROM workspaces
WHERE workspace_key = ?;

-- name: ListWorkspacesByUserID :many
SELECT workspace_table.id, workspace_table.workspace_key, workspace_table.display_name, workspace_table.created_by_user_id, workspace_table.created_at, workspace_table.updated_at
FROM workspaces AS workspace_table
INNER JOIN workspace_memberships AS membership_table
    ON membership_table.workspace_id = workspace_table.id
WHERE membership_table.user_id = ?
ORDER BY workspace_table.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountWorkspacesByUserID :one
SELECT COUNT(*)
FROM workspace_memberships
WHERE user_id = ?;
