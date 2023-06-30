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
	// AddToken stores a new user token or overwrites existing one.
	AddToken(ctx context.Context, tokenID string, subjectID UserID) error

	// GetToken returns stored token.
	//
	// Return ErrTokenNotExists if token doesn't exist.
	GetToken(ctx context.Context, subjectID UserID) (string, error)

	// RemoveToken removes user token from storage.
	RemoveToken(ctx context.Context, subjectID UserID) error
}

type JWTSignParams struct {
	Method     string
	SigningKey *rsa.PublicKey
}

type Service struct {
	log        *zap.Logger
	cfg        config.ResolvedAuthConfig
	tokenStore TokenStorage
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

// ValidateToken checks if provided token is correct.
func (svc Service) ValidateToken(ctx context.Context, token *Claims) (UserID, error) {
	userID, err := ParseUserID(token.Subject)
	if err != nil {
		return 0, fmt.Errorf("invalid token subject: %w", err)
	}

	gotToken, err := svc.tokenStore.GetToken(ctx, userID)
	if errors.Is(err, ErrTokenNotExists) {
		return 0, ErrInvalidToken
	}

	if err != nil {
		return 0, err
	}

	if gotToken != token.ID {
		return 0, ErrInvalidToken
	}

	return userID, nil
}

// ProvideUserToken returns stored user token or creates a new one.
func (svc Service) ProvideUserToken(ctx context.Context, subject UserID) (string, error) {
	token, err := svc.tokenStore.GetToken(ctx, subject)
	if errors.Is(err, ErrTokenNotExists) {
		token, err = svc.generateNewToken(ctx, subject)
	}
	if err != nil {
		return "", err
	}

	claims := NewClaims(subject, token)
	return svc.buildToken(claims)
}

func (svc Service) RemoveUserToken(ctx context.Context, subject UserID) error {
	return svc.tokenStore.RemoveToken(ctx, subject)
}

func (svc Service) generateNewToken(ctx context.Context, subject UserID) (string, error) {
	tokenID := NewTokenID()
	if err := svc.tokenStore.AddToken(ctx, tokenID, subject); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return tokenID, nil
}

func (svc Service) buildToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwtSignMethod, claims)
	signedToken, err := token.SignedString(svc.cfg.JWTPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, err
}
