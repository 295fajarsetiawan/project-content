package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"content-service/internal/delivery/grpc/interceptor"
	"content-service/internal/delivery/grpc/pb"
	"content-service/internal/domain/entity"
	domainErrors "content-service/internal/domain/errors"
	"content-service/internal/usecase/content"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContentHandler struct {
	pb.UnimplementedContentServiceServer
	usecase content.ContentUsecase
}

func NewContentHandler(usecase content.ContentUsecase) *ContentHandler {
	return &ContentHandler{
		usecase: usecase,
	}
}

// =============================================================================
// 1. Banner Handlers
// =============================================================================

func (h *ContentHandler) ListBanners(ctx context.Context, req *pb.ListBannersRequest) (*pb.ListBannersResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	bannerFilter := entity.ListBannersFilter{
		AppID:    appID,
		Position: req.GetPosition(),
		Limit:    perPage,
		Offset:   offset,
	}
	if !interceptor.IsAdmin(ctx) {
		active := true
		bannerFilter.IsActive = &active
	}

	banners, meta, err := h.usecase.ListBanners(ctx, bannerFilter)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoBanners := make([]*pb.BannerData, len(banners))
	for i, b := range banners {
		protoBanners[i] = toProtoBanner(b)
	}

	return &pb.ListBannersResponse{
		Message: "Banners retrieved successfully",
		Data:    protoBanners,
		Pagination: &pb.PaginationMeta{
			CurrentPage: meta.CurrentPage,
			PerPage:     meta.PerPage,
			TotalItems:  meta.TotalItems,
			TotalPages:  meta.TotalPages,
		},
	}, nil
}

func (h *ContentHandler) GetBanner(ctx context.Context, req *pb.GetBannerRequest) (*pb.BannerResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	bannerID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid banner id format")
	}

	b, err := h.usecase.GetBanner(ctx, appID, bannerID)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	if !interceptor.IsAdmin(ctx) && !b.IsActive {
		return nil, mapDomainErrorToGRPC(domainErrors.ErrBannerNotFound)
	}

	return &pb.BannerResponse{
		Message: "Banner retrieved successfully",
		Data:    toProtoBanner(b),
	}, nil
}

func (h *ContentHandler) CreateBanner(ctx context.Context, req *pb.CreateBannerRequest) (*pb.BannerResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	var startDate, endDate *time.Time
	if req.GetStartDate() != "" {
		if t, err := time.Parse(time.RFC3339, req.GetStartDate()); err == nil {
			startDate = &t
		}
	}
	if req.GetEndDate() != "" {
		if t, err := time.Parse(time.RFC3339, req.GetEndDate()); err == nil {
			endDate = &t
		}
	}

	additionalData := make(map[string]interface{})
	if req.GetAdditionalData() != "" {
		_ = json.Unmarshal([]byte(req.GetAdditionalData()), &additionalData)
	}

	var createdBy *uuid.UUID
	if u, ok := interceptor.GetUserContext(ctx); ok && u != nil {
		createdBy = &u.UserID
	}

	banner := &entity.ContentBanner{
		AppID:          appID,
		Title:          req.GetTitle(),
		Subtitle:       req.GetSubtitle(),
		Description:    req.GetDescription(),
		ImageURL:       req.GetImageUrl(),
		MobileImageURL: req.GetMobileImageUrl(),
		TargetURL:      req.GetTargetUrl(),
		ButtonText:     req.GetButtonText(),
		Position:       req.GetPosition(),
		SequenceOrder:  req.GetSequenceOrder(),
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       req.GetIsActive(),
		AdditionalData: additionalData,
		CreatedBy:      createdBy,
	}

	created, err := h.usecase.CreateBanner(ctx, banner)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.BannerResponse{
		Message: "Banner created successfully",
		Data:    toProtoBanner(created),
	}, nil
}

