package handler_test

import (
	"context"
	"testing"
	"time"

	"content-service/internal/delivery/grpc/handler"
	"content-service/internal/delivery/grpc/pb"
	"content-service/internal/domain/entity"
	domainErrors "content-service/internal/domain/errors"
	"content-service/internal/infrastructure/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type mockContentUsecase struct {
	banner    *entity.ContentBanner
	blog      *entity.ContentBlog
	gallery   *entity.ContentGallery
	block     *entity.ContentCustomBlock
	faq       *entity.ContentFAQ
	testimoni *entity.ContentTestimonial
}

func (m *mockContentUsecase) ListBanners(ctx context.Context, filter entity.ListBannersFilter) ([]*entity.ContentBanner, *entity.PaginationMeta, error) {
	if m.banner != nil {
		return []*entity.ContentBanner{m.banner}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 1, TotalPages: 1}, nil
	}
	return []*entity.ContentBanner{}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 0, TotalPages: 0}, nil
}

func (m *mockContentUsecase) GetBanner(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBanner, error) {
	if m.banner != nil && m.banner.ID == id {
		return m.banner, nil
	}
	return nil, domainErrors.ErrBannerNotFound
}

func (m *mockContentUsecase) CreateBanner(ctx context.Context, banner *entity.ContentBanner) (*entity.ContentBanner, error) {
	banner.ID = uuid.New()
	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()
	m.banner = banner
	return banner, nil
}

func (m *mockContentUsecase) UpdateBanner(ctx context.Context, banner *entity.ContentBanner) (*entity.ContentBanner, error) {
	banner.UpdatedAt = time.Now()
	m.banner = banner
	return banner, nil
}

func (m *mockContentUsecase) DeleteBanner(ctx context.Context, appID string, id uuid.UUID) error {
	m.banner = nil
	return nil
}

func (m *mockContentUsecase) ListBlogCategories(ctx context.Context, appID string) ([]*entity.ContentBlogCategory, error) {
	return []*entity.ContentBlogCategory{
		{ID: uuid.New(), AppID: appID, Name: "Tips", Slug: "tips", IsActive: true},
	}, nil
}

func (m *mockContentUsecase) CreateBlogCategory(ctx context.Context, category *entity.ContentBlogCategory) (*entity.ContentBlogCategory, error) {
	category.ID = uuid.New()
	category.CreatedAt = time.Now()
	return category, nil
}

func (m *mockContentUsecase) ListBlogs(ctx context.Context, filter entity.ListBlogsFilter) ([]*entity.ContentBlog, *entity.PaginationMeta, error) {
	if m.blog != nil {
		return []*entity.ContentBlog{m.blog}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 1, TotalPages: 1}, nil
	}
	return []*entity.ContentBlog{}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 0, TotalPages: 0}, nil
}

func (m *mockContentUsecase) GetBlog(ctx context.Context, appID string, id *uuid.UUID, slug string) (*entity.ContentBlog, error) {
	if m.blog != nil && (slug == m.blog.Slug || (id != nil && *id == m.blog.ID)) {
		return m.blog, nil
	}
	return nil, domainErrors.ErrBlogNotFound
}

func (m *mockContentUsecase) CreateBlog(ctx context.Context, blog *entity.ContentBlog) (*entity.ContentBlog, error) {
	blog.ID = uuid.New()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	m.blog = blog
	return blog, nil
}

func (m *mockContentUsecase) UpdateBlog(ctx context.Context, blog *entity.ContentBlog) (*entity.ContentBlog, error) {
	blog.UpdatedAt = time.Now()
	m.blog = blog
	return blog, nil
}

func (m *mockContentUsecase) DeleteBlog(ctx context.Context, appID string, id uuid.UUID) error {
	m.blog = nil
	return nil
}

func (m *mockContentUsecase) ListGalleries(ctx context.Context, filter entity.ListGalleriesFilter) ([]*entity.ContentGallery, *entity.PaginationMeta, error) {
	if m.gallery != nil {
		return []*entity.ContentGallery{m.gallery}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 1, TotalPages: 1}, nil
	}
	return []*entity.ContentGallery{}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 0, TotalPages: 0}, nil
}

func (m *mockContentUsecase) GetGallery(ctx context.Context, appID string, id *uuid.UUID, slug string) (*entity.ContentGallery, error) {
	if m.gallery != nil && (slug == m.gallery.Slug || (id != nil && *id == m.gallery.ID)) {
		return m.gallery, nil
	}
	return nil, domainErrors.ErrGalleryNotFound
}

func (m *mockContentUsecase) CreateGallery(ctx context.Context, gallery *entity.ContentGallery) (*entity.ContentGallery, error) {
	gallery.ID = uuid.New()
	m.gallery = gallery
	return gallery, nil
}

func (m *mockContentUsecase) DeleteGallery(ctx context.Context, appID string, id uuid.UUID) error {
	m.gallery = nil
	return nil
}

