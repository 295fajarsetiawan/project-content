package security

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidType  = errors.New("invalid token type: access token required")
)

type UserContext struct {
	UserID uuid.UUID
	Email  string
	Role   string
	AppID  string
}

type JWTValidator struct {
	secret []byte
}

func NewJWTValidator(secret string) *JWTValidator {
	return &JWTValidator{
		secret: []byte(secret),
	}
}

type Claims struct {
	jwt.RegisteredClaims
	AppID     string `json:"app_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
}

func (v *JWTValidator) ValidateAccessToken(tokenStr string) (*UserContext, error) {
	cleanToken := strings.TrimPrefix(tokenStr, "Bearer ")
	cleanToken = strings.TrimSpace(cleanToken)

	if cleanToken == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(cleanToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verify token type if specified
	if claims.TokenType != "" && claims.TokenType != "access" {
		return nil, ErrInvalidType
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid subject UUID in token: %w", err)
	}

	return &UserContext{
		UserID: userID,
		Email:  claims.Email,
		Role:   claims.Role,
		AppID:  claims.AppID,
	}, nil
}