func (h *ContentHandler) UpdateBanner(ctx context.Context, req *pb.UpdateBannerRequest) (*pb.BannerResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	bannerID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid banner id format")
	}

	var startDate, endDate *time.Time
	if req.GetStartDate() != "" {
		if t, err := time.Parse(time.RFC3339, req.GetStartDate()); err == nil {
			startDate = &t
		}
	}
	if req.GetEndDate() != "" {
		if t, err := time.Parse(time.RFC3339, req.GetEndDate()); err == nil {
			endDate = &t
		}
	}

	additionalData := make(map[string]interface{})
	if req.GetAdditionalData() != "" {
		_ = json.Unmarshal([]byte(req.GetAdditionalData()), &additionalData)
	}

	banner := &entity.ContentBanner{
		ID:             bannerID,
		AppID:          appID,
		Title:          req.GetTitle(),
		Subtitle:       req.GetSubtitle(),
		Description:    req.GetDescription(),
		ImageURL:       req.GetImageUrl(),
		MobileImageURL: req.GetMobileImageUrl(),
		TargetURL:      req.GetTargetUrl(),
		ButtonText:     req.GetButtonText(),
		Position:       req.GetPosition(),
		SequenceOrder:  req.GetSequenceOrder(),
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       req.GetIsActive(),
		AdditionalData: additionalData,
	}

	updated, err := h.usecase.UpdateBanner(ctx, banner)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.BannerResponse{
		Message: "Banner updated successfully",
		Data:    toProtoBanner(updated),
	}, nil
}

func (h *ContentHandler) DeleteBanner(ctx context.Context, req *pb.DeleteBannerRequest) (*pb.MessageResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	bannerID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid banner id format")
	}

	if err := h.usecase.DeleteBanner(ctx, appID, bannerID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "Banner deleted successfully",
	}, nil
}

// =============================================================================
// 2. Blog Category & Article Handlers
// =============================================================================

func (h *ContentHandler) ListBlogCategories(ctx context.Context, req *pb.ListBlogCategoriesRequest) (*pb.ListBlogCategoriesResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	filter := entity.ListBlogCategoriesFilter{
		AppID:  appID,
		Search: req.GetSearch(),
		Limit:  perPage,
		Offset: offset,
	}
	if !interceptor.IsAdmin(ctx) {
		active := true
		filter.IsActive = &active
	}
	categories, meta, err := h.usecase.ListBlogCategories(ctx, filter)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoCategories := make([]*pb.BlogCategoryData, len(categories))
	for i, c := range categories {
		protoCategories[i] = toProtoBlogCategory(c)
	}

	return &pb.ListBlogCategoriesResponse{
		Message:    "Blog categories retrieved successfully",
		Data:       protoCategories,
		Pagination: toProtoPagination(meta),
	}, nil
}

