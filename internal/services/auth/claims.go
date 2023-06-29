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

// NewClaims returns new JWT token claims
func NewClaims(subjectID UserID) Claims {
	// Under the hood, the package never actually returns any error.
	jwtID, _ := uuid.NewUUID()

	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "starbot",
			Subject:  strconv.FormatInt(subjectID, 10),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       jwtID.String(),
		},
	}
}
