package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type ResolvedAuthConfig struct {
	JWTPrivateKey *rsa.PrivateKey
	JWTPublicKey  *rsa.PublicKey
}

type AuthConfig struct {
	JWTSecretFile string `envconfig:"JWT_SECRET_KEY_FILE" required:"true"`
	JWTPublicKey  string `envconfig:"JWT_PUBLIC_KEY_FILE" required:"true"`
}

func (cfg AuthConfig) ResolvedAuthConfig() (*ResolvedAuthConfig, error) {
	privKey, err := readJWTSecretFile(cfg.JWTSecretFile)
	if err != nil {
		return nil, err
	}

	pubKey, err := readJWTPublicKeyFile(cfg.JWTPublicKey)
	if err != nil {
		return nil, err
	}

	if privKey == nil || pubKey == nil {
		return nil, errors.New("invalid JWT config")
	}

	return &ResolvedAuthConfig{
		JWTPublicKey:  pubKey,
		JWTPrivateKey: privKey,
	}, nil
}

// readJWTSecretFile reads JWT secret key from PEM/RSA file.
func readJWTSecretFile(fileName string) (*rsa.PrivateKey, error) {
	if fileName == "" {
		return nil, errors.New("empty JWT private key path")
	}

	keyFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWT private key: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT private key: %w", err)
	}

	return privKey, nil
}

// readJWTPublicKeyFile reads JWT public key from PEM/RSA file.
func readJWTPublicKeyFile(fileName string) (*rsa.PublicKey, error) {
	if fileName == "" {
		return nil, errors.New("empty JWT public key path")
	}

	keyFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWT public key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT public key: %w", err)
	}

	return pubKey, nil
}