func (h *ContentHandler) UpdateBlogCategory(ctx context.Context, req *pb.UpdateBlogCategoryRequest) (*pb.BlogCategoryResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid blog category id format")
	}

	updated, err := h.usecase.UpdateBlogCategory(ctx, &entity.ContentBlogCategory{
		ID:          id,
		AppID:       appID,
		Name:        req.GetName(),
		Slug:        req.GetSlug(),
		Description: req.GetDescription(),
		IsActive:    req.GetIsActive(),
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.BlogCategoryResponse{
		Message: "Blog category updated successfully",
		Data:    toProtoBlogCategory(updated),
	}, nil
}

func (h *ContentHandler) CreateBlogCategory(ctx context.Context, req *pb.CreateBlogCategoryRequest) (*pb.BlogCategoryResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	cat := &entity.ContentBlogCategory{
		AppID:       appID,
		Name:        req.GetName(),
		Slug:        req.GetSlug(),
		Description: req.GetDescription(),
		IsActive:    true,
	}

	created, err := h.usecase.CreateBlogCategory(ctx, cat)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.BlogCategoryResponse{
		Message: "Blog category created successfully",
		Data:    toProtoBlogCategory(created),
	}, nil
}

func (h *ContentHandler) ListBlogs(ctx context.Context, req *pb.ListBlogsRequest) (*pb.ListBlogsResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	var categoryID *uuid.UUID
	if req.GetCategoryId() != "" {
		if parsed, err := uuid.Parse(req.GetCategoryId()); err == nil {
			categoryID = &parsed
		}
	}

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	statusFilter := req.GetStatus()
	if !interceptor.IsAdmin(ctx) {
		statusFilter = entity.BlogStatusPublished
	}

	blogs, meta, err := h.usecase.ListBlogs(ctx, entity.ListBlogsFilter{
		AppID:      appID,
		CategoryID: categoryID,
		Status:     statusFilter,
		Search:     req.GetSearch(),
		Limit:      perPage,
		Offset:     offset,
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoBlogs := make([]*pb.BlogData, len(blogs))
	for i, b := range blogs {
		protoBlogs[i] = toProtoBlog(b)
	}

	return &pb.ListBlogsResponse{
		Message: "Blog articles retrieved successfully",
		Data:    protoBlogs,
		Pagination: &pb.PaginationMeta{
			CurrentPage: meta.CurrentPage,
			PerPage:     meta.PerPage,
			TotalItems:  meta.TotalItems,
			TotalPages:  meta.TotalPages,
		},
	}, nil
}

func (h *ContentHandler) GetBlog(ctx context.Context, req *pb.GetBlogRequest) (*pb.BlogResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	var id *uuid.UUID
	if req.GetId() != "" {
		if parsed, err := uuid.Parse(req.GetId()); err == nil {
			id = &parsed
		}
	}

	b, err := h.usecase.GetBlog(ctx, appID, id, req.GetSlug())
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	if !interceptor.IsAdmin(ctx) && b.Status != entity.BlogStatusPublished {
		return nil, mapDomainErrorToGRPC(domainErrors.ErrBlogNotFound)
	}

	return &pb.BlogResponse{
		Message: "Blog article retrieved successfully",
		Data:    toProtoBlog(b),
	}, nil
}

func (h *ContentHandler) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.BlogResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	catID, err := uuid.Parse(req.GetCategoryId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid category id format")
	}

	var authorID *uuid.UUID
	authorName := req.GetAuthorName()
	if u, ok := interceptor.GetUserContext(ctx); ok && u != nil {
		authorID = &u.UserID
		if authorName == "" {
			authorName = u.Email
		}
	}

	blog := &entity.ContentBlog{
		AppID:           appID,
		CategoryID:      catID,
		AuthorID:        authorID,
		AuthorName:      authorName,
		Title:           req.GetTitle(),
		Slug:            req.GetSlug(),
		Excerpt:         req.GetExcerpt(),
		Body:            req.GetBody(),
		ThumbnailURL:    req.GetThumbnailUrl(),
		Tags:            req.GetTags(),
		MetaTitle:       req.GetMetaTitle(),
		MetaDescription: req.GetMetaDescription(),
		Status:          req.GetStatus(),
		IsFeatured:      req.GetIsFeatured(),
	}

	created, err := h.usecase.CreateBlog(ctx, blog)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.BlogResponse{
		Message: "Blog article created successfully",
		Data:    toProtoBlog(created),
	}, nil
}

func (h *ContentHandler) UpdateBlog(ctx context.Context, req *pb.UpdateBlogRequest) (*pb.BlogResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	blogID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid blog id format")
	}

	var catID uuid.UUID
	if req.GetCategoryId() != "" {
		catID, err = uuid.Parse(req.GetCategoryId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid category id format")
		}
	}

	blog := &entity.ContentBlog{
		ID:              blogID,
		AppID:           appID,
		CategoryID:      catID,
		Title:           req.GetTitle(),
		Slug:            req.GetSlug(),
		Excerpt:         req.GetExcerpt(),
		Body:            req.GetBody(),
		ThumbnailURL:    req.GetThumbnailUrl(),
		Tags:            req.GetTags(),
		MetaTitle:       req.GetMetaTitle(),
		MetaDescription: req.GetMetaDescription(),
		Status:          req.GetStatus(),
		IsFeatured:      req.GetIsFeatured(),
	}

	updated, err := h.usecase.UpdateBlog(ctx, blog)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.BlogResponse{
		Message: "Blog article updated successfully",
		Data:    toProtoBlog(updated),
	}, nil
}

func (h *ContentHandler) DeleteBlog(ctx context.Context, req *pb.DeleteBlogRequest) (*pb.MessageResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	blogID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid blog id format")
	}

	if err := h.usecase.DeleteBlog(ctx, appID, blogID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "Blog article deleted successfully",
	}, nil
}

// =============================================================================
// 3. Gallery Handlers
// =============================================================================

