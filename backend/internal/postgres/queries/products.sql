-- name: CreateProduct :one
INSERT INTO products (name, description, price, category_id, has_stems, is_message_card, is_add_on, is_flowers, image_url, stock_quantity)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: TotalProducts :one
SELECT COALESCE(COUNT(*), 0) AS total_products
FROM products
WHERE deleted_at IS NULL;

-- name: GetProductByID :one
SELECT p.*, 
       c.id AS category_id,
       c.name AS category_name, 
       c.description AS category_description
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.id = $1
GROUP BY p.id, c.id, c.name, c.description;

-- name: ProductExists :one
SELECT EXISTS(SELECT 1 FROM products WHERE id = $1) AS exists;

-- name: UpdateProduct :one
UPDATE products
SET name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    price = coalesce(sqlc.narg('price'), price),
    category_id = coalesce(sqlc.narg('category_id'), category_id),
    has_stems = coalesce(sqlc.narg('has_stems'), has_stems),
    is_message_card = coalesce(sqlc.narg('is_message_card'), is_message_card),
    is_flowers = coalesce(sqlc.narg('is_flowers'), is_flowers),
    is_add_on = coalesce(sqlc.narg('is_add_on'), is_add_on),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    stock_quantity = coalesce(sqlc.narg('stock_quantity'), stock_quantity)
WHERE id = sqlc.arg('id')
RETURNING *;

-- -- name: ListProducts :many
-- SELECT p.*, 
--        c.id AS category_id,
--        c.name AS category_name, 
--        c.description AS category_description
-- FROM products p
-- LEFT JOIN categories c ON p.category_id = c.id
-- WHERE 
--     p.deleted_at IS NULL
--     AND (
--         COALESCE(sqlc.narg('search'), '') = '' 
--         OR LOWER(p.name) LIKE sqlc.narg('search')
--         OR LOWER(p.description) LIKE sqlc.narg('search')
--     )
--     AND (
--         sqlc.narg('price_from')::float IS NULL 
--         OR p.price >= sqlc.narg('price_from')
--     )
--     AND (
--         sqlc.narg('price_to')::float IS NULL 
--         OR p.price <= sqlc.narg('price_to')
--     )
--     AND (
--         sqlc.narg('category_ids')::int[] IS NULL 
--         OR p.category_id = ANY(sqlc.narg('category_ids')::int[])
--     )
-- ORDER BY p.created_at DESC
-- LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- -- name: ListCountProducts :one
-- SELECT COUNT(*) AS total_products
-- FROM products
-- WHERE 
--     deleted_at IS NULL
--     AND (
--         COALESCE(sqlc.narg('search'), '') = '' 
--         OR LOWER(name) LIKE sqlc.narg('search')
--         OR LOWER(description) LIKE sqlc.narg('search')
--     )
--     AND (
--         sqlc.narg('price_from')::float IS NULL 
--         OR price >= sqlc.narg('price_from')
--     )
--     AND (
--         sqlc.narg('price_to')::float IS NULL 
--         OR price <= sqlc.narg('price_to')
--     )
--     AND (
--         sqlc.narg('category_ids')::int[] IS NULL 
--         OR category_id = ANY(sqlc.narg('category_ids')::int[])
--     );

-- name: ListProducts :many
SELECT 
    p.*,
    c.id AS category_id,
    c.name AS category_name, 
    c.description AS category_description,
    COALESCE(
        json_agg(
            jsonb_build_object(
                'id', ps.id,
                'product_id', ps.product_id,
                'stem_count', ps.stem_count,
                'price', ps.price
            )
        ) FILTER (WHERE ps.id IS NOT NULL), '[]'
    ) AS stems
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN product_stems ps ON ps.product_id = p.id
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
        sqlc.narg('is_message_card')::boolean IS NULL 
        OR p.is_message_card = sqlc.narg('is_message_card')
    )
    AND (
        sqlc.narg('is_flowers')::boolean IS NULL 
        OR p.is_flowers = sqlc.narg('is_flowers')
    )
    AND (
        sqlc.narg('is_add_on')::boolean IS NULL 
        OR p.is_add_on = sqlc.narg('is_add_on')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL 
        OR p.price <= sqlc.narg('price_to')
    )
    AND (
        sqlc.narg('category_ids')::int[] IS NULL 
        OR p.category_id = ANY(sqlc.narg('category_ids')::int[])
    )
GROUP BY p.id, c.id, c.name, c.description
ORDER BY p.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountProducts :one
SELECT COUNT(DISTINCT p.id) AS total_products
FROM products p
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
        sqlc.narg('is_message_card')::boolean IS NULL 
        OR p.is_message_card = sqlc.narg('is_message_card')
    )
    AND (
        sqlc.narg('is_flowers')::boolean IS NULL 
        OR p.is_flowers = sqlc.narg('is_flowers')
    )
    AND (
        sqlc.narg('is_add_on')::boolean IS NULL 
        OR p.is_add_on = sqlc.narg('is_add_on')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL 
        OR p.price <= sqlc.narg('price_to')
    )
    AND (
        sqlc.narg('category_ids')::int[] IS NULL 
        OR p.category_id = ANY(sqlc.narg('category_ids')::int[])
    );

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = now()
WHERE id = $1;