package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTTokenGenerator generates and validates JWT tokens.
type JWTTokenGenerator struct {
	secret   string
	validity time.Duration
}

// TokenConfig configures token generation.
type TokenConfig struct {
	Secret   string
	Validity time.Duration
}

// NewJWTTokenGenerator creates a new JWT token generator.
func NewJWTTokenGenerator(cfg TokenConfig) *JWTTokenGenerator {
	validity := cfg.Validity
	if validity == 0 {
		validity = 24 * time.Hour
	}

	return &JWTTokenGenerator{
		secret:   cfg.Secret,
		validity: validity,
	}
}

// Claims represents JWT token claims.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// Generate creates a new JWT token for a user.
func (g *JWTTokenGenerator) Generate(userID uuid.UUID, email string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(g.validity)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

// Validate verifies a JWT token and returns its claims.
func (g *JWTTokenGenerator) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(g.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
