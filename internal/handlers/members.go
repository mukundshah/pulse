package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/middleware"
	"pulse/internal/store"
)

type MembersHandler struct {
	store *store.Store
}

func NewMembersHandler(s *store.Store) *MembersHandler {
	return &MembersHandler{store: s}
}

func (h *MembersHandler) ListMembers(c *gin.Context) {
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

	// Check if user is a member
	isMember, err := h.store.IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	members, err := h.store.GetProjectMembers(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member viewer"`
}

func (h *MembersHandler) UpdateMemberRole(c *gin.Context) {
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

	memberUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Check if requester is admin
	isAdmin, err := h.store.IsProjectAdmin(projectID, userID)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only project admins can update member roles"})
		return
	}

	// Prevent removing the last admin
	if memberUserID == userID {
		members, _ := h.store.GetProjectMembers(projectID)
		adminCount := 0
		for _, m := range members {
			if m.Role == "admin" {
				adminCount++
			}
		}
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change role of the last admin"})
			return
		}
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.UpdateProjectMemberRole(projectID, memberUserID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update member role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member role updated"})
}

func (h *MembersHandler) RemoveMember(c *gin.Context) {
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

	memberUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Check if requester is admin
	isAdmin, err := h.store.IsProjectAdmin(projectID, userID)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only project admins can remove members"})
		return
	}

	// Prevent removing the last admin
	if memberUserID == userID {
		members, _ := h.store.GetProjectMembers(projectID)
		adminCount := 0
		for _, m := range members {
			if m.Role == "admin" {
				adminCount++
			}
		}
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot remove the last admin"})
			return
		}
	}

	if err := h.store.RemoveProjectMember(projectID, memberUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}
