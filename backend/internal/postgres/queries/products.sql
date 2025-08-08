-- name: CreateProduct :one
INSERT INTO products (name, description, price, category_id, image_url, stock_quantity)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductByID :one
SELECT p.*, 
       c.id AS category_id,
       c.name AS category_name, 
       c.description AS category_description
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.id = $1;

-- name: ProductExists :one
SELECT EXISTS(SELECT 1 FROM products WHERE id = $1) AS exists;

-- name: UpdateProduct :one
UPDATE products
SET name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    price = coalesce(sqlc.narg('price'), price),
    category_id = coalesce(sqlc.narg('category_id'), category_id),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    stock_quantity = coalesce(sqlc.narg('stock_quantity'), stock_quantity)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListProducts :many
SELECT p.*, 
       c.id AS category_id,
       c.name AS category_name, 
       c.description AS category_description
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE 
    p.deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(p.name) LIKE sqlc.narg('search')
        OR LOWER(p.description) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('price_from')::float IS NULL 
        OR p.price >= sqlc.narg('price_from')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL 
        OR p.price <= sqlc.narg('price_to')
    )
    AND (
        sqlc.narg('category_ids')::int[] IS NULL 
        OR p.category_id = ANY(sqlc.narg('category_ids')::int[])
    )
ORDER BY p.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountProducts :one
SELECT COUNT(*) AS total_products
FROM products
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
    )
    AND (
        sqlc.narg('category_ids')::int[] IS NULL 
        OR category_id = ANY(sqlc.narg('category_ids')::int[])
    );

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = now()
WHERE id = $1;