func (h *ContentHandler) ListGalleries(ctx context.Context, req *pb.ListGalleriesRequest) (*pb.ListGalleriesResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	filter := entity.ListGalleriesFilter{
		AppID:    appID,
		Category: req.GetCategory(),
		Limit:    perPage,
		Offset:   offset,
	}
	if !interceptor.IsAdmin(ctx) {
		active := true
		filter.IsActive = &active
	}

	galleries, meta, err := h.usecase.ListGalleries(ctx, filter)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoGalleries := make([]*pb.GalleryData, len(galleries))
	for i, g := range galleries {
		protoGalleries[i] = toProtoGallery(g)
	}

	return &pb.ListGalleriesResponse{
		Message: "Galleries retrieved successfully",
		Data:    protoGalleries,
		Pagination: &pb.PaginationMeta{
			CurrentPage: meta.CurrentPage,
			PerPage:     meta.PerPage,
			TotalItems:  meta.TotalItems,
			TotalPages:  meta.TotalPages,
		},
	}, nil
}

func (h *ContentHandler) UpdateGallery(ctx context.Context, req *pb.UpdateGalleryRequest) (*pb.GalleryResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid gallery id format")
	}

	updated, err := h.usecase.UpdateGallery(ctx, &entity.ContentGallery{
		ID:          id,
		AppID:       appID,
		Title:       req.GetTitle(),
		Slug:        req.GetSlug(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
		CoverURL:    req.GetCoverUrl(),
		IsActive:    req.GetIsActive(),
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.GalleryResponse{
		Message: "Gallery updated successfully",
		Data:    toProtoGallery(updated),
	}, nil
}

func (h *ContentHandler) GetGallery(ctx context.Context, req *pb.GetGalleryRequest) (*pb.GalleryResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	var id *uuid.UUID
	if req.GetId() != "" {
		if parsed, err := uuid.Parse(req.GetId()); err == nil {
			id = &parsed
		}
	}

	g, err := h.usecase.GetGallery(ctx, appID, id, req.GetSlug())
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	if !interceptor.IsAdmin(ctx) && !g.IsActive {
		return nil, mapDomainErrorToGRPC(domainErrors.ErrGalleryNotFound)
	}

	return &pb.GalleryResponse{
		Message: "Gallery retrieved successfully",
		Data:    toProtoGallery(g),
	}, nil
}

func (h *ContentHandler) CreateGallery(ctx context.Context, req *pb.CreateGalleryRequest) (*pb.GalleryResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	gallery := &entity.ContentGallery{
		AppID:       appID,
		Title:       req.GetTitle(),
		Slug:        req.GetSlug(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
		CoverURL:    req.GetCoverUrl(),
		IsActive:    true,
	}

	created, err := h.usecase.CreateGallery(ctx, gallery)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.GalleryResponse{
		Message: "Gallery created successfully",
		Data:    toProtoGallery(created),
	}, nil
}

func (h *ContentHandler) DeleteGallery(ctx context.Context, req *pb.DeleteGalleryRequest) (*pb.MessageResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	galleryID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid gallery id format")
	}

	if err := h.usecase.DeleteGallery(ctx, appID, galleryID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "Gallery deleted successfully",
	}, nil
}

func (h *ContentHandler) AddGalleryItem(ctx context.Context, req *pb.AddGalleryItemRequest) (*pb.GalleryItemResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	galleryID, err := uuid.Parse(req.GetGalleryId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid gallery id format")
	}

	item := &entity.ContentGalleryItem{
		AppID:         appID,
		GalleryID:     galleryID,
		Title:         req.GetTitle(),
		Caption:       req.GetCaption(),
		MediaType:     req.GetMediaType(),
		MediaURL:      req.GetMediaUrl(),
		ThumbnailURL:  req.GetThumbnailUrl(),
		SequenceOrder: req.GetSequenceOrder(),
		IsActive:      true,
	}

	created, err := h.usecase.AddGalleryItem(ctx, item)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.GalleryItemResponse{
		Message: "Gallery item added successfully",
		Data:    toProtoGalleryItem(created),
	}, nil
}

func (h *ContentHandler) DeleteGalleryItem(ctx context.Context, req *pb.DeleteGalleryItemRequest) (*pb.MessageResponse, error) {
	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid gallery item id format")
	}

	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	if err := h.usecase.DeleteGalleryItem(ctx, appID, itemID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "Gallery item deleted successfully",
	}, nil
}

// =============================================================================
// 4. Other Content Handlers (Custom Blocks, FAQs, Testimonials)
// =============================================================================

func (h *ContentHandler) ListCustomBlocks(ctx context.Context, req *pb.ListCustomBlocksRequest) (*pb.ListCustomBlocksResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	filter := entity.ListCustomBlocksFilter{
		AppID:       appID,
		ContentType: req.GetContentType(),
		Search:      req.GetSearch(),
		Limit:       perPage,
		Offset:      offset,
	}
	if !interceptor.IsAdmin(ctx) {
		active := true
		filter.IsActive = &active
	}

	blocks, meta, err := h.usecase.ListCustomBlocks(ctx, filter)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoBlocks := make([]*pb.CustomBlockData, len(blocks))
	for i, b := range blocks {
		protoBlocks[i] = toProtoCustomBlock(b)
	}

	return &pb.ListCustomBlocksResponse{
		Message:    "Custom blocks retrieved successfully",
		Data:       protoBlocks,
		Pagination: toProtoPagination(meta),
	}, nil
}

func (h *ContentHandler) GetCustomBlock(ctx context.Context, req *pb.GetCustomBlockRequest) (*pb.CustomBlockResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	b, err := h.usecase.GetCustomBlock(ctx, appID, req.GetSectionKey())
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	if !interceptor.IsAdmin(ctx) && !b.IsActive {
		return nil, mapDomainErrorToGRPC(domainErrors.ErrCustomBlockNotFound)
	}

	return &pb.CustomBlockResponse{
		Message: "Custom block retrieved successfully",
		Data:    toProtoCustomBlock(b),
	}, nil
}

func (h *ContentHandler) UpsertCustomBlock(ctx context.Context, req *pb.UpsertCustomBlockRequest) (*pb.CustomBlockResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	metadata := make(map[string]interface{})
	if req.GetMetadata() != "" {
		_ = json.Unmarshal([]byte(req.GetMetadata()), &metadata)
	}

	block := &entity.ContentCustomBlock{
		AppID:         appID,
		SectionKey:    req.GetSectionKey(),
		ContentType:   req.GetContentType(),
		Title:         req.GetTitle(),
		Subtitle:      req.GetSubtitle(),
		Body:          req.GetBody(),
		ImageURL:      req.GetImageUrl(),
		Metadata:      metadata,
		SequenceOrder: req.GetSequenceOrder(),
		IsActive:      req.GetIsActive(),
	}

	saved, err := h.usecase.UpsertCustomBlock(ctx, block)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.CustomBlockResponse{
		Message: "Custom block saved successfully",
		Data:    toProtoCustomBlock(saved),
	}, nil
}

func (h *ContentHandler) DeleteCustomBlock(ctx context.Context, req *pb.DeleteCustomBlockRequest) (*pb.MessageResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	blockID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid custom block id format")
	}

	if err := h.usecase.DeleteCustomBlock(ctx, appID, blockID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "Custom block deleted successfully",
	}, nil
}

