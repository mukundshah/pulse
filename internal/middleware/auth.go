package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/auth/token"
	"pulse/internal/config"
	"pulse/internal/store"
)

// AuthMiddleware validates JWT tokens, checks session validity, and sets user context
func AuthMiddleware(cfg *config.Config, s *store.Store) gin.HandlerFunc {
	jwtGenerator := token.NewJWTTokenGenerator(token.TokenConfig{
		Secret:   cfg.JWTSecret,
		Validity: 24 * time.Hour,
	})

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtGenerator.Validate(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Validate session using JTI
		if claims.ID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing session identifier"})
			c.Abort()
			return
		}

		// Check if session is valid
		session, err := s.GetSessionByJTI(claims.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			c.Abort()
			return
		}

		// Update session activity (non-blocking, don't fail if this errors)
		go func() {
			_ = s.UpdateSessionActivity(claims.ID)
		}()

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("jti", claims.ID) // Store JTI for logout functionality
		c.Set("session_id", session.ID)

		c.Next()
	}
}

// GetUserID extracts user ID from context (must be called after AuthMiddleware)
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}
