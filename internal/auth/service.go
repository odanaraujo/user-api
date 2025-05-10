package auth

import (
	"context"
	"time"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/infrastructure/exception"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"go.uber.org/zap"
)

type Service interface {
	GenerateToken(ctx context.Context, userID string) (string, *exception.Exception)
	ValidateToken(ctx context.Context, token string) (*Claims, *exception.Exception)
	InvalidateToken(ctx context.Context, token string) *exception.Exception
}

type AuthService struct {
	cache cache.Cache
}

func NewAuthService(cache cache.Cache) *AuthService {
	return &AuthService{
		cache: cache,
	}
}

func (s *AuthService) GenerateToken(ctx context.Context, userID string) (string, *exception.Exception) {
	token, err := GenerateToken(userID)
	if err != nil {
		log := loggers.FromContext(ctx)
		log.Error("failed to generate token", zap.Error(err))
		return "", exception.InternalServerException("failed to generate token")
	}

	s.cache.Set(ctx, "token:"+token, []byte(userID), 24*time.Hour)

	return token, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*Claims, *exception.Exception) {
	if _, exists := s.cache.Get(ctx, "blacklist:"+token); exists {
		return nil, exception.UnauthorizedRequestException("token has been invalidated")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		if err == ErrExpiredToken {
			return nil, exception.UnauthorizedRequestException("token has expired")
		}
		return nil, exception.UnauthorizedRequestException("invalid token")
	}

	return claims, nil
}

func (s *AuthService) InvalidateToken(ctx context.Context, token string) *exception.Exception {
	s.cache.Set(ctx, "blacklist:"+token, []byte(""), 24*time.Hour)
	return nil
}
