package repository

import (
	"context"
	"errors"
	"github.com/personal-project/pitch-league/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type IAuthRepository interface {
	GetAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (models.AuthRefreshToken, error)
	CreateAuthRefreshToken(ctx context.Context, token models.AuthRefreshToken) error
	UpdateAuthRefreshTokenExpires(ctx context.Context, tokenID uuid.UUID, expiresAt time.Time) error
	DeleteAuthRefreshToken(ctx context.Context, userID int64) error
	GenerateTokenPair(userID int64, refreshTokenID uuid.UUID, role float64) (models.AuthTokenPair, error)
	ParseRefreshToken(refreshToken string) (refreshTokenID uuid.UUID, userID int64, role float64, err error)
}

type AuthRepository struct {
	db                     *bun.DB
	jwtSecret              string
	accessTokenExpireTime  time.Duration
	refreshTokenExpireTime time.Duration
}

func NewAuthRepository(
	db *bun.DB,
	jwtSecret string,
	accessTokenExpireTime time.Duration,
	refreshTokenExpireTime time.Duration,
) IAuthRepository {
	return &AuthRepository{
		db:                     db,
		jwtSecret:              jwtSecret,
		accessTokenExpireTime:  accessTokenExpireTime,
		refreshTokenExpireTime: refreshTokenExpireTime,
	}
}

func (r AuthRepository) GetAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (models.AuthRefreshToken, error) {
	var token models.AuthRefreshToken
	err := r.db.NewSelect().
		Model(&token).
		Where("token_id = ?", refreshTokenID).
		Scan(ctx)

	if err != nil {
		return token, errors.New("refresh token not found")
	}

	if token.ExpiresAt.Before(time.Now()) {
		return token, errors.New("refresh token expired")
	}

	return token, nil
}

func (r AuthRepository) CreateAuthRefreshToken(ctx context.Context, token models.AuthRefreshToken) error {
	_, err := r.db.NewInsert().
		Model(&token).
		Exec(ctx)

	if err != nil {
		return errors.New("failed to create auth refresh token: " + err.Error())
	}
	return nil
}

func (r AuthRepository) UpdateAuthRefreshTokenExpires(ctx context.Context, tokenID uuid.UUID, expiresAt time.Time) error {
	_, err := r.db.NewUpdate().
		Model((*models.AuthRefreshToken)(nil)).
		Set("expires_at = ?", expiresAt).
		Where("token_id = ?", tokenID).
		Exec(ctx)

	if err != nil {
		return errors.New("failed to update auth refresh token expiration: " + err.Error())
	}
	return nil
}

func (r AuthRepository) DeleteAuthRefreshToken(ctx context.Context, userID int64) error {
	_, err := r.db.NewDelete().
		Model((*models.AuthRefreshToken)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)

	if err != nil {
		return errors.New("failed to delete auth refresh token: " + err.Error())
	}
	return nil
}

func (r AuthRepository) GenerateTokenPair(userID int64, refreshTokenID uuid.UUID, role float64) (models.AuthTokenPair, error) {
	var m models.AuthTokenPair
	now := time.Now()

	accessClaims := models.AccessTokenClaims{
		ID:     uuid.New(),
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(r.accessTokenExpireTime)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(r.jwtSecret))
	if err != nil {
		return m, err
	}

	refreshClaims := models.RefreshTokenClaims{
		ID:     refreshTokenID,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(r.refreshTokenExpireTime)),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(r.jwtSecret))
	if err != nil {
		return m, err
	}

	m.AccessToken = accessToken
	m.RefreshToken = refreshToken
	return m, nil
}

func (r AuthRepository) ParseRefreshToken(refreshToken string) (refreshTokenID uuid.UUID, userID int64, role float64, err error) {
	refreshClaims := models.RefreshTokenClaims{}
	claims, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.jwtSecret), nil
	})
	if err != nil {
		return
	}
	if !claims.Valid {
		err = errors.New("invalid token")
		return
	}

	now := time.Now()
	rtokenClaims := claims.Claims.(*models.RefreshTokenClaims)
	if rtokenClaims.VerifyExpiresAt(now, false) == false {
		err = errors.New("token expired")
		return
	}

	return rtokenClaims.ID, rtokenClaims.UserID, rtokenClaims.Role, nil
}
