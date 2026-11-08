package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PaginationMeta encapsulates standard pagination details
type PaginationMeta struct {
	CurrentPage int32 `json:"current_page"`
	PerPage     int32 `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int32 `json:"total_pages"`
}

// NormalizePage returns a bounded page size and non-negative row offset.
func NormalizePage(limit, offset int32) (int32, int32) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// NormalizePageRequest converts a 1-based page request into a bounded limit and offset.
func NormalizePageRequest(page, perPage int32) (int32, int32) {
	if page < 1 {
		page = 1
	}
	limit, _ := NormalizePage(perPage, 0)
	offset := (int64(page) - 1) * int64(limit)
	maxOffset := int64(1<<31 - 1 - (1<<31-1)%limit)
	if offset > maxOffset {
		offset = maxOffset
	}
	return limit, int32(offset)
}

// -----------------------------------------------------------------------------
// 1. Content Banner
// -----------------------------------------------------------------------------

type ContentBanner struct {
	ID             uuid.UUID              `json:"id"`
	AppID          string                 `json:"app_id"`
	Title          string                 `json:"title"`
	Subtitle       string                 `json:"subtitle"`
	Description    string                 `json:"description"`
	ImageURL       string                 `json:"image_url"`
	MobileImageURL string                 `json:"mobile_image_url"`
	TargetURL      string                 `json:"target_url"`
	ButtonText     string                 `json:"button_text"`
	Position       string                 `json:"position"` // HERO, PROMO, POPUP, SIDEBAR
	SequenceOrder  int32                  `json:"sequence_order"`
	StartDate      *time.Time             `json:"start_date,omitempty"`
	EndDate        *time.Time             `json:"end_date,omitempty"`
	IsActive       bool                   `json:"is_active"`
	AdditionalData map[string]interface{} `json:"additional_data"`
	CreatedBy      *uuid.UUID             `json:"created_by,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type ListBlogCategoriesFilter struct {
	AppID    string
	Search   string
	IsActive *bool
	Limit    int32
	Offset   int32
}

const (
	BlogStatusDraft     = "DRAFT"
	BlogStatusPublished = "PUBLISHED"
	BlogStatusArchived  = "ARCHIVED"
)

type ListCustomBlocksFilter struct {
	AppID       string
	ContentType string
	IsActive    *bool
	Search      string
	Limit       int32
	Offset      int32
}

type ListFAQsFilter struct {
	AppID    string
	Category string
	IsActive *bool
	Search   string
	Limit    int32
	Offset   int32
}

type ListBannersFilter struct {
	AppID    string
	Position string
	IsActive *bool
	Limit    int32
	Offset   int32
}

// -----------------------------------------------------------------------------
// 2. Content Blog & Categories
// -----------------------------------------------------------------------------

type ContentBlogCategory struct {
	ID          uuid.UUID `json:"id"`
	AppID       string    `json:"app_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ContentBlog struct {
	ID              uuid.UUID  `json:"id"`
	AppID           string     `json:"app_id"`
	CategoryID      uuid.UUID  `json:"category_id"`
	CategoryName    string     `json:"category_name"`
	CategorySlug    string     `json:"category_slug"`
	AuthorID        *uuid.UUID `json:"author_id,omitempty"`
	AuthorName      string     `json:"author_name"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Excerpt         string     `json:"excerpt"`
	Body            string     `json:"body"`
	ThumbnailURL    string     `json:"thumbnail_url"`
	Tags            []string   `json:"tags"`
	MetaTitle       string     `json:"meta_title"`
	MetaDescription string     `json:"meta_description"`
	Status          string     `json:"status"` // DRAFT, PUBLISHED, ARCHIVED
	IsFeatured      bool       `json:"is_featured"`
	ViewCount       int64      `json:"view_count"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ListBlogsFilter struct {
	AppID      string
	CategoryID *uuid.UUID
	Status     string
	IsFeatured *bool
	Search     string
	Limit      int32
	Offset     int32
}

// -----------------------------------------------------------------------------
// 3. Content Galleries & Items
// -----------------------------------------------------------------------------

type ContentGallery struct {
	ID          uuid.UUID            `json:"id"`
	AppID       string               `json:"app_id"`
	Title       string               `json:"title"`
	Slug        string               `json:"slug"`
	Description string               `json:"description"`
	Category    string               `json:"category"` // UMRAH, HAJJ, TOUR, EVENT, GENERAL
	CoverURL    string               `json:"cover_url"`
	IsActive    bool                 `json:"is_active"`
	Items       []ContentGalleryItem `json:"items,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type ContentGalleryItem struct {
	ID            uuid.UUID `json:"id"`
	AppID         string    `json:"app_id"`
	GalleryID     uuid.UUID `json:"gallery_id"`
	Title         string    `json:"title"`
	Caption       string    `json:"caption"`
	MediaType     string    `json:"media_type"` // IMAGE, VIDEO
	MediaURL      string    `json:"media_url"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	SequenceOrder int32     `json:"sequence_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListGalleriesFilter struct {
	AppID    string
	Category string
	IsActive *bool
	Limit    int32
	Offset   int32
}

// -----------------------------------------------------------------------------
// 4. Other Contents (Custom Blocks, FAQs, Testimonials)
// -----------------------------------------------------------------------------

