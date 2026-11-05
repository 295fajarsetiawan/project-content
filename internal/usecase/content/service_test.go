package content_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"content-service/internal/domain/entity"
	domainErrors "content-service/internal/domain/errors"
	"content-service/internal/infrastructure/storage"
	"content-service/internal/usecase/content"

	"github.com/google/uuid"
)

type mockContentRepo struct {
	banners      map[string]*entity.ContentBanner
	categories   map[string]*entity.ContentBlogCategory
	blogs        map[string]*entity.ContentBlog
	galleries    map[string]*entity.ContentGallery
	galleryItems map[uuid.UUID]*entity.ContentGalleryItem
	customBlocks map[string]*entity.ContentCustomBlock
	faqs         map[string]*entity.ContentFAQ
	testimonials map[string]*entity.ContentTestimonial
}

func newMockRepo() *mockContentRepo {
	return &mockContentRepo{
		banners:      make(map[string]*entity.ContentBanner),
		categories:   make(map[string]*entity.ContentBlogCategory),
		blogs:        make(map[string]*entity.ContentBlog),
		galleries:    make(map[string]*entity.ContentGallery),
		galleryItems: make(map[uuid.UUID]*entity.ContentGalleryItem),
		customBlocks: make(map[string]*entity.ContentCustomBlock),
		faqs:         make(map[string]*entity.ContentFAQ),
		testimonials: make(map[string]*entity.ContentTestimonial),
	}
}

func (m *mockContentRepo) WithTx(ctx context.Context, fn func(repo entity.ContentRepository) error) error {
	return fn(m)
}

func (m *mockContentRepo) ListBanners(ctx context.Context, filter entity.ListBannersFilter) ([]*entity.ContentBanner, int64, error) {
	var list []*entity.ContentBanner
	for _, b := range m.banners {
		if b.AppID == filter.AppID {
			list = append(list, b)
		}
	}
	return list, int64(len(list)), nil
}

func (m *mockContentRepo) GetBannerByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBanner, error) {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	b, ok := m.banners[key]
	if !ok {
		return nil, domainErrors.ErrBannerNotFound
	}
	return b, nil
}

func (m *mockContentRepo) CreateBanner(ctx context.Context, banner *entity.ContentBanner) (*entity.ContentBanner, error) {
	if banner.ID == uuid.Nil {
		banner.ID = uuid.New()
	}
	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()
	key := fmt.Sprintf("%s:%s", banner.AppID, banner.ID.String())
	m.banners[key] = banner
	return banner, nil
}

func (m *mockContentRepo) UpdateBanner(ctx context.Context, banner *entity.ContentBanner) (*entity.ContentBanner, error) {
	key := fmt.Sprintf("%s:%s", banner.AppID, banner.ID.String())
	if _, ok := m.banners[key]; !ok {
		return nil, domainErrors.ErrBannerNotFound
	}
	banner.UpdatedAt = time.Now()
	m.banners[key] = banner
	return banner, nil
}

func (m *mockContentRepo) DeleteBanner(ctx context.Context, appID string, id uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	delete(m.banners, key)
	return nil
}

func (m *mockContentRepo) ListBlogCategories(ctx context.Context, appID string) ([]*entity.ContentBlogCategory, error) {
	var list []*entity.ContentBlogCategory
	for _, c := range m.categories {
		if c.AppID == appID {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *mockContentRepo) GetBlogCategoryByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBlogCategory, error) {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	c, ok := m.categories[key]
	if !ok {
		return nil, domainErrors.ErrBlogCategoryNotFound
	}
	return c, nil
}

func (m *mockContentRepo) GetBlogCategoryBySlug(ctx context.Context, appID, slug string) (*entity.ContentBlogCategory, error) {
	for _, c := range m.categories {
		if c.AppID == appID && c.Slug == slug {
			return c, nil
		}
	}
	return nil, domainErrors.ErrBlogCategoryNotFound
}

func (m *mockContentRepo) CreateBlogCategory(ctx context.Context, category *entity.ContentBlogCategory) (*entity.ContentBlogCategory, error) {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	key := fmt.Sprintf("%s:%s", category.AppID, category.ID.String())
	m.categories[key] = category
	return category, nil
}

func (m *mockContentRepo) ListBlogs(ctx context.Context, filter entity.ListBlogsFilter) ([]*entity.ContentBlog, int64, error) {
	var list []*entity.ContentBlog
	for _, b := range m.blogs {
		if b.AppID == filter.AppID {
			list = append(list, b)
		}
	}
	return list, int64(len(list)), nil
}

func (m *mockContentRepo) GetBlogByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBlog, error) {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	b, ok := m.blogs[key]
	if !ok {
		return nil, domainErrors.ErrBlogNotFound
	}
	return b, nil
}

