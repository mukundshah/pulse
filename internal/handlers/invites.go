package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/middleware"
	"pulse/internal/models"
	"pulse/internal/store"
)

type InvitesHandler struct {
	store *store.Store
}

func NewInvitesHandler(s *store.Store) *InvitesHandler {
	return &InvitesHandler{store: s}
}

type CreateInviteRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// generateToken generates a random token for invitations
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *InvitesHandler) CreateInvite(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Check if user is admin
	isAdmin, err := h.store.IsProjectAdmin(projectID, userID)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only project admins can create invitations"})
		return
	}

	var req CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists and is a member
	user, err := h.store.GetUserByEmail(req.Email)
	if err == nil && user != nil {
		isMember, _ := h.store.IsProjectMember(projectID, user.ID)
		if isMember {
			c.JSON(http.StatusConflict, gin.H{"error": "user is already a member of this project"})
			return
		}
	}

	// Generate invitation token
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate invitation token"})
		return
	}

	// Create invitation
	invite := &models.ProjectInvitation{
		ID:        uuid.New(),
		ProjectID: projectID,
		Email:     req.Email,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.store.CreateProjectInvitation(invite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invitation"})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

func (h *InvitesHandler) ListInvites(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Check if user is admin
	isAdmin, err := h.store.IsProjectAdmin(projectID, userID)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only project admins can view invitations"})
		return
	}

	invites, err := h.store.GetProjectInvitations(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invitations"})
		return
	}

	c.JSON(http.StatusOK, invites)
}

type AcceptInviteRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *InvitesHandler) AcceptInvite(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get invitation by token
	invite, err := h.store.GetProjectInvitationByToken(req.Token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid or expired invitation"})
		return
	}

	// Get user to verify email matches
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.Email != invite.Email {
		c.JSON(http.StatusForbidden, gin.H{"error": "invitation email does not match your account"})
		return
	}

	// Check if already a member
	isMember, _ := h.store.IsProjectMember(invite.ProjectID, userID)
	if isMember {
		c.JSON(http.StatusConflict, gin.H{"error": "already a member of this project"})
		return
	}

	// Create project member
	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: invite.ProjectID,
		UserID:    userID,
		Role:      "member", // Default role for invited members
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.store.CreateProjectMember(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to accept invitation"})
		return
	}

	// Delete invitation
	_ = h.store.DeleteProjectInvitation(invite.ID)

	c.JSON(http.StatusOK, gin.H{"message": "invitation accepted", "project_id": invite.ProjectID})
}
