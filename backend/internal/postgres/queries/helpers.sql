-- name: ListAddOns :many
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
WHERE p.deleted_at IS NULL AND p.is_add_on = TRUE
GROUP BY p.id, c.id, c.name, c.description
ORDER BY p.created_at DESC;

-- name: ListMessageCards :many
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
WHERE p.deleted_at IS NULL AND p.is_message_card = TRUE
GROUP BY p.id, c.id, c.name, c.description
ORDER BY p.created_at DESC;