-- name: CreateCategory :one
INSERT INTO categories (name, description, image_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: UpdateCategory :one
UPDATE categories
SET name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    image_url = coalesce(sqlc.narg('image_url'), image_url)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories
WHERE 
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
        OR LOWER(description) LIKE sqlc.narg('search')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCategoriesCount :one
SELECT COUNT(*) AS total_categories
FROM categories
WHERE 
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
        OR LOWER(description) LIKE sqlc.narg('search')
    );

-- name: DeleteCategory :exec
UPDATE categories
SET deleted_at = now()
WHERE id = $1;