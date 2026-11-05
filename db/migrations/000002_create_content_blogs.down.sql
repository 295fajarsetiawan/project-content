-- 000002_create_content_blogs.down.sql

DROP TRIGGER IF EXISTS trg_content_blogs_updated_at ON content_blogs;
DROP TABLE IF EXISTS content_blogs CASCADE;

DROP TRIGGER IF EXISTS trg_content_blog_categories_updated_at ON content_blog_categories;
DROP TABLE IF EXISTS content_blog_categories CASCADE;
