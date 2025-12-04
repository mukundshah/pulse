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

// CreateTag handles POST /tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	if err := h.store.CreateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// ListTags handles GET /tags
func (h *TagHandler) ListTags(c *gin.Context) {
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

	project, err := h.store.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load project"})
		return
	}

	c.JSON(http.StatusOK, project)
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

	project, err := h.store.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load project"})
		return
	}

	c.JSON(http.StatusOK, project)
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