func (h *ContentHandler) ListFAQs(ctx context.Context, req *pb.ListFAQsRequest) (*pb.ListFAQsResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	filter := entity.ListFAQsFilter{
		AppID:    appID,
		Category: req.GetCategory(),
		Search:   req.GetSearch(),
		Limit:    perPage,
		Offset:   offset,
	}
	if !interceptor.IsAdmin(ctx) {
		active := true
		filter.IsActive = &active
	}

	faqs, meta, err := h.usecase.ListFAQs(ctx, filter)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoFAQs := make([]*pb.FAQData, len(faqs))
	for i, f := range faqs {
		protoFAQs[i] = toProtoFAQ(f)
	}

	return &pb.ListFAQsResponse{
		Message:    "FAQs retrieved successfully",
		Data:       protoFAQs,
		Pagination: toProtoPagination(meta),
	}, nil
}

func (h *ContentHandler) CreateFAQ(ctx context.Context, req *pb.CreateFAQRequest) (*pb.FAQResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	faq := &entity.ContentFAQ{
		AppID:         appID,
		Category:      req.GetCategory(),
		Question:      req.GetQuestion(),
		Answer:        req.GetAnswer(),
		SequenceOrder: req.GetSequenceOrder(),
		IsActive:      true,
	}

	created, err := h.usecase.CreateFAQ(ctx, faq)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.FAQResponse{
		Message: "FAQ created successfully",
		Data:    toProtoFAQ(created),
	}, nil
}

