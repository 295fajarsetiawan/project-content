package interceptor

import (
	"context"
	"testing"
	"time"

	"content-service/internal/infrastructure/security"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const testSecret = "test-secret-for-interceptor"

const deleteGalleryItemMethod = "/content.v1.ContentService/DeleteGalleryItem"

type appReq struct {
	AppID string
}

func (r appReq) GetAppId() string { return r.AppID }

func makeToken(t *testing.T, role, appID string) string {
	t.Helper()
	claims := security.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		AppID:     appID,
		Email:     "user@example.com",
		Role:      role,
		TokenType: "access",
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return "Bearer " + tok
}

func invoke(t *testing.T, authHeader string, req interface{}, method string) error {
	t.Helper()
	ic := AuthInterceptor(security.NewJWTValidator(testSecret))
	ctx := context.Background()
	if authHeader != "" {
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", authHeader))
	}
	info := &grpc.UnaryServerInfo{FullMethod: method}
	_, err := ic(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})
	return err
}

func TestAuthInterceptor_AdminMethodRolePolicy(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		wantCode codes.Code
	}{
		{name: "no token is unauthenticated", header: "", wantCode: codes.Unauthenticated},
		{name: "regular user is denied", header: makeToken(t, "user", "umroh"), wantCode: codes.PermissionDenied},
		{name: "admin is allowed", header: makeToken(t, "admin", "umroh"), wantCode: codes.OK},
		{name: "superadmin is allowed", header: makeToken(t, "superadmin", "umroh"), wantCode: codes.OK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := invoke(t, tt.header, appReq{AppID: "umroh"}, deleteGalleryItemMethod)
			if got := status.Code(err); got != tt.wantCode {
				t.Fatalf("code = %v, want %v (err=%v)", got, tt.wantCode, err)
			}
		})
	}
}

func TestAuthInterceptor_AdminCannotActOnOtherApp(t *testing.T) {
	admin := makeToken(t, "admin", "umroh")

	if err := invoke(t, admin, appReq{AppID: "other"}, deleteGalleryItemMethod); status.Code(err) != codes.PermissionDenied {
		t.Fatalf("admin cross-app code = %v, want PermissionDenied", status.Code(err))
	}
	if err := invoke(t, admin, appReq{AppID: "umroh"}, deleteGalleryItemMethod); status.Code(err) != codes.OK {
		t.Fatalf("admin same-app code = %v, want OK", status.Code(err))
	}
}

func TestAuthInterceptor_SuperAdminCanActAcrossApps(t *testing.T) {
	super := makeToken(t, "superadmin", "umroh")
	if err := invoke(t, super, appReq{AppID: "other"}, deleteGalleryItemMethod); status.Code(err) != codes.OK {
		t.Fatalf("superadmin cross-app code = %v, want OK", status.Code(err))
	}
}
