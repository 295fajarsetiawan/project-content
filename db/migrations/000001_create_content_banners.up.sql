-- 000001_create_content_banners.up.sql

CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS content_banners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(50) NOT NULL DEFAULT 'umroh',
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL,
    mobile_image_url TEXT NOT NULL DEFAULT '',
    target_url TEXT NOT NULL DEFAULT '',
    button_text VARCHAR(100) NOT NULL DEFAULT '',
    position VARCHAR(50) NOT NULL DEFAULT 'HERO', -- HERO, PROMO, POPUP, SIDEBAR
    sequence_order INT NOT NULL DEFAULT 0,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    additional_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_content_banners_app_id ON content_banners(app_id);
CREATE INDEX IF NOT EXISTS idx_content_banners_lookup ON content_banners(app_id, position, is_active, sequence_order);
CREATE INDEX IF NOT EXISTS idx_content_banners_dates ON content_banners(start_date, end_date);

CREATE TRIGGER trg_content_banners_updated_at
BEFORE UPDATE ON content_banners
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();