func (h *ContentHandler) UpdateFAQ(ctx context.Context, req *pb.UpdateFAQRequest) (*pb.FAQResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	faqID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid faq id format")
	}

	faq := &entity.ContentFAQ{
		ID:            faqID,
		AppID:         appID,
		Category:      req.GetCategory(),
		Question:      req.GetQuestion(),
		Answer:        req.GetAnswer(),
		SequenceOrder: req.GetSequenceOrder(),
		IsActive:      req.GetIsActive(),
	}

	updated, err := h.usecase.UpdateFAQ(ctx, faq)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.FAQResponse{
		Message: "FAQ updated successfully",
		Data:    toProtoFAQ(updated),
	}, nil
}

func (h *ContentHandler) DeleteFAQ(ctx context.Context, req *pb.DeleteFAQRequest) (*pb.MessageResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	faqID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid faq id format")
	}

	if err := h.usecase.DeleteFAQ(ctx, appID, faqID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "FAQ deleted successfully",
	}, nil
}

func (h *ContentHandler) ListTestimonials(ctx context.Context, req *pb.ListTestimonialsRequest) (*pb.ListTestimonialsResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	perPage, offset := entity.NormalizePageRequest(req.GetPage(), req.GetPerPage())

	testimonialFilter := entity.ListTestimonialsFilter{
		AppID:  appID,
		Limit:  perPage,
		Offset: offset,
	}
	if !interceptor.IsAdmin(ctx) {
		active := true
		testimonialFilter.IsActive = &active
	}

	testimonials, meta, err := h.usecase.ListTestimonials(ctx, testimonialFilter)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	protoTestimonials := make([]*pb.TestimonialData, len(testimonials))
	for i, t := range testimonials {
		protoTestimonials[i] = toProtoTestimonial(t)
	}

	return &pb.ListTestimonialsResponse{
		Message: "Testimonials retrieved successfully",
		Data:    protoTestimonials,
		Pagination: &pb.PaginationMeta{
			CurrentPage: meta.CurrentPage,
			PerPage:     meta.PerPage,
			TotalItems:  meta.TotalItems,
			TotalPages:  meta.TotalPages,
		},
	}, nil
}

func (h *ContentHandler) CreateTestimonial(ctx context.Context, req *pb.CreateTestimonialRequest) (*pb.TestimonialResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	testimonial := &entity.ContentTestimonial{
		AppID:         appID,
		CustomerName:  req.GetCustomerName(),
		CustomerTitle: req.GetCustomerTitle(),
		AvatarURL:     req.GetAvatarUrl(),
		Rating:        req.GetRating(),
		Comment:       req.GetComment(),
		PackageName:   req.GetPackageName(),
		IsFeatured:    req.GetIsFeatured(),
		IsActive:      true,
	}

	created, err := h.usecase.CreateTestimonial(ctx, testimonial)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.TestimonialResponse{
		Message: "Testimonial created successfully",
		Data:    toProtoTestimonial(created),
	}, nil
}

func (h *ContentHandler) UpdateTestimonial(ctx context.Context, req *pb.UpdateTestimonialRequest) (*pb.TestimonialResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	testimonialID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid testimonial id format")
	}

	testimonial := &entity.ContentTestimonial{
		ID:            testimonialID,
		AppID:         appID,
		CustomerName:  req.GetCustomerName(),
		CustomerTitle: req.GetCustomerTitle(),
		AvatarURL:     req.GetAvatarUrl(),
		Rating:        req.GetRating(),
		Comment:       req.GetComment(),
		PackageName:   req.GetPackageName(),
		IsFeatured:    req.GetIsFeatured(),
		IsActive:      req.GetIsActive(),
	}

	updated, err := h.usecase.UpdateTestimonial(ctx, testimonial)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.TestimonialResponse{
		Message: "Testimonial updated successfully",
		Data:    toProtoTestimonial(updated),
	}, nil
}

