package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var epoch = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

// Generator generates and validates password reset tokens.
type PasswordResetTokenGenerator struct {
	secret  string
	timeout time.Duration
}

// Config holds configuration for creating a Generator.
type Config struct {
	Secret  string
	Timeout time.Duration
}

// New creates a new password reset token generator.
// The timeout specifies how long tokens remain valid.
func NewPasswordResetTokenGenerator(cfg Config) *PasswordResetTokenGenerator {
	return &PasswordResetTokenGenerator{
		secret:  cfg.Secret,
		timeout: cfg.Timeout,
	}
}

// User represents the minimal user data needed for token generation.
type User struct {
	ID           string     // User's primary key as string
	PasswordHash string     // Password hash
	LastLogin    *time.Time // Pointer to allow nil for users who never logged in
	Email        string     // User's email address
}

// Generate creates a password reset token for the given user.
func (g *PasswordResetTokenGenerator) Generate(u User) string {
	ts := secondsSinceEpoch(time.Now())
	return g.generateWithTimestamp(u, ts)
}

func (g *PasswordResetTokenGenerator) GetUID(token string) string {
	parts := strings.SplitN(token, "-", 3)
	if len(parts) != 3 {
		return ""
	}
	uidb64 := parts[0]
	uid, err := base64.RawURLEncoding.DecodeString(uidb64)
	if err != nil {
		return ""
	}
	return string(uid)
}

func (g *PasswordResetTokenGenerator) GenerateWithUID(u User) string {
	ts := secondsSinceEpoch(time.Now())
	token := g.generateWithTimestamp(u, ts)
	uidb64 := base64.RawURLEncoding.EncodeToString([]byte(u.ID))
	return fmt.Sprintf("%s-%s", uidb64, token)
}

func (g *PasswordResetTokenGenerator) ValidateWithUID(u User, token string) bool {
	parts := strings.SplitN(token, "-", 2)
	if len(parts) != 2 {
		return false
	}
	uidb64 := parts[0]
	token = parts[1]
	uid, err := base64.RawURLEncoding.DecodeString(uidb64)

	if err != nil {
		return false
	}

	if string(uid) != u.ID {
		return false
	}

	return g.ValidateWithFallbacks(u, token, nil)
}

// Validate checks if a token is valid for the given user.
// Returns true if the token is valid and not expired.
func (g *PasswordResetTokenGenerator) Validate(u User, token string) bool {
	return g.ValidateWithFallbacks(u, token, nil)
}

// ValidateWithFallbacks validates a token against multiple secrets.
// Useful when rotating secrets to support tokens generated with old secrets.
func (g *PasswordResetTokenGenerator) ValidateWithFallbacks(u User, token string, fallbackSecrets []string) bool {
	if token == "" {
		return false
	}

	ts, err := parseToken(token)
	if err != nil {
		return false
	}

	// Try primary secret first, then fallbacks
	secrets := append([]string{g.secret}, fallbackSecrets...)
	tokenValid := false

	for _, secret := range secrets {
		expected := g.generateWithSecret(u, ts, secret)
		if subtle.ConstantTimeCompare([]byte(expected), []byte(token)) == 1 {
			tokenValid = true
			break
		}
	}

	if !tokenValid {
		return false
	}

	// Check if token has expired
	now := secondsSinceEpoch(time.Now())
	elapsed := time.Duration(now-ts) * time.Second

	return elapsed <= g.timeout
}

// generateWithTimestamp creates a token with a specific timestamp and the default secret.
func (g *PasswordResetTokenGenerator) generateWithTimestamp(u User, ts int64) string {
	return g.generateWithSecret(u, ts, g.secret)
}

// generateWithSecret creates a token with a specific timestamp and secret.
func (g *PasswordResetTokenGenerator) generateWithSecret(u User, ts int64, secret string) string {
	// Convert timestamp to base36 for compact representation
	tsB36 := strconv.FormatInt(ts, 36)

	// Create hash of user data
	hashValue := makeHashValue(u, ts)
	hashBytes := saltedHMAC(hashValue, secret)

	// Convert to hex and take every other character
	hexHash := fmt.Sprintf("%x", hashBytes)
	hashStr := takeEveryOther(hexHash)

	return fmt.Sprintf("%s-%s", tsB36, hashStr)
}

// parseToken extracts the timestamp from a token string.
func parseToken(token string) (int64, error) {
	parts := strings.SplitN(token, "-", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}

	ts, err := strconv.ParseInt(parts[0], 36, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp: %w", err)
	}

	return ts, nil
}

// makeHashValue creates the hash input string from user data.
// Format: "{id}{password}{login_timestamp}{timestamp}{email}"
func makeHashValue(u User, ts int64) string {
	loginStr := formatLoginTimestamp(u.LastLogin)
	return fmt.Sprintf("%s%s%s%d%s", u.ID, u.PasswordHash, loginStr, ts, u.Email)
}

// formatLoginTimestamp formats the last login time.
// Returns empty string for nil (user never logged in).
func formatLoginTimestamp(t *time.Time) string {
	if t == nil {
		return ""
	}

	// Truncate to seconds (no microseconds) in UTC
	truncated := time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		0, time.UTC,
	)

	return truncated.Format("2006-01-02 15:04:05")
}

// saltedHMAC creates a salted HMAC hash.
// First derives a key from the salt and secret, then creates an HMAC.
func saltedHMAC(value, secret string) []byte {
	const keySalt = "password.reset.token.generator"

	// Derive key: SHA256(key_salt + secret)
	keyHasher := sha256.New()
	keyHasher.Write([]byte(keySalt + secret))
	key := keyHasher.Sum(nil)

	// Create HMAC with derived key
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))

	return mac.Sum(nil)
}

// takeEveryOther returns every other character from the string.
func takeEveryOther(s string) string {
	var b strings.Builder
	b.Grow(len(s) / 2)

	for i := 0; i < len(s); i += 2 {
		b.WriteByte(s[i])
	}

	return b.String()
}

// secondsSinceEpoch returns seconds since the epoch (2001-01-01).
func secondsSinceEpoch(t time.Time) int64 {
	return int64(t.Sub(epoch).Seconds())
}
