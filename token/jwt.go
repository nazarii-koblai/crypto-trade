package token

import (
	"fmt"

	"github.com/crypto-trade/config"
	"github.com/golang-jwt/jwt"
)

// Token describes token interface.
type Token interface {
	GenerateWithClaims(calims jwt.MapClaims) (string, error)
}

// JWT represents jwt structure.
type JWT struct {
	key []byte
}

// NewJWT returns new JWT token structure.
func NewJWT(config config.JWT) *JWT {
	return &JWT{
		key: config.Key,
	}
}

// GenerateWithClaims generates a new token using HS256 algh.
func (j *JWT) GenerateWithClaims(calims jwt.MapClaims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, calims).SignedString(j.key)
	if err != nil {
		return "", fmt.Errorf("can't generate token")
	}
	return token, nil
}
