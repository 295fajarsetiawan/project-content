-- 000003_create_content_galleries.down.sql

DROP TRIGGER IF EXISTS trg_content_gallery_items_updated_at ON content_gallery_items;
DROP TABLE IF EXISTS content_gallery_items CASCADE;

DROP TRIGGER IF EXISTS trg_content_galleries_updated_at ON content_galleries;
DROP TABLE IF EXISTS content_galleries CASCADE;
