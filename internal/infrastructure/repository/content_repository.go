package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"content-service/internal/domain/entity"
	domainErrors "content-service/internal/domain/errors"
	"content-service/internal/infrastructure/database/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contentRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewContentRepository(pool *pgxpool.Pool) entity.ContentRepository {
	return &contentRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *contentRepository) WithTx(ctx context.Context, fn func(repo entity.ContentRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repository/begin-tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txRepo := &contentRepository{
		pool:    r.pool,
		queries: r.queries.WithTx(tx),
	}

	if err := fn(txRepo); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("repository/commit-tx: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------
// 1. Content Banners
// -----------------------------------------------------------------------------

func (r *contentRepository) ListBanners(ctx context.Context, filter entity.ListBannersFilter) ([]*entity.ContentBanner, int64, error) {
	var isActive pgtype.Bool
	if filter.IsActive != nil {
		isActive = pgtype.Bool{Bool: *filter.IsActive, Valid: true}
	}
	var position pgtype.Text
	if filter.Position != "" {
		position = pgtype.Text{String: filter.Position, Valid: true}
	}

	total, err := r.queries.CountBanners(ctx, sqlc.CountBannersParams{
		AppID:    filter.AppID,
		Position: position,
		IsActive: isActive,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/count-banners: %w", err)
	}

	rows, err := r.queries.ListBanners(ctx, sqlc.ListBannersParams{
		AppID:    filter.AppID,
		Position: position,
		IsActive: isActive,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/list-banners: %w", err)
	}

	banners := make([]*entity.ContentBanner, len(rows))
	for i, b := range rows {
		banners[i] = toDomainBanner(&b)
	}

	return banners, total, nil
}

func (r *contentRepository) GetBannerByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBanner, error) {
	row, err := r.queries.GetBannerByID(ctx, sqlc.GetBannerByIDParams{
		ID:    id,
		AppID: appID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBannerNotFound
		}
		return nil, fmt.Errorf("repository/get-banner: %w", err)
	}
	return toDomainBanner(&row), nil
}

func (r *contentRepository) CreateBanner(ctx context.Context, b *entity.ContentBanner) (*entity.ContentBanner, error) {
	additionalJSON, _ := json.Marshal(b.AdditionalData)
	if additionalJSON == nil {
		additionalJSON = []byte("{}")
	}

	row, err := r.queries.CreateBanner(ctx, sqlc.CreateBannerParams{
		AppID:          b.AppID,
		Title:          b.Title,
		Subtitle:       b.Subtitle,
		Description:    b.Description,
		ImageUrl:       b.ImageURL,
		MobileImageUrl: b.MobileImageURL,
		TargetUrl:      b.TargetURL,
		ButtonText:     b.ButtonText,
		Position:       b.Position,
		SequenceOrder:  b.SequenceOrder,
		StartDate:      ptrTimeToTimestamptz(b.StartDate),
		EndDate:        ptrTimeToTimestamptz(b.EndDate),
		IsActive:       b.IsActive,
		AdditionalData: additionalJSON,
		CreatedBy:      ptrUUIDToPgtypeUUID(b.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("repository/create-banner: %w", err)
	}
	return toDomainBanner(&row), nil
}

func (r *contentRepository) UpdateBanner(ctx context.Context, b *entity.ContentBanner) (*entity.ContentBanner, error) {
	additionalJSON, _ := json.Marshal(b.AdditionalData)
	if additionalJSON == nil {
		additionalJSON = []byte("{}")
	}

	row, err := r.queries.UpdateBanner(ctx, sqlc.UpdateBannerParams{
		ID:             b.ID,
		AppID:          b.AppID,
		Title:          b.Title,
		Subtitle:       b.Subtitle,
		Description:    b.Description,
		ImageUrl:       b.ImageURL,
		MobileImageUrl: b.MobileImageURL,
		TargetUrl:      b.TargetURL,
		ButtonText:     b.ButtonText,
		Position:       b.Position,
		SequenceOrder:  b.SequenceOrder,
		StartDate:      ptrTimeToTimestamptz(b.StartDate),
		EndDate:        ptrTimeToTimestamptz(b.EndDate),
		IsActive:       b.IsActive,
		AdditionalData: additionalJSON,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBannerNotFound
		}
		return nil, fmt.Errorf("repository/update-banner: %w", err)
	}
	return toDomainBanner(&row), nil
}

func (r *contentRepository) DeleteBanner(ctx context.Context, appID string, id uuid.UUID) error {
	return r.queries.DeleteBanner(ctx, sqlc.DeleteBannerParams{
		ID:    id,
		AppID: appID,
	})
}

// -----------------------------------------------------------------------------
// 2. Content Blog & Categories
// -----------------------------------------------------------------------------

func (r *contentRepository) ListBlogCategories(ctx context.Context, appID string) ([]*entity.ContentBlogCategory, error) {
	rows, err := r.queries.ListBlogCategories(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("repository/list-blog-categories: %w", err)
	}

	categories := make([]*entity.ContentBlogCategory, len(rows))
	for i, c := range rows {
		categories[i] = &entity.ContentBlogCategory{
			ID:          c.ID,
			AppID:       c.AppID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			IsActive:    c.IsActive,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
	}
	return categories, nil
}

func (r *contentRepository) GetBlogCategoryByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBlogCategory, error) {
	row, err := r.queries.GetBlogCategoryByID(ctx, sqlc.GetBlogCategoryByIDParams{
		ID:    id,
		AppID: appID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBlogCategoryNotFound
		}
		return nil, fmt.Errorf("repository/get-blog-category: %w", err)
	}
	return &entity.ContentBlogCategory{
		ID:          row.ID,
		AppID:       row.AppID,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: row.Description,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *contentRepository) GetBlogCategoryBySlug(ctx context.Context, appID, slug string) (*entity.ContentBlogCategory, error) {
	row, err := r.queries.GetBlogCategoryBySlug(ctx, sqlc.GetBlogCategoryBySlugParams{
		AppID: appID,
		Slug:  slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBlogCategoryNotFound
		}
		return nil, fmt.Errorf("repository/get-blog-category-by-slug: %w", err)
	}
	return &entity.ContentBlogCategory{
		ID:          row.ID,
		AppID:       row.AppID,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: row.Description,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *contentRepository) CreateBlogCategory(ctx context.Context, c *entity.ContentBlogCategory) (*entity.ContentBlogCategory, error) {
	row, err := r.queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		AppID:       c.AppID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		IsActive:    c.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/create-blog-category: %w", err)
	}
	return &entity.ContentBlogCategory{
		ID:          row.ID,
		AppID:       row.AppID,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: row.Description,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *contentRepository) ListBlogs(ctx context.Context, filter entity.ListBlogsFilter) ([]*entity.ContentBlog, int64, error) {
	var categoryID pgtype.UUID
	if filter.CategoryID != nil {
		categoryID = pgtype.UUID{Bytes: *filter.CategoryID, Valid: true}
	}
	var isFeatured pgtype.Bool
	if filter.IsFeatured != nil {
		isFeatured = pgtype.Bool{Bool: *filter.IsFeatured, Valid: true}
	}
	var statusText pgtype.Text
	if filter.Status != "" {
		statusText = pgtype.Text{String: filter.Status, Valid: true}
	}
	var searchText pgtype.Text
	if filter.Search != "" {
		searchText = pgtype.Text{String: filter.Search, Valid: true}
	}

	total, err := r.queries.CountBlogs(ctx, sqlc.CountBlogsParams{
		AppID:      filter.AppID,
		CategoryID: categoryID,
		Status:     statusText,
		IsFeatured: isFeatured,
		Search:     searchText,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/count-blogs: %w", err)
	}

	rows, err := r.queries.ListBlogs(ctx, sqlc.ListBlogsParams{
		AppID:      filter.AppID,
		CategoryID: categoryID,
		Status:     statusText,
		IsFeatured: isFeatured,
		Search:     searchText,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/list-blogs: %w", err)
	}

	blogs := make([]*entity.ContentBlog, len(rows))
	for i, b := range rows {
		blogs[i] = &entity.ContentBlog{
			ID:              b.ID,
			AppID:           b.AppID,
			CategoryID:      b.CategoryID,
			CategoryName:    b.CategoryName,
			CategorySlug:    b.CategorySlug,
			AuthorID:        pgtypeUUIDToPtrUUID(b.AuthorID),
			AuthorName:      b.AuthorName,
			Title:           b.Title,
			Slug:            b.Slug,
			Excerpt:         b.Excerpt,
			Body:            b.Body,
			ThumbnailURL:    b.ThumbnailUrl,
			Tags:            b.Tags,
			MetaTitle:       b.MetaTitle,
			MetaDescription: b.MetaDescription,
			Status:          b.Status,
			IsFeatured:      b.IsFeatured,
			ViewCount:       b.ViewCount,
			PublishedAt:     timestamptzToPtrTime(b.PublishedAt),
			CreatedAt:       b.CreatedAt,
			UpdatedAt:       b.UpdatedAt,
		}
	}

	return blogs, total, nil
}

func (r *contentRepository) GetBlogByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentBlog, error) {
	b, err := r.queries.GetBlogByID(ctx, sqlc.GetBlogByIDParams{
		ID:    id,
		AppID: appID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBlogNotFound
		}
		return nil, fmt.Errorf("repository/get-blog: %w", err)
	}

	return &entity.ContentBlog{
		ID:              b.ID,
		AppID:           b.AppID,
		CategoryID:      b.CategoryID,
		CategoryName:    b.CategoryName,
		CategorySlug:    b.CategorySlug,
		AuthorID:        pgtypeUUIDToPtrUUID(b.AuthorID),
		AuthorName:      b.AuthorName,
		Title:           b.Title,
		Slug:            b.Slug,
		Excerpt:         b.Excerpt,
		Body:            b.Body,
		ThumbnailURL:    b.ThumbnailUrl,
		Tags:            b.Tags,
		MetaTitle:       b.MetaTitle,
		MetaDescription: b.MetaDescription,
		Status:          b.Status,
		IsFeatured:      b.IsFeatured,
		ViewCount:       b.ViewCount,
		PublishedAt:     timestamptzToPtrTime(b.PublishedAt),
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}, nil
}

func (r *contentRepository) GetBlogBySlug(ctx context.Context, appID, slug string) (*entity.ContentBlog, error) {
	b, err := r.queries.GetBlogBySlug(ctx, sqlc.GetBlogBySlugParams{
		AppID: appID,
		Slug:  slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBlogNotFound
		}
		return nil, fmt.Errorf("repository/get-blog-by-slug: %w", err)
	}

	return &entity.ContentBlog{
		ID:              b.ID,
		AppID:           b.AppID,
		CategoryID:      b.CategoryID,
		CategoryName:    b.CategoryName,
		CategorySlug:    b.CategorySlug,
		AuthorID:        pgtypeUUIDToPtrUUID(b.AuthorID),
		AuthorName:      b.AuthorName,
		Title:           b.Title,
		Slug:            b.Slug,
		Excerpt:         b.Excerpt,
		Body:            b.Body,
		ThumbnailURL:    b.ThumbnailUrl,
		Tags:            b.Tags,
		MetaTitle:       b.MetaTitle,
		MetaDescription: b.MetaDescription,
		Status:          b.Status,
		IsFeatured:      b.IsFeatured,
		ViewCount:       b.ViewCount,
		PublishedAt:     timestamptzToPtrTime(b.PublishedAt),
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}, nil
}

func (r *contentRepository) IncrementBlogViewCount(ctx context.Context, id uuid.UUID) error {
	return r.queries.IncrementBlogViewCount(ctx, id)
}

func (r *contentRepository) CreateBlog(ctx context.Context, b *entity.ContentBlog) (*entity.ContentBlog, error) {
	row, err := r.queries.CreateBlog(ctx, sqlc.CreateBlogParams{
		AppID:           b.AppID,
		CategoryID:      b.CategoryID,
		AuthorID:        ptrUUIDToPgtypeUUID(b.AuthorID),
		AuthorName:      b.AuthorName,
		Title:           b.Title,
		Slug:            b.Slug,
		Excerpt:         b.Excerpt,
		Body:            b.Body,
		ThumbnailUrl:    b.ThumbnailURL,
		Tags:            b.Tags,
		MetaTitle:       b.MetaTitle,
		MetaDescription: b.MetaDescription,
		Status:          b.Status,
		IsFeatured:      b.IsFeatured,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/create-blog: %w", err)
	}

	return &entity.ContentBlog{
		ID:              row.ID,
		AppID:           row.AppID,
		CategoryID:      row.CategoryID,
		AuthorID:        pgtypeUUIDToPtrUUID(row.AuthorID),
		AuthorName:      row.AuthorName,
		Title:           row.Title,
		Slug:            row.Slug,
		Excerpt:         row.Excerpt,
		Body:            row.Body,
		ThumbnailURL:    row.ThumbnailUrl,
		Tags:            row.Tags,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Status:          row.Status,
		IsFeatured:      row.IsFeatured,
		ViewCount:       row.ViewCount,
		PublishedAt:     timestamptzToPtrTime(row.PublishedAt),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *contentRepository) UpdateBlog(ctx context.Context, b *entity.ContentBlog) (*entity.ContentBlog, error) {
	row, err := r.queries.UpdateBlog(ctx, sqlc.UpdateBlogParams{
		ID:              b.ID,
		AppID:           b.AppID,
		CategoryID:      b.CategoryID,
		Title:           b.Title,
		Slug:            b.Slug,
		Excerpt:         b.Excerpt,
		Body:            b.Body,
		ThumbnailUrl:    b.ThumbnailURL,
		Tags:            b.Tags,
		MetaTitle:       b.MetaTitle,
		MetaDescription: b.MetaDescription,
		Status:          b.Status,
		IsFeatured:      b.IsFeatured,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrBlogNotFound
		}
		return nil, fmt.Errorf("repository/update-blog: %w", err)
	}

	return &entity.ContentBlog{
		ID:              row.ID,
		AppID:           row.AppID,
		CategoryID:      row.CategoryID,
		AuthorID:        pgtypeUUIDToPtrUUID(row.AuthorID),
		AuthorName:      row.AuthorName,
		Title:           row.Title,
		Slug:            row.Slug,
		Excerpt:         row.Excerpt,
		Body:            row.Body,
		ThumbnailURL:    row.ThumbnailUrl,
		Tags:            row.Tags,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Status:          row.Status,
		IsFeatured:      row.IsFeatured,
		ViewCount:       row.ViewCount,
		PublishedAt:     timestamptzToPtrTime(row.PublishedAt),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *contentRepository) DeleteBlog(ctx context.Context, appID string, id uuid.UUID) error {
	return r.queries.DeleteBlog(ctx, sqlc.DeleteBlogParams{
		ID:    id,
		AppID: appID,
	})
}

// -----------------------------------------------------------------------------
// 3. Galleries & Items
// -----------------------------------------------------------------------------

func (r *contentRepository) ListGalleries(ctx context.Context, filter entity.ListGalleriesFilter) ([]*entity.ContentGallery, int64, error) {
	var isActive pgtype.Bool
	if filter.IsActive != nil {
		isActive = pgtype.Bool{Bool: *filter.IsActive, Valid: true}
	}
	var cat pgtype.Text
	if filter.Category != "" {
		cat = pgtype.Text{String: filter.Category, Valid: true}
	}

	total, err := r.queries.CountGalleries(ctx, sqlc.CountGalleriesParams{
		AppID:    filter.AppID,
		Category: cat,
		IsActive: isActive,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/count-galleries: %w", err)
	}

	rows, err := r.queries.ListGalleries(ctx, sqlc.ListGalleriesParams{
		AppID:    filter.AppID,
		Category: cat,
		IsActive: isActive,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/list-galleries: %w", err)
	}

	galleries := make([]*entity.ContentGallery, len(rows))
	for i, g := range rows {
		galleries[i] = &entity.ContentGallery{
			ID:          g.ID,
			AppID:       g.AppID,
			Title:       g.Title,
			Slug:        g.Slug,
			Description: g.Description,
			Category:    g.Category,
			CoverURL:    g.CoverUrl,
			IsActive:    g.IsActive,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}
	}

	return galleries, total, nil
}

func (r *contentRepository) GetGalleryByID(ctx context.Context, appID string, id uuid.UUID) (*entity.ContentGallery, error) {
	row, err := r.queries.GetGalleryByID(ctx, sqlc.GetGalleryByIDParams{
		ID:    id,
		AppID: appID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrGalleryNotFound
		}
		return nil, fmt.Errorf("repository/get-gallery: %w", err)
	}

	items, err := r.ListGalleryItems(ctx, row.ID, nil)
	if err != nil {
		items = []entity.ContentGalleryItem{}
	}

	return &entity.ContentGallery{
		ID:          row.ID,
		AppID:       row.AppID,
		Title:       row.Title,
		Slug:        row.Slug,
		Description: row.Description,
		Category:    row.Category,
		CoverURL:    row.CoverUrl,
		IsActive:    row.IsActive,
		Items:       items,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *contentRepository) GetGalleryBySlug(ctx context.Context, appID, slug string) (*entity.ContentGallery, error) {
	row, err := r.queries.GetGalleryBySlug(ctx, sqlc.GetGalleryBySlugParams{
		AppID: appID,
		Slug:  slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrGalleryNotFound
		}
		return nil, fmt.Errorf("repository/get-gallery-by-slug: %w", err)
	}

	items, err := r.ListGalleryItems(ctx, row.ID, nil)
	if err != nil {
		items = []entity.ContentGalleryItem{}
	}

	return &entity.ContentGallery{
		ID:          row.ID,
		AppID:       row.AppID,
		Title:       row.Title,
		Slug:        row.Slug,
		Description: row.Description,
		Category:    row.Category,
		CoverURL:    row.CoverUrl,
		IsActive:    row.IsActive,
		Items:       items,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *contentRepository) CreateGallery(ctx context.Context, g *entity.ContentGallery) (*entity.ContentGallery, error) {
	row, err := r.queries.CreateGallery(ctx, sqlc.CreateGalleryParams{
		AppID:       g.AppID,
		Title:       g.Title,
		Slug:        g.Slug,
		Description: g.Description,
		Category:    g.Category,
		CoverUrl:    g.CoverURL,
		IsActive:    g.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/create-gallery: %w", err)
	}

	return &entity.ContentGallery{
		ID:          row.ID,
		AppID:       row.AppID,
		Title:       row.Title,
		Slug:        row.Slug,
		Description: row.Description,
		Category:    row.Category,
		CoverURL:    row.CoverUrl,
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *contentRepository) DeleteGallery(ctx context.Context, appID string, id uuid.UUID) error {
	return r.queries.DeleteGallery(ctx, sqlc.DeleteGalleryParams{
		ID:    id,
		AppID: appID,
	})
}

func (r *contentRepository) ListGalleryItems(ctx context.Context, galleryID uuid.UUID, isActive *bool) ([]entity.ContentGalleryItem, error) {
	var activeFilter pgtype.Bool
	if isActive != nil {
		activeFilter = pgtype.Bool{Bool: *isActive, Valid: true}
	}

	rows, err := r.queries.ListGalleryItems(ctx, sqlc.ListGalleryItemsParams{
		GalleryID: galleryID,
		IsActive:  activeFilter,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/list-gallery-items: %w", err)
	}

	items := make([]entity.ContentGalleryItem, len(rows))
	for i, it := range rows {
		items[i] = entity.ContentGalleryItem{
			ID:            it.ID,
			AppID:         it.AppID,
			GalleryID:     it.GalleryID,
			Title:         it.Title,
			Caption:       it.Caption,
			MediaType:     it.MediaType,
			MediaURL:      it.MediaUrl,
			ThumbnailURL:  it.ThumbnailUrl,
			SequenceOrder: it.SequenceOrder,
			IsActive:      it.IsActive,
			CreatedAt:     it.CreatedAt,
			UpdatedAt:     it.UpdatedAt,
		}
	}
	return items, nil
}

func (r *contentRepository) AddGalleryItem(ctx context.Context, item *entity.ContentGalleryItem) (*entity.ContentGalleryItem, error) {
	row, err := r.queries.AddGalleryItem(ctx, sqlc.AddGalleryItemParams{
		AppID:         item.AppID,
		GalleryID:     item.GalleryID,
		Title:         item.Title,
		Caption:       item.Caption,
		MediaType:     item.MediaType,
		MediaUrl:      item.MediaURL,
		ThumbnailUrl:  item.ThumbnailURL,
		SequenceOrder: item.SequenceOrder,
		IsActive:      item.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/add-gallery-item: %w", err)
	}

	return &entity.ContentGalleryItem{
		ID:            row.ID,
		AppID:         row.AppID,
		GalleryID:     row.GalleryID,
		Title:         row.Title,
		Caption:       row.Caption,
		MediaType:     row.MediaType,
		MediaURL:      row.MediaUrl,
		ThumbnailURL:  row.ThumbnailUrl,
		SequenceOrder: row.SequenceOrder,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *contentRepository) DeleteGalleryItem(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteGalleryItem(ctx, id)
}

// -----------------------------------------------------------------------------
// 4. Custom Blocks
// -----------------------------------------------------------------------------

func (r *contentRepository) ListCustomBlocks(ctx context.Context, appID, contentType string, isActive *bool) ([]*entity.ContentCustomBlock, error) {
	var activeFilter pgtype.Bool
	if isActive != nil {
		activeFilter = pgtype.Bool{Bool: *isActive, Valid: true}
	}
	var cType pgtype.Text
	if contentType != "" {
		cType = pgtype.Text{String: contentType, Valid: true}
	}

	rows, err := r.queries.ListCustomBlocks(ctx, sqlc.ListCustomBlocksParams{
		AppID:       appID,
		ContentType: cType,
		IsActive:    activeFilter,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/list-custom-blocks: %w", err)
	}

	blocks := make([]*entity.ContentCustomBlock, len(rows))
	for i, b := range rows {
		blocks[i] = toDomainCustomBlock(&b)
	}
	return blocks, nil
}

func (r *contentRepository) GetCustomBlockBySection(ctx context.Context, appID, sectionKey string) (*entity.ContentCustomBlock, error) {
	row, err := r.queries.GetCustomBlockBySection(ctx, sqlc.GetCustomBlockBySectionParams{
		AppID:      appID,
		SectionKey: sectionKey,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrCustomBlockNotFound
		}
		return nil, fmt.Errorf("repository/get-custom-block: %w", err)
	}
	return toDomainCustomBlock(&row), nil
}

func (r *contentRepository) UpsertCustomBlock(ctx context.Context, b *entity.ContentCustomBlock) (*entity.ContentCustomBlock, error) {
	metaJSON, _ := json.Marshal(b.Metadata)
	if metaJSON == nil {
		metaJSON = []byte("{}")
	}

	row, err := r.queries.UpsertCustomBlock(ctx, sqlc.UpsertCustomBlockParams{
		AppID:         b.AppID,
		SectionKey:    b.SectionKey,
		ContentType:   b.ContentType,
		Title:         b.Title,
		Subtitle:      b.Subtitle,
		Body:          b.Body,
		ImageUrl:      b.ImageURL,
		Metadata:      metaJSON,
		SequenceOrder: b.SequenceOrder,
		IsActive:      b.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/upsert-custom-block: %w", err)
	}
	return toDomainCustomBlock(&row), nil
}

func (r *contentRepository) DeleteCustomBlock(ctx context.Context, appID string, id uuid.UUID) error {
	return r.queries.DeleteCustomBlock(ctx, sqlc.DeleteCustomBlockParams{
		ID:    id,
		AppID: appID,
	})
}

// -----------------------------------------------------------------------------
// 5. FAQs
// -----------------------------------------------------------------------------

func (r *contentRepository) ListFAQs(ctx context.Context, appID, category string, isActive *bool) ([]*entity.ContentFAQ, error) {
	var activeFilter pgtype.Bool
	if isActive != nil {
		activeFilter = pgtype.Bool{Bool: *isActive, Valid: true}
	}
	var cat pgtype.Text
	if category != "" {
		cat = pgtype.Text{String: category, Valid: true}
	}

	rows, err := r.queries.ListFAQs(ctx, sqlc.ListFAQsParams{
		AppID:    appID,
		Category: cat,
		IsActive: activeFilter,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/list-faqs: %w", err)
	}

	faqs := make([]*entity.ContentFAQ, len(rows))
	for i, f := range rows {
		faqs[i] = &entity.ContentFAQ{
			ID:            f.ID,
			AppID:         f.AppID,
			Category:      f.Category,
			Question:      f.Question,
			Answer:        f.Answer,
			SequenceOrder: f.SequenceOrder,
			IsActive:      f.IsActive,
			CreatedAt:     f.CreatedAt,
			UpdatedAt:     f.UpdatedAt,
		}
	}
	return faqs, nil
}

func (r *contentRepository) CreateFAQ(ctx context.Context, f *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	row, err := r.queries.CreateFAQ(ctx, sqlc.CreateFAQParams{
		AppID:         f.AppID,
		Category:      f.Category,
		Question:      f.Question,
		Answer:        f.Answer,
		SequenceOrder: f.SequenceOrder,
		IsActive:      f.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/create-faq: %w", err)
	}

	return &entity.ContentFAQ{
		ID:            row.ID,
		AppID:         row.AppID,
		Category:      row.Category,
		Question:      row.Question,
		Answer:        row.Answer,
		SequenceOrder: row.SequenceOrder,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *contentRepository) UpdateFAQ(ctx context.Context, f *entity.ContentFAQ) (*entity.ContentFAQ, error) {
	row, err := r.queries.UpdateFAQ(ctx, sqlc.UpdateFAQParams{
		ID:            f.ID,
		AppID:         f.AppID,
		Category:      f.Category,
		Question:      f.Question,
		Answer:        f.Answer,
		SequenceOrder: f.SequenceOrder,
		IsActive:      f.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrFAQNotFound
		}
		return nil, fmt.Errorf("repository/update-faq: %w", err)
	}

	return &entity.ContentFAQ{
		ID:            row.ID,
		AppID:         row.AppID,
		Category:      row.Category,
		Question:      row.Question,
		Answer:        row.Answer,
		SequenceOrder: row.SequenceOrder,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *contentRepository) DeleteFAQ(ctx context.Context, appID string, id uuid.UUID) error {
	return r.queries.DeleteFAQ(ctx, sqlc.DeleteFAQParams{
		ID:    id,
		AppID: appID,
	})
}

// -----------------------------------------------------------------------------
// 6. Testimonials
// -----------------------------------------------------------------------------

func (r *contentRepository) ListTestimonials(ctx context.Context, filter entity.ListTestimonialsFilter) ([]*entity.ContentTestimonial, int64, error) {
	var activeFilter pgtype.Bool
	if filter.IsActive != nil {
		activeFilter = pgtype.Bool{Bool: *filter.IsActive, Valid: true}
	}
	var featuredFilter pgtype.Bool
	if filter.IsFeatured != nil {
		featuredFilter = pgtype.Bool{Bool: *filter.IsFeatured, Valid: true}
	}

	total, err := r.queries.CountTestimonials(ctx, sqlc.CountTestimonialsParams{
		AppID:      filter.AppID,
		IsActive:   activeFilter,
		IsFeatured: featuredFilter,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/count-testimonials: %w", err)
	}

	rows, err := r.queries.ListTestimonials(ctx, sqlc.ListTestimonialsParams{
		AppID:      filter.AppID,
		IsActive:   activeFilter,
		IsFeatured: featuredFilter,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository/list-testimonials: %w", err)
	}

	testimonials := make([]*entity.ContentTestimonial, len(rows))
	for i, t := range rows {
		testimonials[i] = &entity.ContentTestimonial{
			ID:            t.ID,
			AppID:         t.AppID,
			CustomerName:  t.CustomerName,
			CustomerTitle: t.CustomerTitle,
			AvatarURL:     t.AvatarUrl,
			Rating:        t.Rating,
			Comment:       t.Comment,
			PackageName:   t.PackageName,
			IsFeatured:    t.IsFeatured,
			IsActive:      t.IsActive,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
	}

	return testimonials, total, nil
}

func (r *contentRepository) CreateTestimonial(ctx context.Context, t *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	row, err := r.queries.CreateTestimonial(ctx, sqlc.CreateTestimonialParams{
		AppID:         t.AppID,
		CustomerName:  t.CustomerName,
		CustomerTitle: t.CustomerTitle,
		AvatarUrl:     t.AvatarURL,
		Rating:        t.Rating,
		Comment:       t.Comment,
		PackageName:   t.PackageName,
		IsFeatured:    t.IsFeatured,
		IsActive:      t.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("repository/create-testimonial: %w", err)
	}

	return &entity.ContentTestimonial{
		ID:            row.ID,
		AppID:         row.AppID,
		CustomerName:  row.CustomerName,
		CustomerTitle: row.CustomerTitle,
		AvatarURL:     row.AvatarUrl,
		Rating:        row.Rating,
		Comment:       row.Comment,
		PackageName:   row.PackageName,
		IsFeatured:    row.IsFeatured,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *contentRepository) UpdateTestimonial(ctx context.Context, t *entity.ContentTestimonial) (*entity.ContentTestimonial, error) {
	row, err := r.queries.UpdateTestimonial(ctx, sqlc.UpdateTestimonialParams{
		ID:            t.ID,
		AppID:         t.AppID,
		CustomerName:  t.CustomerName,
		CustomerTitle: t.CustomerTitle,
		AvatarUrl:     t.AvatarURL,
		Rating:        t.Rating,
		Comment:       t.Comment,
		PackageName:   t.PackageName,
		IsFeatured:    t.IsFeatured,
		IsActive:      t.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrTestimonialNotFound
		}
		return nil, fmt.Errorf("repository/update-testimonial: %w", err)
	}

	return &entity.ContentTestimonial{
		ID:            row.ID,
		AppID:         row.AppID,
		CustomerName:  row.CustomerName,
		CustomerTitle: row.CustomerTitle,
		AvatarURL:     row.AvatarUrl,
		Rating:        row.Rating,
		Comment:       row.Comment,
		PackageName:   row.PackageName,
		IsFeatured:    row.IsFeatured,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *contentRepository) DeleteTestimonial(ctx context.Context, appID string, id uuid.UUID) error {
	return r.queries.DeleteTestimonial(ctx, sqlc.DeleteTestimonialParams{
		ID:    id,
		AppID: appID,
	})
}

// -----------------------------------------------------------------------------
// Mapping Helpers
// -----------------------------------------------------------------------------

func toDomainBanner(b *sqlc.ContentBanner) *entity.ContentBanner {
	var additional map[string]interface{}
	if len(b.AdditionalData) > 0 {
		_ = json.Unmarshal(b.AdditionalData, &additional)
	}

	return &entity.ContentBanner{
		ID:             b.ID,
		AppID:          b.AppID,
		Title:          b.Title,
		Subtitle:       b.Subtitle,
		Description:    b.Description,
		ImageURL:       b.ImageUrl,
		MobileImageURL: b.MobileImageUrl,
		TargetURL:      b.TargetUrl,
		ButtonText:     b.ButtonText,
		Position:       b.Position,
		SequenceOrder:  b.SequenceOrder,
		StartDate:      timestamptzToPtrTime(b.StartDate),
		EndDate:        timestamptzToPtrTime(b.EndDate),
		IsActive:       b.IsActive,
		AdditionalData: additional,
		CreatedBy:      pgtypeUUIDToPtrUUID(b.CreatedBy),
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

func toDomainCustomBlock(b *sqlc.ContentCustomBlock) *entity.ContentCustomBlock {
	var meta map[string]interface{}
	if len(b.Metadata) > 0 {
		_ = json.Unmarshal(b.Metadata, &meta)
	}

	return &entity.ContentCustomBlock{
		ID:            b.ID,
		AppID:         b.AppID,
		SectionKey:    b.SectionKey,
		ContentType:   b.ContentType,
		Title:         b.Title,
		Subtitle:      b.Subtitle,
		Body:          b.Body,
		ImageURL:      b.ImageUrl,
		Metadata:      meta,
		SequenceOrder: b.SequenceOrder,
		IsActive:      b.IsActive,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

func ptrTimeToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func timestamptzToPtrTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	res := t.Time
	return &res
}

func ptrUUIDToPgtypeUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil || *u == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func pgtypeUUIDToPtrUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	res := uuid.UUID(u.Bytes)
	return &res
}
