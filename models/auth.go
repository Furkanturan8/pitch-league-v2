package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuthRefreshToken struct {
	bun.BaseModel `bun:"table:auth_refresh_tokens"`
	TokenID       uuid.UUID `bun:"token_id,pk"`
	UserID        int64     `bun:"user_id"`
	Role          float64   `bun:"role"`
	ExpiresAt     time.Time `bun:"expires_at"`
}

func (AuthRefreshToken) ModelName() string {
	return "user_refresh_token"
}

// AuthTokenPair defines the structure for access and refresh tokens
type AuthTokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"uid"`
	Role   float64   `json:"role"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"uid"`
	Role   float64   `json:"role"`
}

type AuthLoginVM struct {
	Email    string `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Password string `json:"password" validate:"required" label:"Parola"`
}

type AuthTokenVM struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshVM struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
