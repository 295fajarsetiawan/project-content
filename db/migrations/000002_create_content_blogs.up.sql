-- 000002_create_content_blogs.up.sql

CREATE TABLE IF NOT EXISTS content_blog_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(150) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_content_blog_categories_app_slug UNIQUE (app_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_content_blog_categories_app_id ON content_blog_categories(app_id);
CREATE INDEX IF NOT EXISTS idx_content_blog_categories_slug ON content_blog_categories(app_id, slug);

CREATE TRIGGER trg_content_blog_categories_updated_at
BEFORE UPDATE ON content_blog_categories
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

CREATE TABLE IF NOT EXISTS content_blogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    category_id UUID NOT NULL REFERENCES content_blog_categories(id) ON DELETE RESTRICT,
    author_id UUID,
    author_name VARCHAR(150) NOT NULL DEFAULT 'Admin',
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    excerpt TEXT NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    thumbnail_url TEXT NOT NULL DEFAULT '',
    tags TEXT[] NOT NULL DEFAULT '{}',
    meta_title VARCHAR(255) NOT NULL DEFAULT '',
    meta_description TEXT NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT', -- DRAFT, PUBLISHED, ARCHIVED
    is_featured BOOLEAN NOT NULL DEFAULT false,
    view_count BIGINT NOT NULL DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_content_blogs_app_slug UNIQUE (app_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_content_blogs_app_id ON content_blogs(app_id);
CREATE INDEX IF NOT EXISTS idx_content_blogs_lookup ON content_blogs(app_id, status, is_featured, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_blogs_category ON content_blogs(category_id);
CREATE INDEX IF NOT EXISTS idx_content_blogs_slug ON content_blogs(app_id, slug);
CREATE INDEX IF NOT EXISTS idx_content_blogs_tags ON content_blogs USING GIN(tags);

CREATE TRIGGER trg_content_blogs_updated_at
BEFORE UPDATE ON content_blogs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();
