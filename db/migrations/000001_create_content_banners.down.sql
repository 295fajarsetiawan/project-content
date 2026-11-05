-- 000001_create_content_banners.down.sql

DROP TRIGGER IF EXISTS trg_content_banners_updated_at ON content_banners;
DROP TABLE IF EXISTS content_banners CASCADE;
