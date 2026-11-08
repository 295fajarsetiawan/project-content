package content

import (
	"context"
	"math"
	"strings"
	"time"

	"content-service/internal/domain/entity"
	domainErrors "content-service/internal/domain/errors"
	"content-service/internal/infrastructure/storage"

	"github.com/google/uuid"
)

type ContentUsecase interface {
	// Banners
	ListBanners(ctx context.Context, filter entity.ListBannersFilter) ([]*entity.ContentBanner, *entity.PaginationMeta, error)
	GetBanner(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBanner, error)
	CreateBanner(ctx context.Context, banner *entity.ContentBanner) (*entity.ContentBanner, error)
	UpdateBanner(ctx context.Context, banner *entity.ContentBanner) (*entity.ContentBanner, error)
	DeleteBanner(ctx context.Context, appID string, id uuid.UUID) error

	// Blog Categories
	ListBlogCategories(ctx context.Context, appID string) ([]*entity.ContentBlogCategory, error)
	CreateBlogCategory(ctx context.Context, category *entity.ContentBlogCategory) (*entity.ContentBlogCategory, error)

	// Blogs
	ListBlogs(ctx context.Context, filter entity.ListBlogsFilter) ([]*entity.ContentBlog, *entity.PaginationMeta, error)
	GetBlog(ctx context.Context, appID string, id *uuid.UUID, slug string) (*entity.ContentBlog, error)
	CreateBlog(ctx context.Context, blog *entity.ContentBlog) (*entity.ContentBlog, error)
	UpdateBlog(ctx context.Context, blog *entity.ContentBlog) (*entity.ContentBlog, error)
	DeleteBlog(ctx context.Context, appID string, id uuid.UUID) error

	// Galleries & Items
	ListGalleries(ctx context.Context, filter entity.ListGalleriesFilter) ([]*entity.ContentGallery, *entity.PaginationMeta, error)
	GetGallery(ctx context.Context, appID string, id *uuid.UUID, slug string) (*entity.ContentGallery, error)
	CreateGallery(ctx context.Context, gallery *entity.ContentGallery) (*entity.ContentGallery, error)
	DeleteGallery(ctx context.Context, appID string, id uuid.UUID) error
	AddGalleryItem(ctx context.Context, item *entity.ContentGalleryItem) (*entity.ContentGalleryItem, error)
	DeleteGalleryItem(ctx context.Context, appID string, id uuid.UUID) error

	// Custom Blocks
	ListCustomBlocks(ctx context.Context, appID, contentType string) ([]*entity.ContentCustomBlock, error)
	GetCustomBlock(ctx context.Context, appID, sectionKey string) (*entity.ContentCustomBlock, error)
	UpsertCustomBlock(ctx context.Context, block *entity.ContentCustomBlock) (*entity.ContentCustomBlock, error)
	DeleteCustomBlock(ctx context.Context, appID string, id uuid.UUID) error

	// FAQs
	ListFAQs(ctx context.Context, appID, category string) ([]*entity.ContentFAQ, error)
	CreateFAQ(ctx context.Context, faq *entity.ContentFAQ) (*entity.ContentFAQ, error)
	UpdateFAQ(ctx context.Context, faq *entity.ContentFAQ) (*entity.ContentFAQ, error)
	DeleteFAQ(ctx context.Context, appID string, id uuid.UUID) error

	// Testimonials
	ListTestimonials(ctx context.Context, filter entity.ListTestimonialsFilter) ([]*entity.ContentTestimonial, *entity.PaginationMeta, error)
	CreateTestimonial(ctx context.Context, testimonial *entity.ContentTestimonial) (*entity.ContentTestimonial, error)
	UpdateTestimonial(ctx context.Context, testimonial *entity.ContentTestimonial) (*entity.ContentTestimonial, error)
	DeleteTestimonial(ctx context.Context, appID string, id uuid.UUID) error

	// Media Storage
	UploadMedia(ctx context.Context, appID, folder, filename, base64Content string) (*storage.MediaUploadResult, error)
}

type contentUsecase struct {
	repo    entity.ContentRepository
	storage storage.StorageService
}

func NewContentUsecase(repo entity.ContentRepository, storage storage.StorageService) ContentUsecase {
	return &contentUsecase{
		repo:    repo,
		storage: storage,
	}
}

// -----------------------------------------------------------------------------
// 1. Banner Implementations
// -----------------------------------------------------------------------------

