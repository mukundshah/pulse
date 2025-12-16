package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"

	"pulse/internal/auth/hasher"
	"pulse/internal/middleware"
	"pulse/internal/store"
)

type AccountHandler struct {
	store  *store.Store
	hasher *hasher.Argon2Hasher
}

func NewAccountHandler(s *store.Store) *AccountHandler {
	hasher := hasher.NewArgon2Hasher()
	return &AccountHandler{
		store:  s,
		hasher: hasher,
	}
}

// computeInitials extracts initials from a name
func computeInitials(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	words := strings.Fields(name)
	if len(words) == 0 {
		return ""
	}

	var initials strings.Builder

	// First letter of first word
	if len(words[0]) > 0 {
		firstRune := []rune(words[0])[0]
		if unicode.IsLetter(firstRune) {
			initials.WriteRune(unicode.ToUpper(firstRune))
		}
	}

	// First letter of last word (if different from first word)
	if len(words) > 1 && len(words[len(words)-1]) > 0 {
		lastRune := []rune(words[len(words)-1])[0]
		if unicode.IsLetter(lastRune) {
			initials.WriteRune(unicode.ToUpper(lastRune))
		}
	}

	result := initials.String()
	if result == "" && len(name) > 0 {
		// Fallback: use first character if no letters found
		firstRune := []rune(name)[0]
		return strings.ToUpper(string(firstRune))
	}

	return result
}

// getGravatarURL generates a Gravatar URL from an email address
func getGravatarURL(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.Sum([]byte(email))
	hashHex := hex.EncodeToString(hash[:])
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=200", hashHex)
}

// GetCurrentUser returns the current authenticated user's information
func (h *AccountHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Get user
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	initials := computeInitials(user.Name)
	avatarURL := getGravatarURL(user.Email)

	c.JSON(http.StatusOK, gin.H{
		"name":       user.Name,
		"email":      user.Email,
		"initials":   initials,
		"avatar_url": avatarURL,
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword allows authenticated users to change their password
func (h *AccountHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Verify old password
	ok, err := h.hasher.Verify(req.CurrentPassword, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify password"})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}

	// Hash new password
	newPasswordHash, err := h.hasher.Hash(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password
	user.PasswordHash = newPasswordHash
	if err := h.store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateProfile allows authenticated users to update their profile information
func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if email is being changed and if it's already taken
	if req.Email != user.Email {
		existingUser, err := h.store.GetUserByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		// If email changed, mark as unverified
		user.EmailVerified = false
	}

	// Update user profile
	user.Name = req.Name
	user.Email = req.Email

	if err := h.store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

// DeleteAccount deactivates the current user's account
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Verify password
	ok, err := h.hasher.Verify(req.Password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify password"})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	// Deactivate account
	user.IsActive = false

	if err := h.store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deactivated successfully"})
}
