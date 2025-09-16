-- name: CreateProductStem :one
INSERT INTO product_stems (product_id, stem_count, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProductStemByID :one
SELECT * FROM product_stems
WHERE id = $1;

-- name: GetProductStemsByProductID :many
SELECT * FROM product_stems
WHERE product_id = $1;

-- name: UpdateProductStem :one
UPDATE product_stems
SET stem_count = coalesce(sqlc.narg('stem_count'), stem_count),
    price = coalesce(sqlc.narg('price'), price)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProductStemsByProductID :exec
DELETE FROM product_stems
WHERE product_id = $1;