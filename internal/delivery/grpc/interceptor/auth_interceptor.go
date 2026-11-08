package interceptor

import (
	"context"
	"strings"

	"content-service/internal/infrastructure/security"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	userContextKey  contextKey = "user_context"
	appIDContextKey contextKey = "app_id"
)

// Public methods that do not strictly require authentication (read-only endpoints)
var publicMethods = map[string]bool{
	"/content.v1.ContentService/ListBanners":        true,
	"/content.v1.ContentService/GetBanner":          true,
	"/content.v1.ContentService/ListBlogCategories": true,
	"/content.v1.ContentService/ListBlogs":          true,
	"/content.v1.ContentService/GetBlog":            true,
	"/content.v1.ContentService/ListGalleries":      true,
	"/content.v1.ContentService/GetGallery":         true,
	"/content.v1.ContentService/ListCustomBlocks":   true,
	"/content.v1.ContentService/GetCustomBlock":     true,
	"/content.v1.ContentService/ListFAQs":           true,
	"/content.v1.ContentService/ListTestimonials":   true,
}

// Mutation methods that require an admin or superadmin role.
var adminMethods = map[string]bool{
	"/content.v1.ContentService/CreateBanner":       true,
	"/content.v1.ContentService/UpdateBanner":       true,
	"/content.v1.ContentService/DeleteBanner":       true,
	"/content.v1.ContentService/CreateBlogCategory": true,
	"/content.v1.ContentService/CreateBlog":         true,
	"/content.v1.ContentService/UpdateBlog":         true,
	"/content.v1.ContentService/DeleteBlog":         true,
	"/content.v1.ContentService/CreateGallery":      true,
	"/content.v1.ContentService/DeleteGallery":      true,
	"/content.v1.ContentService/AddGalleryItem":     true,
	"/content.v1.ContentService/DeleteGalleryItem":  true,
	"/content.v1.ContentService/UpsertCustomBlock":  true,
	"/content.v1.ContentService/DeleteCustomBlock":  true,
	"/content.v1.ContentService/CreateFAQ":          true,
	"/content.v1.ContentService/UpdateFAQ":          true,
	"/content.v1.ContentService/DeleteFAQ":          true,
	"/content.v1.ContentService/CreateTestimonial":  true,
	"/content.v1.ContentService/UpdateTestimonial":  true,
	"/content.v1.ContentService/DeleteTestimonial":  true,
	"/content.v1.ContentService/UploadMedia":        true,
}

// appIDGetter is implemented by request messages that carry an explicit app_id.
type appIDGetter interface {
	GetAppId() string
}

func AuthInterceptor(jwtValidator *security.JWTValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		// 1. Extract app_id from metadata or request
		appID, _ := ExtractAppID(ctx, "")
		if appID != "" {
			ctx = context.WithValue(ctx, appIDContextKey, appID)
		}

		// 2. Extract and validate Bearer Token
		var authHeader string
		if len(md.Get("authorization")) > 0 {
			authHeader = md.Get("authorization")[0]
		}

		var userCtx *security.UserContext
		if authHeader != "" {
			u, err := jwtValidator.ValidateAccessToken(authHeader)
			if err != nil {
				// If not a public method, fail with Unauthenticated
				if !publicMethods[info.FullMethod] {
					return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token: %v", err)
				}
			} else {
				userCtx = u
				ctx = context.WithValue(ctx, userContextKey, userCtx)

				// Fallback app_id from JWT claims if not present in headers
				if appID == "" && u.AppID != "" {
					appID = u.AppID
					ctx = context.WithValue(ctx, appIDContextKey, appID)
				}
			}
		}

		// 3. Enforce authentication on protected mutation methods
		if !publicMethods[info.FullMethod] && userCtx == nil {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid authorization header")
		}

		// 4. Enforce role and tenant scope on mutation methods
		if adminMethods[info.FullMethod] {
			if !isAdminRole(userCtx) {
				return nil, status.Error(codes.PermissionDenied, "admin role required")
			}
		}
		// Superadmin may act across apps; other admins are confined to their token's app_id.
		if adminMethods[info.FullMethod] && userCtx != nil && !isSuperAdmin(userCtx) {
			if g, ok := req.(appIDGetter); ok && g.GetAppId() != "" && g.GetAppId() != userCtx.AppID {
				return nil, status.Error(codes.PermissionDenied, "app_id does not match token scope")
			}
		}

		return handler(ctx, req)
	}
}

func GetUserContext(ctx context.Context) (*security.UserContext, bool) {
	u, ok := ctx.Value(userContextKey).(*security.UserContext)
	return u, ok && u != nil
}

func IsAdmin(ctx context.Context) bool {
	u, ok := GetUserContext(ctx)
	return ok && isAdminRole(u)
}

func isAdminRole(u *security.UserContext) bool {
	if u == nil {
		return false
	}
	role := strings.ToLower(strings.TrimSpace(u.Role))
	return role == "admin" || role == "superadmin"
}

func isSuperAdmin(u *security.UserContext) bool {
	return u != nil && strings.ToLower(strings.TrimSpace(u.Role)) == "superadmin"
}

func ExtractAppID(ctx context.Context, explicitAppID string) (string, error) {
	if explicitAppID != "" {
		return strings.TrimSpace(explicitAppID), nil
	}

	if val, ok := ctx.Value(appIDContextKey).(string); ok && val != "" {
		return val, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if vals := md.Get("x-app-id"); len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
			return strings.TrimSpace(vals[0]), nil
		}
		if vals := md.Get("app-id"); len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
			return strings.TrimSpace(vals[0]), nil
		}
	}

	return "umroh", nil
}