func (u *contentUsecase) ListBanners(ctx context.Context, filter entity.ListBannersFilter) ([]*entity.ContentBanner, *entity.PaginationMeta, error) {
	if filter.AppID == "" {
		filter.AppID = "umroh"
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	filter.Limit = limit
	filter.Offset = offset

	banners, total, err := u.repo.ListBanners(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	page := (offset / limit) + 1
	totalPages := int32(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	meta := &entity.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return banners, meta, nil
}

func (u *contentUsecase) GetBanner(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBanner, error) {
	if id == uuid.Nil {
		return nil, domainErrors.ErrBannerNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.GetBannerByID(ctx, appID, id)
}

func (u *contentUsecase) CreateBanner(ctx context.Context, b *entity.ContentBanner) (*entity.ContentBanner, error) {
	if b.Title == "" || b.ImageURL == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "title and image_url are required", nil)
	}
	if b.AppID == "" {
		b.AppID = "umroh"
	}
	if b.Position == "" {
		b.Position = "HERO"
	}
	b.Position = strings.ToUpper(b.Position)
	return u.repo.CreateBanner(ctx, b)
}

func (u *contentUsecase) UpdateBanner(ctx context.Context, b *entity.ContentBanner) (*entity.ContentBanner, error) {
	if b.ID == uuid.Nil {
		return nil, domainErrors.ErrBannerNotFound
	}
	if b.Title == "" || b.ImageURL == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "title and image_url are required", nil)
	}
	if b.AppID == "" {
		b.AppID = "umroh"
	}
	if b.Position == "" {
		b.Position = "HERO"
	}
	b.Position = strings.ToUpper(b.Position)
	return u.repo.UpdateBanner(ctx, b)
}

func (u *contentUsecase) DeleteBanner(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrBannerNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteBanner(ctx, appID, id)
}

// -----------------------------------------------------------------------------
// 2. Blog Categories & Blogs Implementations
// -----------------------------------------------------------------------------

func (u *contentUsecase) ListBlogCategories(ctx context.Context, appID string) ([]*entity.ContentBlogCategory, error) {
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.ListBlogCategories(ctx, appID)
}

func (u *contentUsecase) CreateBlogCategory(ctx context.Context, c *entity.ContentBlogCategory) (*entity.ContentBlogCategory, error) {
	if c.Name == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "name is required", nil)
	}
	if c.Slug == "" {
		c.Slug = slugify(c.Name)
	}
	if c.AppID == "" {
		c.AppID = "umroh"
	}
	c.Slug = strings.ToLower(strings.TrimSpace(c.Slug))
	c.IsActive = true
	return u.repo.CreateBlogCategory(ctx, c)
}

