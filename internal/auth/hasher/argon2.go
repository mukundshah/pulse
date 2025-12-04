package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidHash        = errors.New("invalid hash format")
	ErrUnsupportedAlgo    = errors.New("unsupported algorithm")
)

// Argon2 parameters matching the reference implementation
const (
	defaultTimeCost    = 2
	defaultMemoryCost  = 102400 // ~100 MB
	defaultParallelism = 8
	defaultSaltSize    = 16
	defaultKeyLength   = 32
)

// Argon2Hasher handles password hashing operations.
type Argon2Hasher struct {
	time        uint32
	memory      uint32
	parallelism uint8
	saltSize    uint32
	keyLength   uint32
}

// Argon2HasherConfig configures an Argon2 Argon2Hasher.
type Argon2HasherConfig struct {
	TimeCost    uint32 // Number of iterations
	MemoryCost  uint32 // Memory in KB
	Parallelism uint8  // Number of threads
	SaltSize    uint32 // Salt length in bytes
	KeyLength   uint32 // Derived key length in bytes
}

// NewArgon2Hasher creates a new Argon2 password Argon2Hasher with default parameters.
func NewArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{
		time:        defaultTimeCost,
		memory:      defaultMemoryCost,
		parallelism: defaultParallelism,
		saltSize:    defaultSaltSize,
		keyLength:   defaultKeyLength,
	}
}

// NewArgon2HasherWithConfig creates a Argon2Hasher with custom parameters.
func NewArgon2HasherWithConfig(cfg Argon2HasherConfig) *Argon2Hasher {
	h := NewArgon2Hasher()

	if cfg.TimeCost > 0 {
		h.time = cfg.TimeCost
	}
	if cfg.MemoryCost > 0 {
		h.memory = cfg.MemoryCost
	}
	if cfg.Parallelism > 0 {
		h.parallelism = cfg.Parallelism
	}
	if cfg.SaltSize > 0 {
		h.saltSize = cfg.SaltSize
	}
	if cfg.KeyLength > 0 {
		h.keyLength = cfg.KeyLength
	}

	return h
}

// Hash creates a password hash.
// Format: argon2$argon2id$v=19$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
func (h *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.time,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	// Use raw base64 encoding (no padding)
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("argon2$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.memory, h.time, h.parallelism, saltB64, hashB64), nil
}

// Verify checks if a password matches the hash.
func (h *Argon2Hasher) Verify(password, encoded string) (bool, error) {
	params, err := parseHash(encoded)
	if err != nil {
		return false, err
	}

	var derived []byte
	switch params.variety {
	case "argon2id":
		derived = argon2.IDKey(
			[]byte(password),
			params.salt,
			params.time,
			params.memory,
			params.parallelism,
			uint32(len(params.hash)),
		)
	case "argon2i":
		derived = argon2.Key(
			[]byte(password),
			params.salt,
			params.time,
			params.memory,
			params.parallelism,
			uint32(len(params.hash)),
		)
	default:
		return false, fmt.Errorf("%w: %s", ErrUnsupportedAlgo, params.variety)
	}

	return subtle.ConstantTimeCompare(derived, params.hash) == 1, nil
}

// hashParams holds parsed hash parameters.
type hashParams struct {
	variety     string
	version     int
	memory      uint32
	time        uint32
	parallelism uint8
	salt        []byte
	hash        []byte
}

// parseHash extracts parameters from an encoded hash string.
func parseHash(encoded string) (*hashParams, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return nil, fmt.Errorf("%w: expected 6 parts, got %d", ErrInvalidHash, len(parts))
	}

	if parts[0] != "argon2" {
		return nil, fmt.Errorf("%w: expected argon2, got %s", ErrInvalidHash, parts[0])
	}

	variety := parts[1]
	if variety != "argon2id" && variety != "argon2i" {
		return nil, fmt.Errorf("%w: unsupported variety %s", ErrInvalidHash, variety)
	}

	// Parse version
	version, err := parseVersion(parts[2])
	if err != nil {
		return nil, err
	}

	// Parse parameters: m=<memory>,t=<time>,p=<parallelism>
	var memory, time, parallelism int
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism); err != nil {
		return nil, fmt.Errorf("%w: invalid parameters: %v", ErrInvalidHash, err)
	}

	// Decode salt and hash
	salt, err := decodeBase64(parts[4])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid salt: %v", ErrInvalidHash, err)
	}

	hash, err := decodeBase64(parts[5])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid hash: %v", ErrInvalidHash, err)
	}

	return &hashParams{
		variety:     variety,
		version:     version,
		memory:      uint32(memory),
		time:        uint32(time),
		parallelism: uint8(parallelism),
		salt:        salt,
		hash:        hash,
	}, nil
}

// parseVersion extracts version number from "v=19" format.
func parseVersion(s string) (int, error) {
	versionStr := strings.TrimPrefix(s, "v=")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid version: %v", ErrInvalidHash, err)
	}
	return version, nil
}

// decodeBase64 handles both padded and unpadded base64.
func decodeBase64(s string) ([]byte, error) {
	// Try raw encoding first (no padding)
	if data, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	// Fall back to standard encoding with padding
	return base64.StdEncoding.DecodeString(s)
}
