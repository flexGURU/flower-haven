-- name: CreateUser :one
INSERT INTO users (name, email, address, phone_number, refresh_token, password, is_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1) AS exists;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name = coalesce(sqlc.narg('name'), name),
    address = coalesce(sqlc.narg('address'), address),
    phone_number = coalesce(sqlc.narg('phone_number'), phone_number),
    password = coalesce(sqlc.narg('password'), password),
    refresh_token = coalesce(sqlc.narg('refresh_token'), refresh_token),
    is_admin = coalesce(sqlc.narg('is_admin'), is_admin),
    is_active = coalesce(sqlc.narg('is_active'), is_active)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
WHERE 
    (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
        OR LOWER(email) LIKE sqlc.narg('search')
        OR LOWER(phone_number) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('is_active')::boolean IS NULL 
        OR is_active = sqlc.narg('is_active')
    )
    AND (
        sqlc.narg('is_admin')::boolean IS NULL 
        OR is_admin = sqlc.narg('is_admin')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListUsersCount :one
SELECT COUNT(*) AS total_users
FROM users
WHERE 
    (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
        OR LOWER(email) LIKE sqlc.narg('search')
        OR LOWER(phone_number) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('is_active')::boolean IS NULL 
        OR is_active = sqlc.narg('is_active')
    )
    AND (
        sqlc.narg('is_admin')::boolean IS NULL 
        OR is_admin = sqlc.narg('is_admin')
    );