func (h *ContentHandler) DeleteTestimonial(ctx context.Context, req *pb.DeleteTestimonialRequest) (*pb.MessageResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	testimonialID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid testimonial id format")
	}

	if err := h.usecase.DeleteTestimonial(ctx, appID, testimonialID); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.MessageResponse{
		Message: "Testimonial deleted successfully",
	}, nil
}

// =============================================================================
// 5. Media Upload Handler
// =============================================================================

// uploadMediaRequest declares input rules for UploadMedia via struct tags.
type uploadMediaRequest struct {
	Filename      string `validate:"required,max=255"`
	Base64Content string `validate:"required"`
	Folder        string `validate:"omitempty,max=100"`
}

func (h *ContentHandler) UploadMedia(ctx context.Context, req *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
	appID, _ := interceptor.ExtractAppID(ctx, req.GetAppId())

	if err := validateRequest(uploadMediaRequest{
		Filename:      req.GetFilename(),
		Base64Content: req.GetBase64Content(),
		Folder:        req.GetFolder(),
	}); err != nil {
		return nil, err
	}

	folder := req.GetFolder()
	if folder == "" {
		folder = "general"
	}

	result, err := h.usecase.UploadMedia(ctx, appID, folder, req.GetFilename(), req.GetBase64Content())
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.UploadMediaResponse{
		Message: "Media uploaded successfully",
		Data: &pb.MediaUploadData{
			FileUrl:  result.FileURL,
			FilePath: result.FilePath,
			Filename: result.Filename,
			FileSize: result.FileSize,
			MimeType: result.MimeType,
		},
	}, nil
}

// =============================================================================
// Helpers & Mapping Functions
// =============================================================================

func toProtoBanner(b *entity.ContentBanner) *pb.BannerData {
	if b == nil {
		return nil
	}

	var startDateStr, endDateStr string
	if b.StartDate != nil {
		startDateStr = b.StartDate.Format(time.RFC3339)
	}
	if b.EndDate != nil {
		endDateStr = b.EndDate.Format(time.RFC3339)
	}

	additionalDataStr := "{}"
	if b.AdditionalData != nil {
		if bytes, err := json.Marshal(b.AdditionalData); err == nil {
			additionalDataStr = string(bytes)
		}
	}

	return &pb.BannerData{
		Id:             b.ID.String(),
		AppId:          b.AppID,
		Title:          b.Title,
		Subtitle:       b.Subtitle,
		Description:    b.Description,
		ImageUrl:       b.ImageURL,
		MobileImageUrl: b.MobileImageURL,
		TargetUrl:      b.TargetURL,
		ButtonText:     b.ButtonText,
		Position:       b.Position,
		SequenceOrder:  b.SequenceOrder,
		StartDate:      startDateStr,
		EndDate:        endDateStr,
		IsActive:       b.IsActive,
		AdditionalData: additionalDataStr,
		CreatedAt:      b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      b.UpdatedAt.Format(time.RFC3339),
	}
}

// toProtoPagination converts usecase pagination metadata to its protobuf form.
func toProtoPagination(m *entity.PaginationMeta) *pb.PaginationMeta {
	if m == nil {
		return nil
	}
	return &pb.PaginationMeta{
		CurrentPage: m.CurrentPage,
		PerPage:     m.PerPage,
		TotalItems:  m.TotalItems,
		TotalPages:  m.TotalPages,
	}
}

func toProtoBlogCategory(c *entity.ContentBlogCategory) *pb.BlogCategoryData {
	if c == nil {
		return nil
	}
	return &pb.BlogCategoryData{
		Id:          c.ID.String(),
		AppId:       c.AppID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
	}
}

func toProtoBlog(b *entity.ContentBlog) *pb.BlogData {
	if b == nil {
		return nil
	}

	var pubStr string
	if b.PublishedAt != nil {
		pubStr = b.PublishedAt.Format(time.RFC3339)
	}

	return &pb.BlogData{
		Id:              b.ID.String(),
		AppId:           b.AppID,
		CategoryId:      b.CategoryID.String(),
		CategoryName:    b.CategoryName,
		CategorySlug:    b.CategorySlug,
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
		ViewCount:       b.ViewCount,
		PublishedAt:     pubStr,
		CreatedAt:       b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       b.UpdatedAt.Format(time.RFC3339),
	}
}

