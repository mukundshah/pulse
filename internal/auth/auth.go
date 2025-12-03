package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

const (
	// Algorithm is the algorithm name (matches Django's Argon2PasswordHasher)
	Algorithm = "argon2"
	// Variety is the Argon2 variant (argon2id is the most secure, matches Django)
	Variety = "argon2id"
	// Argon2Version is the Argon2 version (19 = 0x13, standard version)
	Argon2Version = 19
	// DefaultTime is the default number of iterations (time cost) for Argon2 (matches Django default: 2)
	DefaultTime = 2
	// DefaultMemory is the default memory cost in KB (matches Django default: 102400 = ~100 MB)
	DefaultMemory = 102400
	// DefaultParallelism is the default number of threads (matches Django default: 8)
	DefaultParallelism = 8
	// SaltSize is the size of the salt in bytes (16 bytes recommended for Argon2)
	SaltSize = 16
	// KeyLength is the length of the derived key in bytes (32 bytes, Django default)
	KeyLength = 32
)

// decodeBase64 decodes a base64 string, handling both padded and unpadded formats
func decodeBase64(s string) ([]byte, error) {
	// Try decoding without padding first (most common for argon2-cffi)
	if data, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	// If that fails, try with padding
	return base64.StdEncoding.DecodeString(s)
}

// HashPassword hashes a password using Argon2id in Django's format
// Format: argon2$argon2id$v=<version>$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
func HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using Argon2id
	hash := argon2.IDKey([]byte(password), salt, DefaultTime, DefaultMemory, DefaultParallelism, KeyLength)

	// Encode salt and hash to base64 (standard encoding, no padding like argon2-cffi)
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	// Format matches Django's argon2-cffi output: argon2$argon2id$v=<version>$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
	return fmt.Sprintf("%s$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		Algorithm, Variety, Argon2Version, DefaultMemory, DefaultTime, DefaultParallelism, saltB64, hashB64), nil
}

// CheckPassword compares a password with a Django-formatted Argon2 hash
// Format: argon2$argon2id$v=<version>$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
func CheckPassword(password, hash string) bool {
	// Parse format: argon2$argon2id$v=<version>$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	algorithm := parts[0]
	if algorithm != Algorithm {
		return false
	}

	variety := parts[1]
	if variety != Variety && variety != "argon2i" {
		return false
	}

	// Parse version: v=<version>
	versionStr := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return false
	}
	_ = version // Version is validated but not used in Go's argon2 implementation

	// Parse parameters: m=<memory>,t=<time>,p=<parallelism>
	paramsStr := parts[3]
	var memory, time, parallelism int
	_, err = fmt.Sscanf(paramsStr, "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false
	}

	// Decode salt and hash from base64 (handle both with and without padding)
	// Django adds padding before decoding, but argon2-cffi may store without padding
	salt, err := decodeBase64(parts[4])
	if err != nil {
		return false
	}

	storedHash, err := decodeBase64(parts[5])
	if err != nil {
		return false
	}

	// Ensure stored hash has the expected length
	if len(storedHash) != KeyLength {
		return false
	}

	// Derive key using Argon2 with stored parameters
	var derivedHash []byte
	switch variety {
	case "argon2id":
		derivedHash = argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(parallelism), KeyLength)
	case "argon2i":
		derivedHash = argon2.Key([]byte(password), salt, uint32(time), uint32(memory), uint8(parallelism), KeyLength)
	default:
		return false
	}

	// Constant-time comparison
	return subtle.ConstantTimeCompare(derivedHash, storedHash) == 1
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID uuid.UUID, email, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
