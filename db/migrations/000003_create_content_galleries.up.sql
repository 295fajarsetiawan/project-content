-- 000003_create_content_galleries.up.sql

CREATE TABLE IF NOT EXISTS content_galleries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category VARCHAR(50) NOT NULL DEFAULT 'GENERAL', -- UMRAH, HAJJ, TOUR, EVENT, GENERAL
    cover_url TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_content_galleries_app_slug UNIQUE (app_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_content_galleries_app_id ON content_galleries(app_id);
CREATE INDEX IF NOT EXISTS idx_content_galleries_slug ON content_galleries(app_id, slug);
CREATE INDEX IF NOT EXISTS idx_content_galleries_category ON content_galleries(app_id, category, is_active);

CREATE TRIGGER trg_content_galleries_updated_at
BEFORE UPDATE ON content_galleries
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

CREATE TABLE IF NOT EXISTS content_gallery_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    gallery_id UUID NOT NULL REFERENCES content_galleries(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL DEFAULT '',
    caption TEXT NOT NULL DEFAULT '',
    media_type VARCHAR(20) NOT NULL DEFAULT 'IMAGE', -- IMAGE, VIDEO
    media_url TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL DEFAULT '',
    sequence_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_content_gallery_items_app_id ON content_gallery_items(app_id);
CREATE INDEX IF NOT EXISTS idx_content_gallery_items_gallery ON content_gallery_items(gallery_id, is_active, sequence_order);

CREATE TRIGGER trg_content_gallery_items_updated_at
BEFORE UPDATE ON content_gallery_items
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();
