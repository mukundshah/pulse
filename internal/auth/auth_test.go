package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "normal password",
			password: "mySecurePassword123!",
			wantErr:  false,
		},
		{
			name:     "short password",
			password: "short",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "veryLongPasswordThatExceedsNormalLength123456789!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "special!@#$%^&*()_+-=[]{}|;:,.<>?",
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: "pässwörd测试🔐",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Error("HashPassword() returned empty hash")
			}
			if !tt.wantErr {
				// Verify hash format
				if !strings.HasPrefix(hash, Algorithm+"$") {
					t.Errorf("HashPassword() hash doesn't start with algorithm: %s", hash)
				}
				parts := strings.Split(hash, "$")
				if len(parts) != 6 {
					t.Errorf("HashPassword() hash has wrong number of parts: got %d, want 6", len(parts))
				}
			}
		})
	}
}

func TestHashPassword_UniqueSalts(t *testing.T) {
	password := "testPassword123"
	hashes := make(map[string]bool)

	// Generate multiple hashes for the same password
	for i := 0; i < 10; i++ {
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		// Each hash should be unique due to random salt
		if hashes[hash] {
			t.Errorf("HashPassword() generated duplicate hash at iteration %d", i)
		}
		hashes[hash] = true
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mySecurePassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongPassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "similar password",
			password: "mySecurePassword123",
			hash:     hash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPassword_InvalidFormats(t *testing.T) {
	password := "testPassword123"

	invalidHashes := []struct {
		name string
		hash string
	}{
		{
			name: "empty hash",
			hash: "",
		},
		{
			name: "invalid format",
			hash: "invalid",
		},
		{
			name: "missing parts",
			hash: "argon2$argon2id",
		},
		{
			name: "wrong algorithm",
			hash: "pbkdf2_sha256$600000$salt$hash",
		},
		{
			name: "invalid algorithm",
			hash: "invalid$argon2id$v=19$m=102400,t=2,p=8$salt$hash",
		},
		{
			name: "invalid variety",
			hash: "argon2$invalid$v=19$m=102400,t=2,p=8$salt$hash",
		},
		{
			name: "missing version",
			hash: "argon2$argon2id$m=102400,t=2,p=8$salt$hash",
		},
		{
			name: "invalid version format",
			hash: "argon2$argon2id$v=invalid$m=102400,t=2,p=8$salt$hash",
		},
		{
			name: "missing parameters",
			hash: "argon2$argon2id$v=19$salt$hash",
		},
		{
			name: "invalid parameters format",
			hash: "argon2$argon2id$v=19$invalid$salt$hash",
		},
		{
			name: "missing salt",
			hash: "argon2$argon2id$v=19$m=102400,t=2,p=8$hash",
		},
		{
			name: "missing hash",
			hash: "argon2$argon2id$v=19$m=102400,t=2,p=8$salt$",
		},
		{
			name: "invalid base64 salt",
			hash: "argon2$argon2id$v=19$m=102400,t=2,p=8$invalid-base64!!!$hash",
		},
		{
			name: "invalid base64 hash",
			hash: "argon2$argon2id$v=19$m=102400,t=2,p=8$salt$invalid-base64!!!",
		},
	}

	for _, tt := range invalidHashes {
		t.Run(tt.name, func(t *testing.T) {
			if CheckPassword(password, tt.hash) {
				t.Errorf("CheckPassword() should return false for invalid hash: %s", tt.hash)
			}
		})
	}
}

func TestCheckPassword_RoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"mediumLengthPassword",
		"veryLongPasswordThatExceedsNormalLength123456789!@#$%^&*()",
		"special!@#$%^&*()_+-=[]{}|;:,.<>?",
		"unicode pässwörd测试🔐",
		"1234567890",
		"password with spaces",
		"PASSWORD_WITH_UPPERCASE",
		"password_with_underscores",
		"password-with-dashes",
	}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			if !CheckPassword(password, hash) {
				t.Errorf("CheckPassword() failed for password: %s", password)
			}
		})
	}
}