func (m *mockContentRepo) GetBlogBySlug(ctx context.Context, appID, slug string) (*entity.ContentBlog, error) {
	for _, b := range m.blogs {
		if b.AppID == appID && b.Slug == slug {
			return b, nil
		}
	}
	return nil, domainErrors.ErrBlogNotFound
}

func (m *mockContentRepo) IncrementBlogViewCount(ctx context.Context, id uuid.UUID) error {
	for _, b := range m.blogs {
		if b.ID == id {
			b.ViewCount++
			return nil
		}
	}
	return nil
}

func (m *mockContentRepo) CreateBlog(ctx context.Context, blog *entity.ContentBlog) (*entity.ContentBlog, error) {
	if blog.ID == uuid.Nil {
		blog.ID = uuid.New()
	}
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	key := fmt.Sprintf("%s:%s", blog.AppID, blog.ID.String())
	m.blogs[key] = blog
	return blog, nil
}

func (m *mockContentRepo) UpdateBlog(ctx context.Context, blog *entity.ContentBlog) (*entity.ContentBlog, error) {
	key := fmt.Sprintf("%s:%s", blog.AppID, blog.ID.String())
	if _, ok := m.blogs[key]; !ok {
		return nil, domainErrors.ErrBlogNotFound
	}
	m.blogs[key] = blog
	return blog, nil
}

func (m *mockContentRepo) DeleteBlog(ctx context.Context, appID string, id uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	delete(m.blogs, key)
	return nil
}

func (m *mockContentRepo) ListGalleries(ctx context.Context, filter entity.ListGalleriesFilter) ([]*entity.ContentGallery, int64, error) {
	var list []*entity.ContentGallery
	for _, g := range m.galleries {
		if g.AppID == filter.AppID {
			list = append(list, g)
		}
	}
	return list, int64(len(list)), nil
}

func (m *mockContentRepo) GetGalleryByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentGallery, error) {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	g, ok := m.galleries[key]
	if !ok {
		return nil, domainErrors.ErrGalleryNotFound
	}
	return g, nil
}

func (m *mockContentRepo) GetGalleryBySlug(ctx context.Context, appID, slug string) (*entity.ContentGallery, error) {
	for _, g := range m.galleries {
		if g.AppID == appID && g.Slug == slug {
			return g, nil
		}
	}
	return nil, domainErrors.ErrGalleryNotFound
}

func (m *mockContentRepo) CreateGallery(ctx context.Context, gallery *entity.ContentGallery) (*entity.ContentGallery, error) {
	if gallery.ID == uuid.Nil {
		gallery.ID = uuid.New()
	}
	key := fmt.Sprintf("%s:%s", gallery.AppID, gallery.ID.String())
	m.galleries[key] = gallery
	return gallery, nil
}

func (m *mockContentRepo) DeleteGallery(ctx context.Context, appID string, id uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	delete(m.galleries, key)
	return nil
}

func (m *mockContentRepo) ListGalleryItems(ctx context.Context, galleryID uuid.UUID, isActive *bool) ([]entity.ContentGalleryItem, error) {
	var list []entity.ContentGalleryItem
	for _, it := range m.galleryItems {
		if it.GalleryID == galleryID {
			list = append(list, *it)
		}
	}
	return list, nil
}

func (m *mockContentRepo) AddGalleryItem(ctx context.Context, item *entity.ContentGalleryItem) (*entity.ContentGalleryItem, error) {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	m.galleryItems[item.ID] = item
	return item, nil
}

func (m *mockContentRepo) DeleteGalleryItem(ctx context.Context, id uuid.UUID) error {
	delete(m.galleryItems, id)
	return nil
}

func (m *mockContentRepo) ListCustomBlocks(ctx context.Context, appID, contentType string, isActive *bool) ([]*entity.ContentCustomBlock, error) {
	var list []*entity.ContentCustomBlock
	for _, b := range m.customBlocks {
		if b.AppID == appID {
			list = append(list, b)
		}
	}
	return list, nil
}

