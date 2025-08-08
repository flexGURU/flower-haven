-- name: CreateOrder :one
INSERT INTO orders (user_name, user_phone_number, payment_status, status, shipping_address)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: OrderExists :one
SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1) AS exists;

-- name: UpdateOrder :one
UPDATE orders
SET user_name = coalesce(sqlc.narg('user_name'), user_name),
    user_phone_number = coalesce(sqlc.narg('user_phone_number'), user_phone_number),
    payment_status = coalesce(sqlc.narg('payment_status'), payment_status),
    status = coalesce(sqlc.narg('status'), status),
    shipping_address = coalesce(sqlc.narg('shipping_address'), shipping_address)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrder :one
UPDATE orders
SET deleted_at = now()
WHERE id = $1
RETURNING *;

-- name: ListOrder :many
SELECT * FROM orders
WHERE
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(user_name) LIKE sqlc.narg('search')
        OR LOWER(user_phone_number) LIKE sqlc.narg('search')
        OR LOWER(shipping_address) LIKE sqlc.narg('search')
    )
    AND (
        COALESCE(sqlc.narg('payment_status'), '') = '' 
        OR LOWER(payment_status) LIKE sqlc.narg('payment_status')
    )
    AND (
        COALESCE(sqlc.narg('status'), '') = '' 
        OR LOWER(status) LIKE sqlc.narg('status')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');