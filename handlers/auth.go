package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
	"github.com/personal-project/pitch-league/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authRepository         repository.IAuthRepository
	userRepository         repository.IUserRepository
	refreshTokenExpireTime time.Duration
	jwtSecret              string
}

func NewAuthHandler(ar repository.IAuthRepository, ur repository.IUserRepository, refreshTokenExpireTime time.Duration, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authRepository:         ar,
		userRepository:         ur,
		refreshTokenExpireTime: refreshTokenExpireTime,
		jwtSecret:              jwtSecret,
	}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var vm models.AuthLoginVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	user, err := h.userRepository.GetByEmail(ctx.Context(), utils.CleanEmail(vm.Email))
	if err != nil {
		return errorResult(ctx, err)
	}

	ok := utils.CheckPasswordHash(strings.TrimSpace(vm.Password), user.Password)
	if !ok {
		return errorResult(ctx, errors.New("hatalı email veya parola"))
	}

	refreshTokenID := uuid.New()
	tokens, err := h.authRepository.GenerateTokenPair(user.ID, refreshTokenID, float64(user.Role))
	if err != nil {
		return errorResult(ctx, err)
	}

	refreshToken := models.AuthRefreshToken{
		TokenID:   refreshTokenID,
		UserID:    user.ID,
		Role:      float64(user.Role),
		ExpiresAt: time.Now().Add(h.refreshTokenExpireTime),
	}

	err = h.authRepository.CreateAuthRefreshToken(ctx.Context(), refreshToken)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, models.AuthTokenVM{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	var vm models.AuthRefreshVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	refreshTokenID, userID, role, err := h.authRepository.ParseRefreshToken(vm.RefreshToken)
	if err != nil {
		return errorResult(ctx, errors.New("yetkisiz: "+err.Error()))
	}

	authRefreshToken, err := h.authRepository.GetAuthRefreshToken(ctx.Context(), refreshTokenID)
	if err != nil {
		return errorResult(ctx, err)
	}

	newTokenPair, err := h.authRepository.GenerateTokenPair(userID, refreshTokenID, role)
	if err != nil {
		return errorResult(ctx, err)
	}

	err = h.authRepository.UpdateAuthRefreshTokenExpires(ctx.Context(), authRefreshToken.TokenID, time.Now().Add(h.refreshTokenExpireTime))
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, models.AuthTokenVM{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
	})
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	// Token'ı al
	token := ctx.Get("Authorization")
	if token == "" {
		return errorResult(ctx, errors.New("token bulunamadı"))
	}

	// "Bearer " prefix'ini kaldır
	token = strings.TrimPrefix(token, "Bearer ")

	// Token'ı parse et
	claims := &models.AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil {
		return errorResult(ctx, errors.New("geçersiz token"))
	}

	// Kullanıcının tüm refresh token'larını sil
	err = h.authRepository.DeleteAuthRefreshToken(ctx.Context(), claims.UserID)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Başarıyla çıkış yapıldı")
}

// ... diğer handler metodları ve yardımcı fonksiyonlar buraya gelecek ...