func (m *mockContentRepo) GetCustomBlockBySection(ctx context.Context, appID, sectionKey string) (*entity.ContentCustomBlock, error) {
	key := fmt.Sprintf("%s:%s", appID, sectionKey)
	b, ok := m.customBlocks[key]
	if !ok {
		return nil, domainErrors.ErrCustomBlockNotFound
	}
	return b, nil
}

func (m *mockContentRepo) UpsertCustomBlock(ctx context.Context, block *entity.ContentCustomBlock) (*entity.ContentCustomBlock, error) {
	if block.ID == uuid.Nil {
		block.ID = uuid.New()
	}
	key := fmt.Sprintf("%s:%s", block.AppID, block.SectionKey)
	m.customBlocks[key] = block
	return block, nil
}

func (m *mockContentRepo) DeleteCustomBlock(ctx context.Context, appID string, id uuid.UUID) error {
	for k, b := range m.customBlocks {
		if b.AppID == appID && b.ID == id {
			delete(m.customBlocks, k)
			return nil
		}
	}
	return nil
}

func (m *mockContentRepo) ListFAQs(ctx context.Context, appID, category string, isActive *bool) ([]*entity.ContentFAQ, error) {
	var list []*entity.ContentFAQ
	for _, f := range m.faqs {
		if f.AppID == appID {
			list = append(list, f)
		}
	}
	return list, nil
}

func (m *mockContentRepo) CreateFAQ(ctx context.Context, faq *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	if faq.ID == uuid.Nil {
		faq.ID = uuid.New()
	}
	key := fmt.Sprintf("%s:%s", faq.AppID, faq.ID.String())
	m.faqs[key] = faq
	return faq, nil
}

func (m *mockContentRepo) UpdateFAQ(ctx context.Context, faq *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	key := fmt.Sprintf("%s:%s", faq.AppID, faq.ID.String())
	if _, ok := m.faqs[key]; !ok {
		return nil, domainErrors.ErrFAQNotFound
	}
	m.faqs[key] = faq
	return faq, nil
}

func (m *mockContentRepo) DeleteFAQ(ctx context.Context, appID string, id uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	delete(m.faqs, key)
	return nil
}

func (m *mockContentRepo) ListTestimonials(ctx context.Context, filter entity.ListTestimonialsFilter) ([]*entity.ContentTestimonial, int64, error) {
	var list []*entity.ContentTestimonial
	for _, t := range m.testimonials {
		if t.AppID == filter.AppID {
			list = append(list, t)
		}
	}
	return list, int64(len(list)), nil
}

func (m *mockContentRepo) CreateTestimonial(ctx context.Context, testimonial *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	if testimonial.ID == uuid.Nil {
		testimonial.ID = uuid.New()
	}
	key := fmt.Sprintf("%s:%s", testimonial.AppID, testimonial.ID.String())
	m.testimonials[key] = testimonial
	return testimonial, nil
}

func (m *mockContentRepo) UpdateTestimonial(ctx context.Context, testimonial *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	key := fmt.Sprintf("%s:%s", testimonial.AppID, testimonial.ID.String())
	if _, ok := m.testimonials[key]; !ok {
		return nil, domainErrors.ErrTestimonialNotFound
	}
	m.testimonials[key] = testimonial
	return testimonial, nil
}

