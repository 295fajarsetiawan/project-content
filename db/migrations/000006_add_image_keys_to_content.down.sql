-- 000006_add_image_keys_to_content.down.sql

ALTER TABLE content_testimonials DROP COLUMN IF EXISTS avatar_key;
ALTER TABLE content_custom_blocks DROP COLUMN IF EXISTS image_key;
ALTER TABLE content_gallery_items DROP COLUMN IF EXISTS thumbnail_key, DROP COLUMN IF EXISTS media_key;
ALTER TABLE content_galleries DROP COLUMN IF EXISTS cover_key;
ALTER TABLE content_blogs DROP COLUMN IF EXISTS thumbnail_key;
ALTER TABLE content_banners DROP COLUMN IF EXISTS mobile_image_key, DROP COLUMN IF EXISTS image_key;
