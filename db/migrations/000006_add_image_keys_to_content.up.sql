-- 000006_add_image_keys_to_content.up.sql
-- Simpan object key storage untuk setiap gambar. Kolom URL lama tetap dipertahankan.

ALTER TABLE content_banners
    ADD COLUMN IF NOT EXISTS image_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS mobile_image_key TEXT NOT NULL DEFAULT '';

ALTER TABLE content_blogs
    ADD COLUMN IF NOT EXISTS thumbnail_key TEXT NOT NULL DEFAULT '';

ALTER TABLE content_galleries
    ADD COLUMN IF NOT EXISTS cover_key TEXT NOT NULL DEFAULT '';

ALTER TABLE content_gallery_items
    ADD COLUMN IF NOT EXISTS media_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS thumbnail_key TEXT NOT NULL DEFAULT '';

ALTER TABLE content_custom_blocks
    ADD COLUMN IF NOT EXISTS image_key TEXT NOT NULL DEFAULT '';

ALTER TABLE content_testimonials
    ADD COLUMN IF NOT EXISTS avatar_key TEXT NOT NULL DEFAULT '';
