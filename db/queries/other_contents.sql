-- 1. Custom Blocks

-- name: ListCustomBlocks :many
SELECT id, app_id, section_key, content_type, title, subtitle, body, image_url, metadata, sequence_order, is_active, created_at, updated_at
FROM content_custom_blocks
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('content_type')::text IS NULL OR sqlc.narg('content_type')::text = '' OR content_type = sqlc.narg('content_type'))
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
ORDER BY sequence_order ASC, created_at ASC;

-- name: GetCustomBlockBySection :one
SELECT id, app_id, section_key, content_type, title, subtitle, body, image_url, metadata, sequence_order, is_active, created_at, updated_at
FROM content_custom_blocks
WHERE app_id = $1 AND section_key = $2;

-- name: UpsertCustomBlock :one
INSERT INTO content_custom_blocks (
    app_id, section_key, content_type, title, subtitle, body, image_url, metadata, sequence_order, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (app_id, section_key) DO UPDATE SET
    content_type = EXCLUDED.content_type,
    title = EXCLUDED.title,
    subtitle = EXCLUDED.subtitle,
    body = EXCLUDED.body,
    image_url = EXCLUDED.image_url,
    metadata = EXCLUDED.metadata,
    sequence_order = EXCLUDED.sequence_order,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: DeleteCustomBlock :exec
DELETE FROM content_custom_blocks
WHERE id = $1 AND app_id = $2;

-- 2. FAQs

-- name: ListFAQs :many
SELECT id, app_id, category, question, answer, sequence_order, is_active, created_at, updated_at
FROM content_faqs
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('category')::text IS NULL OR sqlc.narg('category')::text = '' OR category = sqlc.narg('category'))
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
ORDER BY sequence_order ASC, created_at ASC;

-- name: CreateFAQ :one
INSERT INTO content_faqs (
    app_id, category, question, answer, sequence_order, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateFAQ :one
UPDATE content_faqs
SET
    category = $3,
    question = $4,
    answer = $5,
    sequence_order = $6,
    is_active = $7
WHERE id = $1 AND app_id = $2
RETURNING *;

-- name: DeleteFAQ :exec
DELETE FROM content_faqs
WHERE id = $1 AND app_id = $2;

-- 3. Testimonials

-- name: ListTestimonials :many
SELECT id, app_id, customer_name, customer_title, avatar_url, rating, comment, package_name, is_featured, is_active, created_at, updated_at
FROM content_testimonials
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
  AND (sqlc.narg('is_featured')::boolean IS NULL OR is_featured = sqlc.narg('is_featured'))
ORDER BY is_featured DESC, created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountTestimonials :one
SELECT COUNT(*)
FROM content_testimonials
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
  AND (sqlc.narg('is_featured')::boolean IS NULL OR is_featured = sqlc.narg('is_featured'));

-- name: CreateTestimonial :one
INSERT INTO content_testimonials (
    app_id, customer_name, customer_title, avatar_url, rating, comment, package_name, is_featured, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateTestimonial :one
UPDATE content_testimonials
SET
    customer_name = $3,
    customer_title = $4,
    avatar_url = $5,
    rating = $6,
    comment = $7,
    package_name = $8,
    is_featured = $9,
    is_active = $10
WHERE id = $1 AND app_id = $2
RETURNING *;

-- name: DeleteTestimonial :exec
DELETE FROM content_testimonials
WHERE id = $1 AND app_id = $2;
