-- 000005_seed_initial_contents.down.sql

DELETE FROM content_testimonials WHERE app_id = 'umroh';
DELETE FROM content_faqs WHERE app_id = 'umroh';
DELETE FROM content_custom_blocks WHERE app_id = 'umroh';
DELETE FROM content_gallery_items WHERE app_id = 'umroh';
DELETE FROM content_galleries WHERE app_id = 'umroh';
DELETE FROM content_blogs WHERE app_id = 'umroh';
DELETE FROM content_blog_categories WHERE app_id = 'umroh';
DELETE FROM content_banners WHERE app_id = 'umroh';