func (m *mockContentUsecase) AddGalleryItem(ctx context.Context, item *entity.ContentGalleryItem) (*entity.ContentGalleryItem, error) {
	item.ID = uuid.New()
	return item, nil
}

func (m *mockContentUsecase) DeleteGalleryItem(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockContentUsecase) ListCustomBlocks(ctx context.Context, appID, contentType string) ([]*entity.ContentCustomBlock, error) {
	if m.block != nil {
		return []*entity.ContentCustomBlock{m.block}, nil
	}
	return []*entity.ContentCustomBlock{}, nil
}

func (m *mockContentUsecase) GetCustomBlock(ctx context.Context, appID, sectionKey string) (*entity.ContentCustomBlock, error) {
	if m.block != nil && m.block.SectionKey == sectionKey {
		return m.block, nil
	}
	return nil, domainErrors.ErrCustomBlockNotFound
}

func (m *mockContentUsecase) UpsertCustomBlock(ctx context.Context, block *entity.ContentCustomBlock) (*entity.ContentCustomBlock, error) {
	block.ID = uuid.New()
	m.block = block
	return block, nil
}

func (m *mockContentUsecase) DeleteCustomBlock(ctx context.Context, appID string, id uuid.UUID) error {
	m.block = nil
	return nil
}

func (m *mockContentUsecase) ListFAQs(ctx context.Context, appID, category string) ([]*entity.ContentFAQ, error) {
	if m.faq != nil {
		return []*entity.ContentFAQ{m.faq}, nil
	}
	return []*entity.ContentFAQ{}, nil
}

func (m *mockContentUsecase) CreateFAQ(ctx context.Context, faq *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	faq.ID = uuid.New()
	m.faq = faq
	return faq, nil
}

func (m *mockContentUsecase) UpdateFAQ(ctx context.Context, faq *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	m.faq = faq
	return faq, nil
}

func (m *mockContentUsecase) DeleteFAQ(ctx context.Context, appID string, id uuid.UUID) error {
	m.faq = nil
	return nil
}

func (m *mockContentUsecase) ListTestimonials(ctx context.Context, filter entity.ListTestimonialsFilter) ([]*entity.ContentTestimonial, *entity.PaginationMeta, error) {
	if m.testimoni != nil {
		return []*entity.ContentTestimonial{m.testimoni}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 1, TotalPages: 1}, nil
	}
	return []*entity.ContentTestimonial{}, &entity.PaginationMeta{CurrentPage: 1, PerPage: 10, TotalItems: 0, TotalPages: 0}, nil
}

func (m *mockContentUsecase) CreateTestimonial(ctx context.Context, testimonial *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	testimonial.ID = uuid.New()
	m.testimoni = testimonial
	return testimonial, nil
}

func (m *mockContentUsecase) UpdateTestimonial(ctx context.Context, testimonial *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	m.testimoni = testimonial
	return testimonial, nil
}

func (m *mockContentUsecase) DeleteTestimonial(ctx context.Context, appID string, id uuid.UUID) error {
	m.testimoni = nil
	return nil
}

func (m *mockContentUsecase) UploadMedia(ctx context.Context, appID, folder, filename, base64Content string) (*storage.MediaUploadResult, error) {
	return &storage.MediaUploadResult{
		FileURL:  "http://localhost:50055/storage/uploads/umroh/banners/test.png",
		FilePath: "storage/uploads/umroh/banners/test.png",
		Filename: "test.png",
		FileSize: 1234,
		MimeType: "image/png",
	}, nil
}

func TestContentHandler_CreateAndListBanners(t *testing.T) {
	uc := &mockContentUsecase{}
	h := handler.NewContentHandler(uc)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-app-id", "umroh"))

	// Create
	createResp, err := h.CreateBanner(ctx, &pb.CreateBannerRequest{
		Title:    "Ramadhan Promo",
		ImageUrl: "https://example.com/banner.jpg",
		Position: "HERO",
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createResp.Data.Title != "Ramadhan Promo" {
		t.Errorf("expected Ramadhan Promo, got %s", createResp.Data.Title)
	}

	// List
	listResp, err := h.ListBanners(ctx, &pb.ListBannersRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(listResp.Data) != 1 {
		t.Errorf("expected 1 banner, got %d", len(listResp.Data))
	}
	if listResp.Pagination.TotalItems != 1 {
		t.Errorf("expected total_items 1, got %d", listResp.Pagination.TotalItems)
	}

	// Get Not Found
	_, err = h.GetBanner(ctx, &pb.GetBannerRequest{Id: uuid.New().String()})
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound error, got %v", status.Code(err))
	}
}

func TestContentHandler_UploadMedia(t *testing.T) {
	uc := &mockContentUsecase{}
	h := handler.NewContentHandler(uc)

	ctx := context.Background()
	resp, err := h.UploadMedia(ctx, &pb.UploadMediaRequest{
		Filename:      "test.png",
		Base64Content: "data:image/png;base64,abc",
		Folder:        "banners",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Data.MimeType != "image/png" {
		t.Errorf("expected image/png, got %s", resp.Data.MimeType)
	}
}
