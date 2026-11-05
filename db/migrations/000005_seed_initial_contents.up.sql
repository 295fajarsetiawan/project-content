-- 000005_seed_initial_contents.up.sql

-- 1. Seed Banners
INSERT INTO content_banners (app_id, title, subtitle, description, image_url, target_url, button_text, position, sequence_order, is_active)
VALUES
('umroh', 'Raih Kekhusyukan Ibadah di Baitullah', 'Paket Umroh Syawal 1448 H', 'Dapatkan pengalaman spiritual terbaik bersama pembimbing ibadah terpercaya sesuai sunnah dengan fasilitas hotel bintang 5.', 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?auto=format&fit=crop&w=1600&q=80', '/packages/umroh-reguler-syawal-berkah-9d', 'Lihat Paket', 'HERO', 1, true),
('umroh', 'Haji Furoda VIP Tanpa Antri', 'Keberangkatan Musim 2027', 'Pendaftaran Haji Khusus Mujamalah Resmi Kemenag dengan maktab VIP AC dekat jamarat Mina.', 'https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?auto=format&fit=crop&w=1600&q=80', '/packages/haji-khusus-furoda-vip-25d', 'Konsultasi Sekarang', 'HERO', 2, true),
('umroh', 'Promo Early Bird Umroh Awal Musim', 'Diskon Spesial Rp 2.000.000', 'Segera daftarkan diri dan keluarga Anda sebelum kuota promo habis!', 'https://images.unsplash.com/photo-1542314831-068cd1dbfeeb?auto=format&fit=crop&w=1200&q=80', '/packages', 'Klaim Promo', 'PROMO', 1, true)
ON CONFLICT DO NOTHING;

-- 2. Seed Blog Categories
INSERT INTO content_blog_categories (id, app_id, name, slug, description, is_active)
VALUES
('11111111-1111-1111-1111-111111111101', 'umroh', 'Panduan Ibadah', 'panduan-ibadah', 'Tata cara manasik, rukun, wajib, dan sunnah umroh serta haji', true),
('11111111-1111-1111-1111-111111111102', 'umroh', 'Destinasi & Sejarah', 'destinasi-sejarah', 'Jejak sejarah Islam di Makkah, Madinah, dan kota-kota bersejarah', true),
('11111111-1111-1111-1111-111111111103', 'umroh', 'Tips & Informasi Perjalanan', 'tips-informasi', 'Persiapan paspor, vaksin, perlengkapan, dan tips hemat selama di tanah suci', true)
ON CONFLICT (app_id, slug) DO NOTHING;

-- 3. Seed Blogs
INSERT INTO content_blogs (app_id, category_id, author_name, title, slug, excerpt, body, thumbnail_url, tags, meta_title, meta_description, status, is_featured, published_at)
VALUES
('umroh', '11111111-1111-1111-1111-111111111101', 'Ustadz Pembimbing', 'Panduan Lengkap Tata Cara Manasik Umroh Sesuai Sunnah', 'panduan-lengkap-manasik-umroh-sunnah', 'Pelajari tahapan ihram, tawaf, sa''i, dan tahallul dengan penjelasan ringkas dan dalil shahih.', 'Ibadah Umroh adalah salah satu ibadah agung yang mendatangkan pahala berlipat ganda. Rukun umroh terdiri dari empat hal: 1. Ihram dari miqat disertai niat, 2. Tawaf mengelilingi Ka''bah sebanyak 7 putaran, 3. Sa''i antara Shafa dan Marwah, 4. Tahallul (memotong atau mencukur rambut). Pastikan jamaah menjaga adab dan larangan ihram selama menjalankan ibadah.', 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?auto=format&fit=crop&w=800&q=80', ARRAY['manasik', 'umroh', 'sunnah', 'panduan'], 'Panduan Lengkap Manasik Umroh Sesuai Sunnah', 'Panduan ringkas dan lengkap tata cara manasik ibadah umroh sesuai sunnah Rasulullah SAW.', 'PUBLISHED', true, CURRENT_TIMESTAMP),
('umroh', '11111111-1111-1111-1111-111111111103', 'Tim Travel', 'Daftar Perlengkapan Wajib yang Harus Dibawa Jamaah Umroh', 'daftar-perlengkapan-wajib-jamaah-umroh', 'Daftar checklist barang bawaan penting mulai dari pakaian ihram, dokumen hingga obat-obatan pribadi.', 'Menyiapkan perlengkapan sebelum berangkat ke tanah suci sangat krusial. Jamaah pria disarankan membawa minimal 2 set kain ihram dan sabuk. Jamaah wanita membawa mukena serta pakaian gamis syar''i yang nyaman. Jangan lupa membawa sandal, kacamata hitam, pelembab bibir, serta dokumen paspor dan buku vaksin.', 'https://images.unsplash.com/photo-1542314831-068cd1dbfeeb?auto=format&fit=crop&w=800&q=80', ARRAY['perlengkapan', 'tips', 'checklist'], 'Perlengkapan Wajib Jamaah Umroh', 'Tips persiapan dan checklist perlengkapan ibadah umroh agar perjalanan lancar dan nyaman.', 'PUBLISHED', false, CURRENT_TIMESTAMP)
ON CONFLICT (app_id, slug) DO NOTHING;

-- 4. Seed Galleries
INSERT INTO content_galleries (id, app_id, title, slug, description, category, cover_url, is_active)
VALUES
('22222222-2222-2222-2222-222222222201', 'umroh', 'Dokumentasi Jamaah Umroh Musim 1447 H', 'dokumentasi-umroh-1447h', 'Koleksi momen kekhusyukan jamaah kami di Masjidil Haram dan Masjid Nabawi', 'UMRAH', 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?auto=format&fit=crop&w=800&q=80', true),
('22222222-2222-2222-2222-222222222202', 'umroh', 'Ziarah Sejarah Islam Kota Madinah', 'ziarah-sejarah-madinah', 'Momen kunjungan ke Masjid Quba, Jabal Uhud, dan Kebun Kurma', 'EVENT', 'https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?auto=format&fit=crop&w=800&q=80', true)
ON CONFLICT (app_id, slug) DO NOTHING;

-- 5. Seed Gallery Items
INSERT INTO content_gallery_items (app_id, gallery_id, title, caption, media_type, media_url, sequence_order, is_active)
VALUES
('umroh', '22222222-2222-2222-2222-222222222201', 'Tawaf di Depan Ka''bah', 'Suasana tawaf jamaah kloter Syawal di pelataran Ka''bah', 'IMAGE', 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?auto=format&fit=crop&w=800&q=80', 1, true),
('umroh', '22222222-2222-2222-2222-222222222201', 'Masjid Nabawi di Waktu Subuh', 'Keindahan payung hidrolik Masjid Nabawi Madinah', 'IMAGE', 'https://images.unsplash.com/photo-1542314831-068cd1dbfeeb?auto=format&fit=crop&w=800&q=80', 2, true),
('umroh', '22222222-2222-2222-2222-222222222202', 'Ziarah Syuhada Uhud', 'Doa bersama jamaah di pemakaman syuhada perang Uhud', 'IMAGE', 'https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?auto=format&fit=crop&w=800&q=80', 1, true)
ON CONFLICT DO NOTHING;

-- 6. Seed Custom Blocks
INSERT INTO content_custom_blocks (app_id, section_key, content_type, title, subtitle, body, image_url, metadata, sequence_order, is_active)
VALUES
('umroh', 'about-us', 'PAGE', 'Tentang Layanan Umroh & Haji Terpercaya', 'Komitmen Kami Melayani Tamu Allah', 'Kami adalah penyelenggara perjalanan ibadah umroh dan haji khusus resmi dengan komitmen memberikan bimbingan ibadah sesuai sunnah dan kenyamanan terbaik bagi seluruh jamaah.', 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?auto=format&fit=crop&w=800&q=80', '{"legal_kemenag": "PPIU No. 123/2020", "since_year": 2012, "jamaah_count": 25000}'::jsonb, 1, true),
('umroh', 'why-choose-us', 'SECTION', 'Mengapa Memilih Kami?', 'Kenyamanan & Kepastian Ibadah Anda', 'Kepastian jadwal terbang, hotel bintang 5 dekat masjid, muthawif bersertifikat, dan pendampingan medis.', '', '{"features": ["Pembimbing Sunnah", "Hotel Dekat Masjid", "Pasti Berangkat", "Layanan VIP"]}'::jsonb, 2, true)
ON CONFLICT (app_id, section_key) DO NOTHING;

-- 7. Seed FAQs
INSERT INTO content_faqs (app_id, category, question, answer, sequence_order, is_active)
VALUES
('umroh', 'PENDAFTARAN', 'Apa saja syarat dokumen untuk mendaftar paket umroh?', 'Syarat utama meliputi paspor asli masa berlaku minimal 7 bulan, KTP & Kartu Keluarga, buku nikah (bagi suami istri), dan sertifikat vaksin meningitis.', 1, true),
('umroh', 'PEMBAYARAN', 'Bagaimana skema pembayaran dan pelunasan paket?', 'Uang muka (DP) sebesar Rp 5.000.000 per jamaah saat pendaftaran. Pelunasan paling lambat 30 hari sebelum tanggal keberangkatan.', 2, true),
('umroh', 'KEBERANGKATAN', 'Apakah ada bimbingan manasik sebelum berangkat?', 'Ya, kami menyelenggarakan 2 kali manasik tatap muka di hotel bintang 4/5 dan bimbingan daring sebelum hari keberangkatan.', 3, true)
ON CONFLICT DO NOTHING;

-- 8. Seed Testimonials
INSERT INTO content_testimonials (app_id, customer_name, customer_title, avatar_url, rating, comment, package_name, is_featured, is_active)
VALUES
('umroh', 'H. Ahmad Fauzi', 'Jamaah Umroh Jakarta', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=200&q=80', 5, 'Alhamdulillah perjalanan ibadah umroh bersama keluarga sangat berkesan. Muthawif membimbing dengan sabar sesuai tuntunan sunnah dan hotelnya sangat dekat dengan pelataran Ka''bah.', 'Umroh Reguler Syawal Berkah 9 Hari', true, true),
('umroh', 'Hj. Siti Rahmah', 'Jamaah Umroh Bandung', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&w=200&q=80', 5, 'Fasilitas dan makanannya sangat terjamin. Pelayanan tim handling sejak di bandara CGK sampai hotel Jeddah/Madinah sangat profesional.', 'Umroh Plus Wisata Turki 12 Hari', true, true)
ON CONFLICT DO NOTHING;
