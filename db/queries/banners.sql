-- name: ListBanners :many
SELECT 
    id, app_id, title, subtitle, description, image_url, mobile_image_url,
    target_url, button_text, position, sequence_order, start_date, end_date,
    is_active, additional_data, created_by, created_at, updated_at
FROM content_banners
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('position')::text IS NULL OR sqlc.narg('position')::text = '' OR position = sqlc.narg('position'))
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
ORDER BY sequence_order ASC, created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountBanners :one
SELECT COUNT(*)
FROM content_banners
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('position')::text IS NULL OR sqlc.narg('position')::text = '' OR position = sqlc.narg('position'))
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'));

-- name: GetBannerByID :one
SELECT 
    id, app_id, title, subtitle, description, image_url, mobile_image_url,
    target_url, button_text, position, sequence_order, start_date, end_date,
    is_active, additional_data, created_by, created_at, updated_at
FROM content_banners
WHERE id = $1 AND app_id = $2;

-- name: CreateBanner :one
INSERT INTO content_banners (
    app_id, title, subtitle, description, image_url, mobile_image_url,
    target_url, button_text, position, sequence_order, start_date, end_date,
    is_active, additional_data, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12,
    $13, $14, $15
)
RETURNING *;

-- name: UpdateBanner :one
UPDATE content_banners
SET
    title = $3,
    subtitle = $4,
    description = $5,
    image_url = $6,
    mobile_image_url = $7,
    target_url = $8,
    button_text = $9,
    position = $10,
    sequence_order = $11,
    start_date = $12,
    end_date = $13,
    is_active = $14,
    additional_data = $15
WHERE id = $1 AND app_id = $2
RETURNING *;

-- name: DeleteBanner :exec
DELETE FROM content_banners
WHERE id = $1 AND app_id = $2;
