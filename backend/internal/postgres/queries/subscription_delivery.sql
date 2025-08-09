-- name: CreateSubscriptionDelivery :one
INSERT INTO subscription_deliveries (description, user_subscription_id, delivered_on)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSubscriptionDeliveryByUserSubscriptionID :many
SELECT * FROM subscription_deliveries WHERE deleted_at IS NULL AND user_subscription_id = $1;

-- name: UpdateSubscriptionDelivery :one
UPDATE subscription_deliveries
SET description = coalesce(sqlc.narg('description'), description),
    delivered_on = coalesce(sqlc.narg('delivered_on'), delivered_on)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListSubscriptionDelivery :many
SELECT * FROM subscription_deliveries
WHERE 
    deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountSubscriptionDelivery :one
SELECT COUNT(*) AS total_subscription_deliveries 
FROM subscription_deliveries
WHERE 
    deleted_at IS NULL;

-- name: DeleteSubscriptionDelivery :exec
UPDATE subscription_deliveries
SET deleted_at = now()
WHERE id = $1;