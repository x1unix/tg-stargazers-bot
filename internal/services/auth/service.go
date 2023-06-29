package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"go.uber.org/zap"
)

// ErrInvalidToken error occurs when token is missing, expired or not valid.
var ErrInvalidToken = errors.New("invalid token")

var jwtSignMethod = jwt.SigningMethodRS256

type TokenStorage interface {
	TokenExists(ctx context.Context, tokenID string, subjectID UserID) (bool, error)
	AddToken(ctx context.Context, tokenID string, subjectID UserID) error
}

type JWTSignParams struct {
	Method     string
	SigningKey *rsa.PublicKey
}

type Service struct {
	log        *zap.Logger
	tokenStore TokenStorage
	cfg        config.ResolvedAuthConfig
}

func NewService(log *zap.Logger, cfg config.ResolvedAuthConfig, tokenStore TokenStorage) *Service {
	return &Service{log: log, cfg: cfg, tokenStore: tokenStore}
}

func (svc Service) JWTSignParams() JWTSignParams {
	return JWTSignParams{
		Method:     jwtSignMethod.Alg(),
		SigningKey: svc.cfg.JWTPublicKey,
	}
}

func (svc Service) CreateUserToken(ctx context.Context, subject UserID) (string, error) {
	claims := NewClaims(subject)
	token := jwt.NewWithClaims(jwtSignMethod, claims)
	signedToken, err := token.SignedString(svc.cfg.JWTPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	if err := svc.tokenStore.AddToken(ctx, claims.ID, subject); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return signedToken, nil
}

func (svc Service) ValidateToken(ctx context.Context, token *Claims) (UserID, error) {
	userID, err := ParseUserID(token.Subject)
	if err != nil {
		return 0, fmt.Errorf("invalid token subject: %w", err)
	}

	ok, err := svc.tokenStore.TokenExists(ctx, token.ID, userID)
	if !ok {
		return 0, ErrInvalidToken
	}

	return userID, nil
}