func toProtoGalleryItem(i *entity.ContentGalleryItem) *pb.GalleryItemData {
	if i == nil {
		return nil
	}
	return &pb.GalleryItemData{
		Id:            i.ID.String(),
		AppId:         i.AppID,
		GalleryId:     i.GalleryID.String(),
		Title:         i.Title,
		Caption:       i.Caption,
		MediaType:     i.MediaType,
		MediaUrl:      i.MediaURL,
		ThumbnailUrl:  i.ThumbnailURL,
		SequenceOrder: i.SequenceOrder,
		IsActive:      i.IsActive,
	}
}

func toProtoGallery(g *entity.ContentGallery) *pb.GalleryData {
	if g == nil {
		return nil
	}
	protoItems := make([]*pb.GalleryItemData, len(g.Items))
	for idx, item := range g.Items {
		protoItems[idx] = toProtoGalleryItem(&item)
	}

	return &pb.GalleryData{
		Id:          g.ID.String(),
		AppId:       g.AppID,
		Title:       g.Title,
		Slug:        g.Slug,
		Description: g.Description,
		Category:    g.Category,
		CoverUrl:    g.CoverURL,
		IsActive:    g.IsActive,
		Items:       protoItems,
		CreatedAt:   g.CreatedAt.Format(time.RFC3339),
	}
}

func toProtoCustomBlock(b *entity.ContentCustomBlock) *pb.CustomBlockData {
	if b == nil {
		return nil
	}
	metaStr := "{}"
	if b.Metadata != nil {
		if bytes, err := json.Marshal(b.Metadata); err == nil {
			metaStr = string(bytes)
		}
	}
	return &pb.CustomBlockData{
		Id:            b.ID.String(),
		AppId:         b.AppID,
		SectionKey:    b.SectionKey,
		ContentType:   b.ContentType,
		Title:         b.Title,
		Subtitle:      b.Subtitle,
		Body:          b.Body,
		ImageUrl:      b.ImageURL,
		Metadata:      metaStr,
		SequenceOrder: b.SequenceOrder,
		IsActive:      b.IsActive,
	}
}

func toProtoFAQ(f *entity.ContentFAQ) *pb.FAQData {
	if f == nil {
		return nil
	}
	return &pb.FAQData{
		Id:            f.ID.String(),
		AppId:         f.AppID,
		Category:      f.Category,
		Question:      f.Question,
		Answer:        f.Answer,
		SequenceOrder: f.SequenceOrder,
		IsActive:      f.IsActive,
	}
}

func toProtoTestimonial(t *entity.ContentTestimonial) *pb.TestimonialData {
	if t == nil {
		return nil
	}
	return &pb.TestimonialData{
		Id:            t.ID.String(),
		AppId:         t.AppID,
		CustomerName:  t.CustomerName,
		CustomerTitle: t.CustomerTitle,
		AvatarUrl:     t.AvatarURL,
		Rating:        t.Rating,
		Comment:       t.Comment,
		PackageName:   t.PackageName,
		IsFeatured:    t.IsFeatured,
		IsActive:      t.IsActive,
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
	}
}

func mapDomainErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var dErr *domainErrors.DomainError
	if errors.As(err, &dErr) {
		return status.Error(codes.InvalidArgument, dErr.Message)
	}

	switch {
	case errors.Is(err, domainErrors.ErrAppIDRequired),
		errors.Is(err, domainErrors.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, domainErrors.ErrNotFound),
		errors.Is(err, domainErrors.ErrBannerNotFound),
		errors.Is(err, domainErrors.ErrBlogCategoryNotFound),
		errors.Is(err, domainErrors.ErrBlogNotFound),
		errors.Is(err, domainErrors.ErrGalleryNotFound),
		errors.Is(err, domainErrors.ErrGalleryItemNotFound),
		errors.Is(err, domainErrors.ErrCustomBlockNotFound),
		errors.Is(err, domainErrors.ErrFAQNotFound),
		errors.Is(err, domainErrors.ErrTestimonialNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domainErrors.ErrDuplicateSlug):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, domainErrors.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())

	default:
		return status.Error(codes.Internal, fmt.Sprintf("internal service error: %v", err))
	}
}
