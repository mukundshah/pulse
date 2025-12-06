package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/auth/hasher"
	"pulse/internal/auth/token"
	"pulse/internal/config"
	"pulse/internal/email"
	"pulse/internal/models"
	"pulse/internal/store"
)

type AuthHandler struct {
	store              *store.Store
	config             *config.Config
	hasher             *hasher.Argon2Hasher
	jwtGenerator       *token.JWTTokenGenerator
	passwordResetToken *token.PasswordResetTokenGenerator
	emailService       *email.Service
}

func NewAuthHandler(s *store.Store, cfg *config.Config, emailService *email.Service) *AuthHandler {

	hasher := hasher.NewArgon2Hasher()

	jwtGen := token.NewJWTTokenGenerator(token.TokenConfig{
		Secret:   cfg.JWTSecret,
		Validity: 24 * time.Hour,
	})

	passwordResetGen := token.NewPasswordResetTokenGenerator(token.Config{
		Secret:  cfg.JWTSecret,
		Timeout: time.Duration(cfg.PasswordResetTimeout) * time.Second,
	})

	return &AuthHandler{
		store:              s,
		config:             cfg,
		hasher:             hasher,
		jwtGenerator:       jwtGen,
		passwordResetToken: passwordResetGen,
		emailService:       emailService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	existingUser, err := h.store.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
		return
	}

	// Hash password
	passwordHash, err := h.hasher.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user
	user := &models.User{
		ID:            uuid.New(),
		Name:          req.Name,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Generate token
	jwtToken, err := h.jwtGenerator.Generate(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, AuthResponse{
		Token: jwtToken,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
		return
	}

	// Verify password
	ok, err := h.hasher.Verify(req.Password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify password"})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := h.store.UpdateUser(user); err != nil {
		// Log error but don't fail the login
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update last login"})
		return
	}

	// Generate token
	jwtToken, err := h.jwtGenerator.Generate(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, AuthResponse{
		Token: jwtToken,
		User:  user,
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// ForgotPassword generates a password reset token for a user
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not (security best practice)
		c.JSON(http.StatusOK, gin.H{"message": "if the email exists, a password reset link has been sent"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusOK, gin.H{"message": "if the email exists, a password reset link has been sent"})
		return
	}

	// Generate reset token
	resetToken := h.passwordResetToken.GenerateWithUID(token.User{
		ID:           user.ID.String(),
		PasswordHash: user.PasswordHash,
		LastLogin:    user.LastLogin,
		Email:        user.Email,
	})

	// Send password reset email asynchronously (non-blocking)
	h.emailService.SendPasswordResetEmailAsync(user.Email, resetToken)

	// Always return success to avoid revealing if user exists (security best practice)
	c.JSON(http.StatusOK, gin.H{"message": "if the email exists, a password reset link has been sent"})
}

// ResetPassword resets a user's password using a reset token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := h.passwordResetToken.GetUID(req.Token)

	// Get user by email
	user, err := h.store.GetUserByID(uuid.MustParse(uid))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token or user not found"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
		return
	}

	// Validate token
	if !h.passwordResetToken.ValidateWithUID(token.User{
		ID:           user.ID.String(),
		PasswordHash: user.PasswordHash,
		LastLogin:    user.LastLogin,
		Email:        user.Email,
	}, req.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// Hash new password
	newPasswordHash, err := h.hasher.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password (this will invalidate the token since password hash changes)
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now()
	if err := h.store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
