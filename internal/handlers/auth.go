package handlers

import (
	"context"
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
	store                  *store.Store
	config                 *config.Config
	hasher                 *hasher.Argon2Hasher
	jwtGenerator           *token.JWTTokenGenerator
	passwordResetToken     *token.PasswordResetTokenGenerator
	emailVerificationToken *token.EmailVerificationTokenGenerator
	emailService           *email.Service
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

	emailVerificationGen := token.NewEmailVerificationTokenGenerator(token.EmailVerificationConfig{
		Secret:  cfg.JWTSecret,
		Timeout: 7 * 24 * time.Hour, // 7 days
	})

	return &AuthHandler{
		store:                  s,
		config:                 cfg,
		hasher:                 hasher,
		jwtGenerator:           jwtGen,
		passwordResetToken:     passwordResetGen,
		emailVerificationToken: emailVerificationGen,
		emailService:           emailService,
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

type RegisterResponse struct {
	Message string `json:"message"`
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

	// Generate email verification token
	verificationToken := h.emailVerificationToken.Generate(user.Email)

	// Send verification email asynchronously
	h.emailService.SendEmailVerificationAsync(context.Background(), user.Email, verificationToken)

	// Return success message (don't authenticate yet)
	c.JSON(http.StatusCreated, RegisterResponse{
		Message: "Registration successful. Please check your email to verify your account.",
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

	// Check if email is verified
	if !user.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "email not verified. Please check your email for verification link."})
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

	// Generate token with JTI
	jwtToken, jti, err := h.jwtGenerator.Generate(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Create session
	session := &models.Session{
		UserID:       user.ID,
		JTI:          jti,
		UserAgent:    c.GetHeader("User-Agent"),
		IPAddress:    c.ClientIP(),
		IsActive:     true,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // Match JWT validity
		LastActivity: time.Now(),
	}

	if err := h.store.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
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
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=8"`
	LogoutFromAllDevices bool   `json:"logout_from_all_devices" binding:"omitempty"`
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
	h.emailService.SendPasswordResetEmailAsync(context.Background(), user.Email, resetToken)

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

	if req.LogoutFromAllDevices {
		// Invalidate all existing sessions for security (user must log in again)
		if err := h.store.InvalidateAllUserSessions(user.ID); err != nil {
			// Log error but don't fail the password reset
			// The password has been changed, which is the critical operation
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmail verifies a user's email address using an HMAC-based token
// The email is extracted from the token itself, so it doesn't need to be provided.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract email from token
	email, err := h.emailVerificationToken.GetEmailFromToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
		return
	}

	// Validate token (this also checks expiration)
	if !h.emailVerificationToken.ValidateWithEmail(req.Token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired verification token"})
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if already verified
	if user.EmailVerified {
		c.JSON(http.StatusOK, gin.H{"message": "email already verified"})
		return
	}

	// Mark email as verified
	user.EmailVerified = true
	user.UpdatedAt = time.Now()
	if err := h.store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendVerificationEmail resends the email verification link
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not (security best practice)
		c.JSON(http.StatusOK, gin.H{"message": "if the email exists and is not verified, a verification link has been sent"})
		return
	}

	// Check if already verified
	if user.EmailVerified {
		c.JSON(http.StatusOK, gin.H{"message": "email already verified"})
		return
	}

	// Generate new verification token
	verificationToken := h.emailVerificationToken.Generate(user.Email)

	// Send verification email asynchronously
	h.emailService.SendEmailVerificationAsync(context.Background(), user.Email, verificationToken)

	// Always return success to avoid revealing if user exists
	c.JSON(http.StatusOK, gin.H{"message": "if the email exists and is not verified, a verification link has been sent"})
}
