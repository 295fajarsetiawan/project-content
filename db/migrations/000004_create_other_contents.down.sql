-- 000004_create_other_contents.down.sql

DROP TRIGGER IF EXISTS trg_content_testimonials_updated_at ON content_testimonials;
DROP TABLE IF EXISTS content_testimonials CASCADE;

DROP TRIGGER IF EXISTS trg_content_faqs_updated_at ON content_faqs;
DROP TABLE IF EXISTS content_faqs CASCADE;

DROP TRIGGER IF EXISTS trg_content_custom_blocks_updated_at ON content_custom_blocks;
DROP TABLE IF EXISTS content_custom_blocks CASCADE;