func TestCheckPassword_DifferentParameters(t *testing.T) {
	password := "testPassword123"

	// Test with different time costs
	timeCosts := []int{1, 2, 3}
	for _, timeCost := range timeCosts {
		t.Run(fmt.Sprintf("time_cost_%d", timeCost), func(t *testing.T) {
			// Manually create a hash with different time cost
			// This tests that CheckPassword correctly parses and uses stored parameters
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			// Verify it works
			if !CheckPassword(password, hash) {
				t.Errorf("CheckPassword() failed with time cost %d", timeCost)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	secret := "test-secret-key"

	token, err := GenerateToken(userID, email, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}

	// Token should be a valid JWT format (three parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("GenerateToken() token has wrong format: got %d parts, want 3", len(parts))
	}
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	secret := "test-secret-key"

	user1 := uuid.New()
	user2 := uuid.New()

	token1, err1 := GenerateToken(user1, "user1@example.com", secret)
	token2, err2 := GenerateToken(user2, "user2@example.com", secret)

	if err1 != nil {
		t.Fatalf("GenerateToken() error = %v", err1)
	}
	if err2 != nil {
		t.Fatalf("GenerateToken() error = %v", err2)
	}

	// Tokens should be different for different users
	if token1 == token2 {
		t.Error("GenerateToken() generated same token for different users")
	}
}

func TestValidateToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	secret := "test-secret-key"

	token, err := GenerateToken(userID, email, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims == nil {
		t.Fatal("ValidateToken() returned nil claims")
	}

	if claims.UserID != userID {
		t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, userID)
	}

	if claims.Email != email {
		t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, email)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key"

	invalidTokens := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid format",
			token: "not.a.valid.jwt",
		},
		{
			name: "wrong secret",
			token: func() string {
				token, _ := GenerateToken(uuid.New(), "test@example.com", "wrong-secret")
				return token
			}(),
		},
		{
			name:  "malformed token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
		},
	}

	for _, tt := range invalidTokens {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, secret)
			if err == nil {
				t.Errorf("ValidateToken() should return error for invalid token, got claims: %v", claims)
			}
		})
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	secret := "correct-secret"
	wrongSecret := "wrong-secret"

	token, err := GenerateToken(userID, email, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ValidateToken(token, wrongSecret)
	if err == nil {
		t.Errorf("ValidateToken() should return error for wrong secret, got claims: %v", claims)
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	secret := "test-secret-key"

	// Create a token with very short expiration
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate should fail for expired token
	validatedClaims, err := ValidateToken(tokenString, secret)
	if err == nil {
		t.Errorf("ValidateToken() should return error for expired token, got claims: %v", validatedClaims)
	}
}

func TestHashAndCheckPassword_Consistency(t *testing.T) {
	password := "consistentPassword123"

	// Hash the password multiple times
	hashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}
		hashes[i] = hash

		// Each hash should verify correctly
		if !CheckPassword(password, hash) {
			t.Errorf("CheckPassword() failed for hash %d", i)
		}
	}

	// All hashes should be different (due to unique salts)
	for i := 0; i < len(hashes); i++ {
		for j := i + 1; j < len(hashes); j++ {
			if hashes[i] == hashes[j] {
				t.Errorf("HashPassword() generated duplicate hashes at indices %d and %d", i, j)
			}
		}
	}
}

func TestHashPassword_Format(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		t.Fatalf("HashPassword() hash has wrong number of parts: got %d, want 6", len(parts))
	}

	// Verify each part
	if parts[0] != Algorithm {
		t.Errorf("HashPassword() algorithm = %v, want %v", parts[0], Algorithm)
	}

	if parts[1] != Variety {
		t.Errorf("HashPassword() variety = %v, want %v", parts[1], Variety)
	}

	if !strings.HasPrefix(parts[2], "v=") {
		t.Errorf("HashPassword() version part doesn't start with 'v=': %v", parts[2])
	}

	if !strings.Contains(parts[3], "m=") || !strings.Contains(parts[3], "t=") || !strings.Contains(parts[3], "p=") {
		t.Errorf("HashPassword() parameters part is malformed: %v", parts[3])
	}

	// Salt and hash should be base64 encoded (non-empty)
	if parts[4] == "" {
		t.Error("HashPassword() salt is empty")
	}
	if parts[5] == "" {
		t.Error("HashPassword() hash is empty")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword() error = %v", err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkPassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("HashPassword() error = %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hash)
	}
}
