-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, amount)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: GetOrderItemsByProductID :many
SELECT * FROM order_items WHERE product_id = $1;
