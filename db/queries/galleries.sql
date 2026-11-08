-- Galleries

-- name: ListGalleries :many
SELECT id, app_id, title, slug, description, category, cover_url, is_active, created_at, updated_at
FROM content_galleries
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('category')::text IS NULL OR sqlc.narg('category')::text = '' OR category = sqlc.narg('category'))
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountGalleries :one
SELECT COUNT(*)
FROM content_galleries
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('category')::text IS NULL OR sqlc.narg('category')::text = '' OR category = sqlc.narg('category'))
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'));

-- name: GetGalleryByID :one
SELECT id, app_id, title, slug, description, category, cover_url, is_active, created_at, updated_at
FROM content_galleries
WHERE id = $1 AND app_id = $2;

-- name: GetGalleryBySlug :one
SELECT id, app_id, title, slug, description, category, cover_url, is_active, created_at, updated_at
FROM content_galleries
WHERE app_id = $1 AND slug = $2;

-- name: CreateGallery :one
INSERT INTO content_galleries (
    app_id, title, slug, description, category, cover_url, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteGallery :exec
DELETE FROM content_galleries
WHERE id = $1 AND app_id = $2;

-- Gallery Items

-- name: ListGalleryItems :many
SELECT id, app_id, gallery_id, title, caption, media_type, media_url, thumbnail_url, sequence_order, is_active, created_at, updated_at
FROM content_gallery_items
WHERE gallery_id = sqlc.arg('gallery_id')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
ORDER BY sequence_order ASC, created_at ASC;

-- name: AddGalleryItem :one
INSERT INTO content_gallery_items (
    app_id, gallery_id, title, caption, media_type, media_url, thumbnail_url, sequence_order, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: DeleteGalleryItem :exec
DELETE FROM content_gallery_items
WHERE id = $1 AND app_id = $2;
