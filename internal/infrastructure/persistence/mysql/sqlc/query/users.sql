-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    password_hash,
    role,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?);

-- name: FindUserByEmail :one
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = ?;

-- name: FindUserByID :one
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE id = ?;
