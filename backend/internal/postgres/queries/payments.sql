-- name: CreatePayment :one
INSERT INTO payments (order_id, description, user_subscription_id, amount, payment_method, paid_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1;

-- name: GetPaymentsByOrderID :one
SELECT * FROM payments WHERE order_id = $1 LIMIT 1;

-- name: GetPaymentsByUserSubscriptionID :one
SELECT * FROM payments WHERE order_id = $1 LIMIT 1;

-- name: UpdatePayment :one
UPDATE payments
SET payment_method = coalesce(sqlc.narg('payment_method'), payment_method),
    amount = coalesce(sqlc.narg('amount'), amount),
    description = coalesce(sqlc.narg('description'), description),
    paid_at = coalesce(sqlc.narg('paid_at'), paid_at)
WHERE id = sqlc.arg('id')
RETURNING id;

-- name: ListPayments :many
SELECT * FROM payments
WHERE 
    (
        COALESCE(sqlc.narg('payment_method'), '') = '' 
        OR LOWER(payment_method) LIKE sqlc.narg('payment_method')
    )
    AND paid_at BETWEEN sqlc.narg('start_date') AND sqlc.narg('end_date')
ORDER BY paid_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountPayments :one
SELECT COUNT(*) AS total_payments
FROM payments
WHERE 
    (
        COALESCE(sqlc.narg('payment_method'), '') = '' 
        OR LOWER(payment_method) LIKE sqlc.narg('payment_method')
    )
    AND paid_at BETWEEN sqlc.narg('start_date') AND sqlc.narg('end_date');