func (u *contentUsecase) ListBlogs(ctx context.Context, filter entity.ListBlogsFilter) ([]*entity.ContentBlog, *entity.PaginationMeta, error) {
	if filter.AppID == "" {
		filter.AppID = "umroh"
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	filter.Limit = limit
	filter.Offset = offset

	blogs, total, err := u.repo.ListBlogs(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	page := (offset / limit) + 1
	totalPages := int32(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	meta := &entity.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return blogs, meta, nil
}

func (u *contentUsecase) GetBlog(ctx context.Context, appID string, id *uuid.UUID, slug string) (*entity.ContentBlog, error) {
	if appID == "" {
		appID = "umroh"
	}

	var blog *entity.ContentBlog
	var err error

	if id != nil && *id != uuid.Nil {
		blog, err = u.repo.GetBlogByID(ctx, appID, *id)
	} else if slug != "" {
		blog, err = u.repo.GetBlogBySlug(ctx, appID, slug)
	} else {
		return nil, domainErrors.ErrBlogNotFound
	}

	if err != nil {
		return nil, err
	}

	// Increment view count asynchronously
	go func(bID uuid.UUID) {
		_ = u.repo.IncrementBlogViewCount(context.Background(), bID)
	}(blog.ID)

	return blog, nil
}

func (u *contentUsecase) CreateBlog(ctx context.Context, b *entity.ContentBlog) (*entity.ContentBlog, error) {
	if b.Title == "" || b.CategoryID == uuid.Nil {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "title and category_id are required", nil)
	}
	if b.Slug == "" {
		b.Slug = slugify(b.Title)
	}
	if b.AppID == "" {
		b.AppID = "umroh"
	}
	b.Slug = strings.ToLower(strings.TrimSpace(b.Slug))
	if b.Status == "" {
		b.Status = "PUBLISHED"
	}
	b.Status = strings.ToUpper(b.Status)
	if b.Status == "PUBLISHED" && b.PublishedAt == nil {
		now := time.Now()
		b.PublishedAt = &now
	}
	return u.repo.CreateBlog(ctx, b)
}

func (u *contentUsecase) UpdateBlog(ctx context.Context, b *entity.ContentBlog) (*entity.ContentBlog, error) {
	if b.ID == uuid.Nil {
		return nil, domainErrors.ErrBlogNotFound
	}
	if b.Title == "" || b.Slug == "" || b.CategoryID == uuid.Nil {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "title, slug, and category_id are required", nil)
	}
	if b.AppID == "" {
		b.AppID = "umroh"
	}
	b.Slug = strings.ToLower(strings.TrimSpace(b.Slug))
	b.Status = strings.ToUpper(b.Status)
	return u.repo.UpdateBlog(ctx, b)
}

func (u *contentUsecase) DeleteBlog(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrBlogNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteBlog(ctx, appID, id)
}

// -----------------------------------------------------------------------------
// 3. Galleries & Items Implementations
// -----------------------------------------------------------------------------

func (u *contentUsecase) ListGalleries(ctx context.Context, filter entity.ListGalleriesFilter) ([]*entity.ContentGallery, *entity.PaginationMeta, error) {
	if filter.AppID == "" {
		filter.AppID = "umroh"
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	filter.Limit = limit
	filter.Offset = offset

	galleries, total, err := u.repo.ListGalleries(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	page := (offset / limit) + 1
	totalPages := int32(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	meta := &entity.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return galleries, meta, nil
}

func (u *contentUsecase) GetGallery(ctx context.Context, appID string, id *uuid.UUID, slug string) (*entity.ContentGallery, error) {
	if appID == "" {
		appID = "umroh"
	}
	if id != nil && *id != uuid.Nil {
		return u.repo.GetGalleryByID(ctx, appID, *id)
	}
	if slug != "" {
		return u.repo.GetGalleryBySlug(ctx, appID, slug)
	}
	return nil, domainErrors.ErrGalleryNotFound
}

func (u *contentUsecase) CreateGallery(ctx context.Context, g *entity.ContentGallery) (*entity.ContentGallery, error) {
	if g.Title == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "title is required", nil)
	}
	if g.Slug == "" {
		g.Slug = slugify(g.Title)
	}
	if g.AppID == "" {
		g.AppID = "umroh"
	}
	if g.Category == "" {
		g.Category = "GENERAL"
	}
	g.Slug = strings.ToLower(strings.TrimSpace(g.Slug))
	g.Category = strings.ToUpper(g.Category)
	g.IsActive = true
	return u.repo.CreateGallery(ctx, g)
}

func (u *contentUsecase) DeleteGallery(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrGalleryNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteGallery(ctx, appID, id)
}

func (u *contentUsecase) AddGalleryItem(ctx context.Context, it *entity.ContentGalleryItem) (*entity.ContentGalleryItem, error) {
	if it.GalleryID == uuid.Nil || it.MediaURL == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "gallery_id and media_url are required", nil)
	}
	if it.AppID == "" {
		it.AppID = "umroh"
	}
	if it.MediaType == "" {
		it.MediaType = "IMAGE"
	}
	it.MediaType = strings.ToUpper(it.MediaType)
	it.IsActive = true
	return u.repo.AddGalleryItem(ctx, it)
}

func (u *contentUsecase) DeleteGalleryItem(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrGalleryItemNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteGalleryItem(ctx, appID, id)
}

// -----------------------------------------------------------------------------
// 4. Custom Blocks Implementations
// -----------------------------------------------------------------------------

func (u *contentUsecase) ListCustomBlocks(ctx context.Context, appID, contentType string) ([]*entity.ContentCustomBlock, error) {
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.ListCustomBlocks(ctx, appID, contentType, nil)
}

func (u *contentUsecase) GetCustomBlock(ctx context.Context, appID, sectionKey string) (*entity.ContentCustomBlock, error) {
	if appID == "" {
		appID = "umroh"
	}
	if sectionKey == "" {
		return nil, domainErrors.ErrCustomBlockNotFound
	}
	return u.repo.GetCustomBlockBySection(ctx, appID, sectionKey)
}

func (u *contentUsecase) UpsertCustomBlock(ctx context.Context, b *entity.ContentCustomBlock) (*entity.ContentCustomBlock, error) {
	if b.SectionKey == "" || b.Title == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "section_key and title are required", nil)
	}
	if b.AppID == "" {
		b.AppID = "umroh"
	}
	if b.ContentType == "" {
		b.ContentType = "SECTION"
	}
	b.SectionKey = strings.ToLower(strings.TrimSpace(b.SectionKey))
	b.ContentType = strings.ToUpper(b.ContentType)
	b.IsActive = true
	return u.repo.UpsertCustomBlock(ctx, b)
}