func (m *mockContentRepo) DeleteTestimonial(ctx context.Context, appID string, id uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", appID, id.String())
	delete(m.testimonials, key)
	return nil
}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func TestContentUsecase_Banners(t *testing.T) {
	repo := newMockRepo()
	storageSvc := storage.NewLocalStorageService(t.TempDir(), "http://localhost:50055/storage/uploads")
	uc := content.NewContentUsecase(repo, storageSvc)
	ctx := context.Background()

	// 1. Validation error when title is empty
	_, err := uc.CreateBanner(ctx, &entity.ContentBanner{
		AppID: "umroh",
	})
	if err == nil {
		t.Fatalf("expected error when title is empty, got nil")
	}

	// 2. Success create banner
	b, err := uc.CreateBanner(ctx, &entity.ContentBanner{
		AppID:    "umroh",
		Title:    "Promo Umroh Ramadhan",
		ImageURL: "https://example.com/promo.jpg",
		Position: "HERO",
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("create banner failed: %v", err)
	}

	// 3. Get banner
	got, err := uc.GetBanner(ctx, "umroh", b.ID)
	if err != nil {
		t.Fatalf("get banner failed: %v", err)
	}
	if got.Title != "Promo Umroh Ramadhan" {
		t.Errorf("expected title 'Promo Umroh Ramadhan', got %s", got.Title)
	}

	// 4. List banners
	list, meta, err := uc.ListBanners(ctx, entity.ListBannersFilter{
		AppID: "umroh",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("list banners failed: %v", err)
	}
	if len(list) != 1 || meta.TotalItems != 1 {
		t.Errorf("expected 1 banner, got %d", len(list))
	}

	// 5. Delete banner
	if err := uc.DeleteBanner(ctx, "umroh", b.ID); err != nil {
		t.Fatalf("delete banner failed: %v", err)
	}
}

func TestContentUsecase_Blogs(t *testing.T) {
	repo := newMockRepo()
	storageSvc := storage.NewLocalStorageService(t.TempDir(), "http://localhost:50055/storage/uploads")
	uc := content.NewContentUsecase(repo, storageSvc)
	ctx := context.Background()

	// 1. Create Category
	cat, err := uc.CreateBlogCategory(ctx, &entity.ContentBlogCategory{
		AppID: "umroh",
		Name:  "Tips Ibadah",
	})
	if err != nil {
		t.Fatalf("create category failed: %v", err)
	}
	if cat.Slug != "tips-ibadah" {
		t.Errorf("expected slug 'tips-ibadah', got %s", cat.Slug)
	}

	// 2. Create Blog
	blog, err := uc.CreateBlog(ctx, &entity.ContentBlog{
		AppID:      "umroh",
		CategoryID: cat.ID,
		Title:      "Panduan Lengkap Tawaf",
		Body:       "<p>Tawaf adalah...</p>",
	})
	if err != nil {
		t.Fatalf("create blog failed: %v", err)
	}
	if blog.Slug != "panduan-lengkap-tawaf" {
		t.Errorf("expected slug 'panduan-lengkap-tawaf', got %s", blog.Slug)
	}

	// 3. Get Blog by Slug
	got, err := uc.GetBlog(ctx, "umroh", nil, "panduan-lengkap-tawaf")
	if err != nil {
		t.Fatalf("get blog by slug failed: %v", err)
	}
	if got.ID != blog.ID {
		t.Errorf("expected id %v, got %v", blog.ID, got.ID)
	}
}

func TestContentUsecase_CustomBlocksAndFAQs(t *testing.T) {
	repo := newMockRepo()
	storageSvc := storage.NewLocalStorageService(t.TempDir(), "http://localhost:50055/storage/uploads")
	uc := content.NewContentUsecase(repo, storageSvc)
	ctx := context.Background()

	// 1. Custom block upsert
	block, err := uc.UpsertCustomBlock(ctx, &entity.ContentCustomBlock{
		AppID:       "umroh",
		SectionKey:  "tentang-kami",
		ContentType: "PAGE",
		Title:       "Tentang Kami",
		Body:        "Kami adalah penyedia layanan travel terpercaya.",
	})
	if err != nil {
		t.Fatalf("upsert custom block failed: %v", err)
	}

	gotBlock, err := uc.GetCustomBlock(ctx, "umroh", "tentang-kami")
	if err != nil {
		t.Fatalf("get custom block failed: %v", err)
	}
	if gotBlock.Title != block.Title {
		t.Errorf("expected title %s, got %s", block.Title, gotBlock.Title)
	}

	// 2. FAQ
	faq, err := uc.CreateFAQ(ctx, &entity.ContentFAQ{
		AppID:    "umroh",
		Category: "UMRAH",
		Question: "Kapan waktu terbaik untuk mendaftar?",
		Answer:   "Minimal 3 bulan sebelum jadwal keberangkatan.",
	})
	if err != nil {
		t.Fatalf("create faq failed: %v", err)
	}

	faqs, err := uc.ListFAQs(ctx, "umroh", "UMRAH")
	if err != nil {
		t.Fatalf("list faqs failed: %v", err)
	}
	if len(faqs) != 1 || faqs[0].ID != faq.ID {
		t.Errorf("expected 1 faq matching, got %d", len(faqs))
	}
}
