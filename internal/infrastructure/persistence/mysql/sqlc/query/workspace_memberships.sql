-- name: CreateWorkspaceMembership :exec
INSERT INTO workspace_memberships (
    workspace_id,
    user_id,
    membership_role,
    invited_by_user_id,
    joined_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?);

-- name: FindWorkspaceMembership :one
SELECT workspace_id, user_id, membership_role, invited_by_user_id, joined_at, updated_at
FROM workspace_memberships
WHERE workspace_id = ? AND user_id = ?;

-- name: UpdateWorkspaceMembershipRole :exec
UPDATE workspace_memberships
SET membership_role = ?, updated_at = ?
WHERE workspace_id = ? AND user_id = ?;
