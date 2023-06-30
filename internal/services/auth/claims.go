package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
}

func NewTokenID() string {
	// Under the hood, the package never actually returns any error.
	tokenID, _ := uuid.NewUUID()
	return tokenID.String()
}

// NewClaims returns new JWT token claims
func NewClaims(subjectID UserID, tokenID string) Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "starbot",
			Subject:  strconv.FormatInt(subjectID, 10),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       tokenID,
		},
	}
}
