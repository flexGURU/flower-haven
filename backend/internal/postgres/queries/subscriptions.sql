-- name: CreateSubscription :one
INSERT INTO subscriptions (name, description, product_ids, add_ons, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: SubscriptionExists :one
SELECT EXISTS(SELECT 1 FROM subscriptions WHERE id = $1) as exists;

-- name: GetSubscriptionByID :one
SELECT 
  s.*,
  COALESCE(p1.products_json, '[]') AS products_data,
  COALESCE(p2.add_ons_json, '[]') AS add_ons_data
FROM subscriptions s
LEFT JOIN LATERAL (
  SELECT json_agg(p.*) AS products_json
  FROM products p
  WHERE p.id = ANY(s.product_ids)
) p1 ON true
LEFT JOIN LATERAL (
  SELECT json_agg(p.*) AS add_ons_json
  FROM products p
  WHERE p.id = ANY(s.add_ons)
) p2 ON true
WHERE s.id = $1;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    product_ids = coalesce(sqlc.narg('product_ids'), product_ids),
    add_ons = coalesce(sqlc.narg('add_ons'), add_ons),
    price = coalesce(sqlc.narg('price'), price)
WHERE id = sqlc.arg('id')
RETURNING id;

-- name: DeleteSubscription :exec
UPDATE subscriptions
SET deleted_at = now()
WHERE id = $1;

-- name: ListSubscriptions :many
SELECT 
  s.*,
  COALESCE(p.products, '[]') AS products_data,
  COALESCE(a.add_ons, '[]') AS add_ons_data
FROM subscriptions s
LEFT JOIN LATERAL (
  SELECT json_agg(json_build_object(
    'id', p.id,
    'name', p.name,
    'price', p.price
  )) AS products
  FROM products p
  WHERE p.id = ANY(s.product_ids)
) p ON true
LEFT JOIN LATERAL (
  SELECT json_agg(json_build_object(
    'id', p.id,
    'name', p.name,
    'price', p.price
  )) AS add_ons
  FROM products p
  WHERE p.id = ANY(s.add_ons)
) a ON true
WHERE 
    s.deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(s.name) LIKE sqlc.narg('search')
        OR LOWER(s.description) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('price_from')::float IS NULL 
        OR s.price >= sqlc.narg('price_from')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL 
        OR s.price <= sqlc.narg('price_to')
    )
ORDER BY s.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListSubscriptionsCount :one
SELECT COUNT(*) AS total_subscriptions
WHERE 
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
        OR LOWER(description) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('price_from')::float IS NULL 
        OR price >= sqlc.narg('price_from')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL 
        OR price <= sqlc.narg('price_to')
    );