type ContentCustomBlock struct {
	ID            uuid.UUID              `json:"id"`
	AppID         string                 `json:"app_id"`
	SectionKey    string                 `json:"section_key"`
	ContentType   string                 `json:"content_type"` // PAGE, SECTION, ANNOUNCEMENT, PARTNER, WIDGET
	Title         string                 `json:"title"`
	Subtitle      string                 `json:"subtitle"`
	Body          string                 `json:"body"`
	ImageURL      string                 `json:"image_url"`
	Metadata      map[string]interface{} `json:"metadata"`
	SequenceOrder int32                  `json:"sequence_order"`
	IsActive      bool                   `json:"is_active"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type ContentFAQ struct {
	ID            uuid.UUID `json:"id"`
	AppID         string    `json:"app_id"`
	Category      string    `json:"category"`
	Question      string    `json:"question"`
	Answer        string    `json:"answer"`
	SequenceOrder int32     `json:"sequence_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ContentTestimonial struct {
	ID            uuid.UUID `json:"id"`
	AppID         string    `json:"app_id"`
	CustomerName  string    `json:"customer_name"`
	CustomerTitle string    `json:"customer_title"`
	AvatarURL     string    `json:"avatar_url"`
	Rating        int32     `json:"rating"`
	Comment       string    `json:"comment"`
	PackageName   string    `json:"package_name"`
	IsFeatured    bool      `json:"is_featured"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListTestimonialsFilter struct {
	AppID      string
	IsActive   *bool
	IsFeatured *bool
	Limit      int32
	Offset     int32
}

// -----------------------------------------------------------------------------
// Content Repository Interface
// -----------------------------------------------------------------------------

type ContentRepository interface {
	WithTx(ctx context.Context, fn func(repo ContentRepository) error) error

	// Banners
	ListBanners(ctx context.Context, filter ListBannersFilter) ([]*ContentBanner, int64, error)
	GetBannerByID(ctx context.Context, appID string, id uuid.UUID) (*ContentBanner, error)
	CreateBanner(ctx context.Context, banner *ContentBanner) (*ContentBanner, error)
	UpdateBanner(ctx context.Context, banner *ContentBanner) (*ContentBanner, error)
	DeleteBanner(ctx context.Context, appID string, id uuid.UUID) error

	// Blog Categories
	ListBlogCategories(ctx context.Context, filter ListBlogCategoriesFilter) ([]*ContentBlogCategory, int64, error)
	GetBlogCategoryByID(ctx context.Context, appID string, id uuid.UUID) (*ContentBlogCategory, error)
	GetBlogCategoryBySlug(ctx context.Context, appID, slug string) (*ContentBlogCategory, error)
	CreateBlogCategory(ctx context.Context, category *ContentBlogCategory) (*ContentBlogCategory, error)
	UpdateBlogCategory(ctx context.Context, category *ContentBlogCategory) (*ContentBlogCategory, error)

	// Blogs
	ListBlogs(ctx context.Context, filter ListBlogsFilter) ([]*ContentBlog, int64, error)
	GetBlogByID(ctx context.Context, appID string, id uuid.UUID) (*ContentBlog, error)
	GetBlogBySlug(ctx context.Context, appID, slug string) (*ContentBlog, error)
	IncrementBlogViewCount(ctx context.Context, id uuid.UUID) error
	CreateBlog(ctx context.Context, blog *ContentBlog) (*ContentBlog, error)
	UpdateBlog(ctx context.Context, blog *ContentBlog) (*ContentBlog, error)
	DeleteBlog(ctx context.Context, appID string, id uuid.UUID) error

	// Galleries & Items
	ListGalleries(ctx context.Context, filter ListGalleriesFilter) ([]*ContentGallery, int64, error)
	GetGalleryByID(ctx context.Context, appID string, id uuid.UUID) (*ContentGallery, error)
	GetGalleryBySlug(ctx context.Context, appID, slug string) (*ContentGallery, error)
	CreateGallery(ctx context.Context, gallery *ContentGallery) (*ContentGallery, error)
	UpdateGallery(ctx context.Context, gallery *ContentGallery) (*ContentGallery, error)
	DeleteGallery(ctx context.Context, appID string, id uuid.UUID) error
	ListGalleryItems(ctx context.Context, galleryID uuid.UUID, isActive *bool) ([]ContentGalleryItem, error)
	AddGalleryItem(ctx context.Context, item *ContentGalleryItem) (*ContentGalleryItem, error)
	DeleteGalleryItem(ctx context.Context, appID string, id uuid.UUID) error

	// Custom Blocks
	ListCustomBlocks(ctx context.Context, filter ListCustomBlocksFilter) ([]*ContentCustomBlock, int64, error)
	GetCustomBlockBySection(ctx context.Context, appID, sectionKey string) (*ContentCustomBlock, error)
	UpsertCustomBlock(ctx context.Context, block *ContentCustomBlock) (*ContentCustomBlock, error)
	DeleteCustomBlock(ctx context.Context, appID string, id uuid.UUID) error

	// FAQs
	ListFAQs(ctx context.Context, filter ListFAQsFilter) ([]*ContentFAQ, int64, error)
	CreateFAQ(ctx context.Context, faq *ContentFAQ) (*ContentFAQ, error)
	UpdateFAQ(ctx context.Context, faq *ContentFAQ) (*ContentFAQ, error)
	DeleteFAQ(ctx context.Context, appID string, id uuid.UUID) error

	// Testimonials
	ListTestimonials(ctx context.Context, filter ListTestimonialsFilter) ([]*ContentTestimonial, int64, error)
	CreateTestimonial(ctx context.Context, testimonial *ContentTestimonial) (*ContentTestimonial, error)
	UpdateTestimonial(ctx context.Context, testimonial *ContentTestimonial) (*ContentTestimonial, error)
	DeleteTestimonial(ctx context.Context, appID string, id uuid.UUID) error
}
