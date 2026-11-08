-- Blog Categories

-- name: ListBlogCategories :many
SELECT id, app_id, name, slug, description, is_active, created_at, updated_at
FROM content_blog_categories
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('search')::text IS NULL OR sqlc.narg('search')::text = '' OR name ILIKE '%' || sqlc.narg('search') || '%' OR description ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY name ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountBlogCategories :one
SELECT COUNT(*)
FROM content_blog_categories
WHERE app_id = sqlc.arg('app_id')
  AND (sqlc.narg('search')::text IS NULL OR sqlc.narg('search')::text = '' OR name ILIKE '%' || sqlc.narg('search') || '%' OR description ILIKE '%' || sqlc.narg('search') || '%');

-- name: GetBlogCategoryByID :one
SELECT id, app_id, name, slug, description, is_active, created_at, updated_at
FROM content_blog_categories
WHERE id = $1 AND app_id = $2;

-- name: GetBlogCategoryBySlug :one
SELECT id, app_id, name, slug, description, is_active, created_at, updated_at
FROM content_blog_categories
WHERE app_id = $1 AND slug = $2;

-- name: CreateBlogCategory :one
INSERT INTO content_blog_categories (
    app_id, name, slug, description, is_active
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateBlogCategory :one
UPDATE content_blog_categories
SET
    name = $3,
    slug = $4,
    description = $5,
    is_active = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND app_id = $2
RETURNING *;

-- Blogs

-- name: ListBlogs :many
SELECT 
    b.id, b.app_id, b.category_id, b.author_id, b.author_name,
    b.title, b.slug, b.excerpt, b.body, b.thumbnail_url, b.tags,
    b.meta_title, b.meta_description, b.status, b.is_featured,
    b.view_count, b.published_at, b.created_at, b.updated_at,
    c.name AS category_name, c.slug AS category_slug
FROM content_blogs b
JOIN content_blog_categories c ON b.category_id = c.id
WHERE b.app_id = sqlc.arg('app_id')
  AND (sqlc.narg('category_id')::uuid IS NULL OR b.category_id = sqlc.narg('category_id'))
  AND (sqlc.narg('status')::text IS NULL OR sqlc.narg('status')::text = '' OR b.status = sqlc.narg('status'))
  AND (sqlc.narg('is_featured')::boolean IS NULL OR b.is_featured = sqlc.narg('is_featured'))
  AND (sqlc.narg('search')::text IS NULL OR sqlc.narg('search')::text = '' OR b.title ILIKE '%' || sqlc.narg('search') || '%' OR b.excerpt ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY b.is_featured DESC, b.published_at DESC NULLS LAST, b.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountBlogs :one
SELECT COUNT(*)
FROM content_blogs b
WHERE b.app_id = sqlc.arg('app_id')
  AND (sqlc.narg('category_id')::uuid IS NULL OR b.category_id = sqlc.narg('category_id'))
  AND (sqlc.narg('status')::text IS NULL OR sqlc.narg('status')::text = '' OR b.status = sqlc.narg('status'))
  AND (sqlc.narg('is_featured')::boolean IS NULL OR b.is_featured = sqlc.narg('is_featured'))
  AND (sqlc.narg('search')::text IS NULL OR sqlc.narg('search')::text = '' OR b.title ILIKE '%' || sqlc.narg('search') || '%' OR b.excerpt ILIKE '%' || sqlc.narg('search') || '%');

-- name: GetBlogByID :one
SELECT 
    b.id, b.app_id, b.category_id, b.author_id, b.author_name,
    b.title, b.slug, b.excerpt, b.body, b.thumbnail_url, b.tags,
    b.meta_title, b.meta_description, b.status, b.is_featured,
    b.view_count, b.published_at, b.created_at, b.updated_at,
    c.name AS category_name, c.slug AS category_slug
FROM content_blogs b
JOIN content_blog_categories c ON b.category_id = c.id
WHERE b.id = $1 AND b.app_id = $2;

-- name: GetBlogBySlug :one
SELECT 
    b.id, b.app_id, b.category_id, b.author_id, b.author_name,
    b.title, b.slug, b.excerpt, b.body, b.thumbnail_url, b.tags,
    b.meta_title, b.meta_description, b.status, b.is_featured,
    b.view_count, b.published_at, b.created_at, b.updated_at,
    c.name AS category_name, c.slug AS category_slug
FROM content_blogs b
JOIN content_blog_categories c ON b.category_id = c.id
WHERE b.app_id = $1 AND b.slug = $2;

-- name: IncrementBlogViewCount :exec
UPDATE content_blogs
SET view_count = view_count + 1
WHERE id = $1;

-- name: CreateBlog :one
INSERT INTO content_blogs (
    app_id, category_id, author_id, author_name,
    title, slug, excerpt, body, thumbnail_url, tags,
    meta_title, meta_description, status, is_featured, published_at
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15
)
RETURNING *;

-- name: UpdateBlog :one
UPDATE content_blogs
SET
    category_id = $3,
    title = $4,
    slug = $5,
    excerpt = $6,
    body = $7,
    thumbnail_url = $8,
    tags = $9,
    meta_title = $10,
    meta_description = $11,
    status = $12,
    is_featured = $13
WHERE id = $1 AND app_id = $2
RETURNING *;

-- name: DeleteBlog :exec
DELETE FROM content_blogs
WHERE id = $1 AND app_id = $2;
