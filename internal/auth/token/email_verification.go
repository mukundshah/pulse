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

// EmailVerificationTokenGenerator generates and validates email verification tokens.
// Uses HMAC-based tokens (no database storage needed) similar to Django AllAuth.
type EmailVerificationTokenGenerator struct {
	secret  string
	timeout time.Duration
}

// EmailVerificationConfig holds configuration for creating a generator.
type EmailVerificationConfig struct {
	Secret  string
	Timeout time.Duration
}

// NewEmailVerificationTokenGenerator creates a new email verification token generator.
// The timeout specifies how long tokens remain valid (default: 7 days).
func NewEmailVerificationTokenGenerator(cfg EmailVerificationConfig) *EmailVerificationTokenGenerator {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 7 * 24 * time.Hour // Default: 7 days
	}

	return &EmailVerificationTokenGenerator{
		secret:  cfg.Secret,
		timeout: timeout,
	}
}

// Generate creates an email verification token for the given email address.
// Returns a token in the format: {emailb64}-{timestamp}-{hmac_signature}
// The email is included in the token so it can be extracted during verification.
func (g *EmailVerificationTokenGenerator) Generate(email string) string {
	ts := secondsSinceEpoch(time.Now())
	token := g.generateWithTimestamp(email, ts)
	emailb64 := base64.RawURLEncoding.EncodeToString([]byte(email))
	return fmt.Sprintf("%s-%s", emailb64, token)
}

// Validate checks if a token is valid for the given email address.
// Returns true if the token is valid and not expired.
func (g *EmailVerificationTokenGenerator) Validate(email, token string) bool {
	if token == "" || email == "" {
		return false
	}

	// Extract email and validate signature
	extractedEmail, tokenPart, err := g.splitToken(token)
	if err != nil {
		return false
	}

	// Verify email matches
	if extractedEmail != email {
		return false
	}

	// Check expiration
	ts, err := parseEmailToken(tokenPart)
	if err != nil {
		return false
	}

	now := secondsSinceEpoch(time.Now())
	elapsed := time.Duration(now-ts) * time.Second

	return elapsed <= g.timeout
}

// ValidateWithEmail validates a token and extracts the email from it.
// Returns true if the token is valid and not expired.
// The email is extracted from the token itself.
func (g *EmailVerificationTokenGenerator) ValidateWithEmail(token string) bool {
	if token == "" {
		return false
	}

	// Extract email and validate signature
	_, tokenPart, err := g.splitToken(token)
	if err != nil {
		return false
	}

	// Check expiration
	ts, err := parseEmailToken(tokenPart)
	if err != nil {
		return false
	}

	now := secondsSinceEpoch(time.Now())
	elapsed := time.Duration(now-ts) * time.Second

	return elapsed <= g.timeout
}

// GetEmailFromToken extracts the email from a token.
// Returns the email address encoded in the token.
func (g *EmailVerificationTokenGenerator) GetEmailFromToken(token string) (string, error) {
	email, _, err := g.splitToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token format: %w", err)
	}
	return email, nil
}

// splitToken splits a token into email and token parts.
// Token format: {emailb64}-{timestamp}-{hash}
func (g *EmailVerificationTokenGenerator) splitToken(token string) (string, string, error) {
	parts := strings.SplitN(token, "-", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid token format")
	}

	emailb64 := parts[0]
	tokenPart := parts[1]

	emailBytes, err := base64.RawURLEncoding.DecodeString(emailb64)
	if err != nil {
		return "", "", fmt.Errorf("invalid email encoding: %w", err)
	}

	email := string(emailBytes)

	// Validate the token part by extracting email and regenerating expected token
	ts, err := parseEmailToken(tokenPart)
	if err != nil {
		return "", "", fmt.Errorf("invalid token: %w", err)
	}

	// Generate expected token with extracted email and timestamp
	expected := g.generateWithTimestamp(email, ts)

	// Constant-time comparison to verify HMAC signature
	if subtle.ConstantTimeCompare([]byte(expected), []byte(tokenPart)) != 1 {
		return "", "", fmt.Errorf("invalid token signature")
	}

	return email, tokenPart, nil
}

// generateWithTimestamp creates a token with a specific timestamp.
func (g *EmailVerificationTokenGenerator) generateWithTimestamp(email string, ts int64) string {
	// Convert timestamp to base36 for compact representation
	tsB36 := strconv.FormatInt(ts, 36)

	// Create hash value: email + timestamp
	hashValue := fmt.Sprintf("%s:%d", email, ts)
	hashBytes := saltedHMACEmail(hashValue, g.secret)

	// Convert to base64 URL-safe encoding
	hashStr := base64.RawURLEncoding.EncodeToString(hashBytes)

	// Return: {timestamp}-{hash}
	return fmt.Sprintf("%s-%s", tsB36, hashStr)
}

// parseEmailToken extracts the timestamp from a token string.
func parseEmailToken(token string) (int64, error) {
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

// saltedHMACEmail creates a salted HMAC hash for email verification.
func saltedHMACEmail(value, secret string) []byte {
	const keySalt = "email.verification.token.generator"

	// Derive key: SHA256(key_salt + secret)
	keyHasher := sha256.New()
	keyHasher.Write([]byte(keySalt + secret))
	key := keyHasher.Sum(nil)

	// Create HMAC with derived key
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))

	return mac.Sum(nil)
}
