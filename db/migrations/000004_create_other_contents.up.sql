-- 000004_create_other_contents.up.sql

-- 1. Custom CMS Blocks & Static Pages (about-us, terms, privacy, hero-stats, etc.)
CREATE TABLE IF NOT EXISTS content_custom_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    section_key VARCHAR(100) NOT NULL,
    content_type VARCHAR(50) NOT NULL DEFAULT 'SECTION', -- PAGE, SECTION, ANNOUNCEMENT, PARTNER, WIDGET
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    sequence_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_content_custom_blocks_app_section UNIQUE (app_id, section_key)
);

CREATE INDEX IF NOT EXISTS idx_content_custom_blocks_app_id ON content_custom_blocks(app_id);
CREATE INDEX IF NOT EXISTS idx_content_custom_blocks_type ON content_custom_blocks(app_id, content_type, is_active);
CREATE INDEX IF NOT EXISTS idx_content_custom_blocks_metadata ON content_custom_blocks USING GIN(metadata);

CREATE TRIGGER trg_content_custom_blocks_updated_at
BEFORE UPDATE ON content_custom_blocks
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

-- 2. Frequently Asked Questions (FAQ)
CREATE TABLE IF NOT EXISTS content_faqs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    category VARCHAR(100) NOT NULL DEFAULT 'UMUM',
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    sequence_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_content_faqs_app_id ON content_faqs(app_id);
CREATE INDEX IF NOT EXISTS idx_content_faqs_category ON content_faqs(app_id, category, is_active, sequence_order);

CREATE TRIGGER trg_content_faqs_updated_at
BEFORE UPDATE ON content_faqs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

-- 3. Testimonials & Customer Reviews
CREATE TABLE IF NOT EXISTS content_testimonials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    customer_name VARCHAR(150) NOT NULL,
    customer_title VARCHAR(150) NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    rating INT NOT NULL DEFAULT 5 CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NOT NULL,
    package_name VARCHAR(255) NOT NULL DEFAULT '',
    is_featured BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_content_testimonials_app_id ON content_testimonials(app_id);
CREATE INDEX IF NOT EXISTS idx_content_testimonials_lookup ON content_testimonials(app_id, is_active, is_featured, created_at DESC);

CREATE TRIGGER trg_content_testimonials_updated_at
BEFORE UPDATE ON content_testimonials
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();
