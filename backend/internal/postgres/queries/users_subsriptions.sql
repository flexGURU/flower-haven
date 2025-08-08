-- name: CreateUserSubscription :one
INSERT INTO user_subscriptions (user_id, subscription_id, start_date, end_date, day_of_week)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UserSubscriptionExists :one
SELECT EXISTS(SELECT 1 FROM user_subscriptions WHERE id = $1) AS exists;

-- name: GetUserSubscriptionByID :one
SELECT 
    us.*,
    COALESCE(p1.user_json, '[]') AS user_data,
    COALESCE(p2.subscription_json, '[]') AS subscription_data,
    COALESCE(p3.payment_json, '[]') AS payment_data
FROM user_subscriptions us
LEFT JOIN LATERAL (
    SELECT json_agg(u.*) AS user_json
    FROM users u
    WHERE u.id = us.user_id
) p1 ON true
LEFT JOIN LATERAL (
    SELECT json_agg(s.*) AS subscription_json
    FROM subscriptions s
    WHERE s.id = us.subscription_id
) p2 ON true
LEFT JOIN LATERAL (
    SELECT json_agg(p.*) AS payment_json
    FROM payments p
    WHERE p.user_subscription_id IS NOT NULL
        AND p.user_subscription_id = sqlc.narg('user_subscription_id')
) p3 ON true
WHERE us.id = sqlc.arg('id');

-- name: GetUserSubscriptionsByUserID :many
SELECT 
    us.*,
    COALESCE(p1.subscription_json, '[]') AS subscription_data
FROM user_subscriptions us
LEFT JOIN LATERAL (
    SELECT json_agg(json_build_object(
        'id', s.id,
        'name', s.name,
        'price', s.price 
    )) AS subscription_json
    FROM subscriptions s
    WHERE s.id = us.subscription_id
) p1 ON true
WHERE us.user_id = $1;

-- name: UpdateUserSubscription :one
UPDATE user_subscriptions
SET start_date = coalesce(sqlc.narg('start_date'), start_date),
    end_date = coalesce(sqlc.narg('end_date'), end_date),
    day_of_week = coalesce(sqlc.narg('day_of_week'), day_of_week),
    status = coalesce(sqlc.narg('status'), status)
WHERE id = sqlc.arg('id')
RETURNING id;

-- name: ListUserSubscriptions :many
SELECT 
    us.*,
    COALESCE(p1.user_json, '[]') AS user_data,
    COALESCE(p2.subscription_json, '[]') AS subscription_data
FROM user_subscriptions us
LEFT JOIN LATERAL (
    SELECT json_agg(json_build_object(
        'id', u.id,
        'name', u.name,
        'email', u.email,
        'phone_number', u.phone_number
    )) AS user_json
    FROM users u
    WHERE u.id = us.user_id
) p1 ON true
LEFT JOIN LATERAL (
    SELECT json_agg(json_build_object(
        'id', s.id,
        'name', s.name,
        'price', s.price 
    )) AS subscription_json
    FROM subscriptions s
    WHERE s.id = us.subscription_id
) p2 ON true
WHERE
    us.deleted_at IS NULL
    AND (
        sqlc.narg('status')::boolean IS NULL 
        OR us.status = sqlc.narg('status')
    )
ORDER BY us.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListCountUserSubscriptions :one
SELECT COUNT(*) AS total_user_subscriptions
WHERE
    deleted_at IS NULL
    AND (
        sqlc.narg('status')::boolean IS NULL 
        OR status = sqlc.narg('status')
    );

-- name: DeleteUserSubscription :exec
UPDATE user_subscriptions
SET deleted_at = now()
WHERE id = $1;
