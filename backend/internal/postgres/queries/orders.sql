-- name: CreateOrder :one
INSERT INTO orders (user_name, user_phone_number, total_amount, payment_status, status, shipping_address, delivery_date, time_slot, by_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: TotalOrders :one
SELECT COALESCE(COUNT(*), 0) AS total_orders
FROM orders
WHERE deleted_at IS NULL;

-- name: GetRecentOrders :many
SELECT * FROM orders
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 7;

-- name: GetOrderByFullDataID :one
SELECT 
  o.*,
  COALESCE(items.items, '[]') AS order_item_data
FROM orders o
LEFT JOIN LATERAL (
  SELECT json_agg(json_build_object(
    'id', oi.id,
    'order_id', oi.order_id,
    'product_id', oi.product_id,
    'quantity', oi.quantity,
    'amount', oi.amount,
    'current_product_details', json_build_object(
      'id', p.id,
      'name', p.name,
      'description', p.description,
      'price', p.price,
      'stock_quantity', p.stock_quantity,
      'image_url', p.image_url,
      'category_id', p.category_id,
      'created_at', p.created_at
    )
  )) AS items
  FROM order_items oi
  JOIN products p ON p.id = oi.product_id
  WHERE oi.order_id = o.id
) items ON true
WHERE o.id = $1;

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
RETURNING id;

-- name: DeleteOrder :exec
UPDATE orders
SET deleted_at = now()
WHERE id = $1;

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
        sqlc.narg('payment_status')::boolean IS NULL 
        OR payment_status = sqlc.narg('payment_status')
    )
    AND (
        COALESCE(sqlc.narg('status'), '') = '' 
        OR LOWER(status) LIKE sqlc.narg('status')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountOrder :one
SELECT COUNT(*) AS total_orders
FROM orders
WHERE
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(user_name) LIKE sqlc.narg('search')
        OR LOWER(user_phone_number) LIKE sqlc.narg('search')
        OR LOWER(shipping_address) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('payment_status')::boolean IS NULL 
        OR payment_status = sqlc.narg('payment_status')
    )
    AND (
        COALESCE(sqlc.narg('status'), '') = '' 
        OR LOWER(status) LIKE sqlc.narg('status')
    );