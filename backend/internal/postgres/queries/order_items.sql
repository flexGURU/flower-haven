-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, amount)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetOrderItemsByProductID :many
SELECT 
    oi.*,
    COALESCE(p1.order_json, '{}') AS order_data
FROM order_items oi
LEFT JOIN LATERAL (
    SELECT json_build_object(
        'id', o.id,
        'user_name', o.user_name,
        'user_phone_number', o.user_phone_number,
        'total_amount', o.total_amount,
        'payment_status', o.payment_status,
        'status', o.status
    ) AS order_json
    FROM orders o
    WHERE o.id = oi.order_id
) p1 ON true
WHERE oi.product_id = $1
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetCountOrderItemsByProductID :one
SELECT COUNT(*) AS total_order_items
FROM order_items
WHERE product_id = $1;