func (u *contentUsecase) DeleteCustomBlock(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrCustomBlockNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteCustomBlock(ctx, appID, id)
}

// -----------------------------------------------------------------------------
// 5. FAQs Implementations
// -----------------------------------------------------------------------------

func (u *contentUsecase) ListFAQs(ctx context.Context, appID, category string) ([]*entity.ContentFAQ, error) {
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.ListFAQs(ctx, appID, category, nil)
}

func (u *contentUsecase) CreateFAQ(ctx context.Context, f *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	if f.Question == "" || f.Answer == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "question and answer are required", nil)
	}
	if f.AppID == "" {
		f.AppID = "umroh"
	}
	if f.Category == "" {
		f.Category = "UMUM"
	}
	f.Category = strings.ToUpper(f.Category)
	f.IsActive = true
	return u.repo.CreateFAQ(ctx, f)
}

func (u *contentUsecase) UpdateFAQ(ctx context.Context, f *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	if f.ID == uuid.Nil {
		return nil, domainErrors.ErrFAQNotFound
	}
	if f.Question == "" || f.Answer == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "question and answer are required", nil)
	}
	if f.AppID == "" {
		f.AppID = "umroh"
	}
	if f.Category == "" {
		f.Category = "UMUM"
	}
	f.Category = strings.ToUpper(f.Category)
	return u.repo.UpdateFAQ(ctx, f)
}

func (u *contentUsecase) DeleteFAQ(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrFAQNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteFAQ(ctx, appID, id)
}

// -----------------------------------------------------------------------------
// 6. Testimonials Implementations
// -----------------------------------------------------------------------------

func (u *contentUsecase) ListTestimonials(ctx context.Context, filter entity.ListTestimonialsFilter) ([]*entity.ContentTestimonial, *entity.PaginationMeta, error) {
	if filter.AppID == "" {
		filter.AppID = "umroh"
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	filter.Limit = limit
	filter.Offset = offset

	testimonials, total, err := u.repo.ListTestimonials(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	page := (offset / limit) + 1
	totalPages := int32(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	meta := &entity.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return testimonials, meta, nil
}

func (u *contentUsecase) CreateTestimonial(ctx context.Context, t *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	if t.CustomerName == "" || t.Comment == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "customer_name and comment are required", nil)
	}
	if t.Rating < 1 {
		t.Rating = 1
	}
	if t.Rating > 5 {
		t.Rating = 5
	}
	if t.AppID == "" {
		t.AppID = "umroh"
	}
	t.IsActive = true
	return u.repo.CreateTestimonial(ctx, t)
}

func (u *contentUsecase) UpdateTestimonial(ctx context.Context, t *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	if t.ID == uuid.Nil {
		return nil, domainErrors.ErrTestimonialNotFound
	}
	if t.CustomerName == "" || t.Comment == "" {
		return nil, domainErrors.NewDomainError("INVALID_INPUT", "customer_name and comment are required", nil)
	}
	if t.Rating < 1 {
		t.Rating = 1
	}
	if t.Rating > 5 {
		t.Rating = 5
	}
	if t.AppID == "" {
		t.AppID = "umroh"
	}
	return u.repo.UpdateTestimonial(ctx, t)
}

func (u *contentUsecase) DeleteTestimonial(ctx context.Context, appID string, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainErrors.ErrTestimonialNotFound
	}
	if appID == "" {
		appID = "umroh"
	}
	return u.repo.DeleteTestimonial(ctx, appID, id)
}

// -----------------------------------------------------------------------------
// 7. Media Storage Implementation
// -----------------------------------------------------------------------------

func (u *contentUsecase) UploadMedia(ctx context.Context, appID, folder, filename, base64Content string) (*storage.MediaUploadResult, error) {
	if u.storage == nil {
		return nil, domainErrors.NewDomainError("STORAGE_UNAVAILABLE", "storage service is not configured", nil)
	}
	if appID == "" {
		appID = "umroh"
	}
	if folder == "" {
		folder = "general"
	}
	return u.storage.SaveBase64Media(ctx, appID, folder, filename, base64Content)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			if b.Len() > 0 && !strings.HasSuffix(b.String(), "-") {
				b.WriteRune('-')
			}
		}
	}
	res := strings.Trim(b.String(), "-")
	if res == "" {
		res = "content"
	}
	return res
}
