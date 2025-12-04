package handlers

import (
	"net/http"

	"pulse/internal/models"
	"pulse/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagHandler struct {
	store *store.Store
}

func NewTagHandler(s *store.Store) *TagHandler {
	return &TagHandler{store: s}
}

// CreateTag handles POST /projects/:projectId/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Verify project exists
	_, err = h.store.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name:      req.Name,
		ProjectID: projectID,
	}

	if err := h.store.CreateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// ListTags handles GET /tags or GET /projects/:projectId/tags
func (h *TagHandler) ListTags(c *gin.Context) {
	// Check if projectId is provided in the path
	if projectIDStr := c.Param("projectId"); projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		tags, err := h.store.GetTagsByProject(projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tags"})
			return
		}

		c.JSON(http.StatusOK, tags)
		return
	}

	// List all tags if no project ID provided
	tags, err := h.store.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// AddTagToProject handles POST /projects/:projectId/tags/:tagId
func (h *TagHandler) AddTagToProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.store.AddTagToProject(projectID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to project"})
		return
	}

	// Reload tag to return updated tag
	tag, err := h.store.GetTag(tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tag"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// RemoveTagFromProject handles DELETE /projects/:projectId/tags/:tagId
func (h *TagHandler) RemoveTagFromProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.store.RemoveTagFromProject(projectID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed from project"})
}

// AddTagToCheck handles POST /checks/:checkId/tags/:tagId
func (h *TagHandler) AddTagToCheck(c *gin.Context) {
	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.store.AddTagToCheck(checkID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to check"})
		return
	}

	check, err := h.store.GetCheck(checkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load check"})
		return
	}

	c.JSON(http.StatusOK, check)
}

// RemoveTagFromCheck handles DELETE /checks/:checkId/tags/:tagId
func (h *TagHandler) RemoveTagFromCheck(c *gin.Context) {
	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := h.store.RemoveTagFromCheck(checkID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from check"})
		return
	}

	check, err := h.store.GetCheck(checkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load check"})
		return
	}

	c.JSON(http.StatusOK, check)
}
