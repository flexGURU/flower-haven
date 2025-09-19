-- name: CreatePaystackPayment :exec
INSERT INTO paystack_payments (amount, email, reference)
VALUES ($1, $2, $3);

-- name: UpdatePaystackPaymentStatus :exec
UPDATE paystack_payments
SET status = $2, updated_at = now()
WHERE reference = $1;

-- name: GetPaystackPaymentByReference :one
SELECT * FROM paystack_payments WHERE reference = $1;

-- name: ListPaystackPayments :many
SELECT * FROM paystack_payments
WHERE 
    (
        COALESCE(sqlc.narg('status')::text, '') = '' 
        OR LOWER(status) LIKE sqlc.narg('status')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountPaystackPayments :one
SELECT COUNT(*) AS total_paystack_payments
FROM paystack_payments
WHERE 
    (
        COALESCE(sqlc.narg('status')::text, '') = '' 
        OR LOWER(status) LIKE sqlc.narg('status')
    );

-- name: CreatePaystackEvent :exec
INSERT INTO paystack_events (event, data)
VALUES ($1, $2);

-- name: ListPaystackEvents :many
SELECT * FROM paystack_events
WHERE 
    (
        COALESCE(sqlc.narg('event')::text, '') = '' 
        OR LOWER(event) LIKE sqlc.narg('event')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountPaystackEvents :one
SELECT COUNT(*) AS total_paystack_events
FROM paystack_events
WHERE 
    (
        COALESCE(sqlc.narg('event')::text, '') = '' 
        OR LOWER(event) LIKE sqlc.narg('event')
    );