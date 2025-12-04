package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/store"
)

type CheckRunHandler struct {
	store *store.Store
}

func NewCheckRunHandler(s *store.Store) *CheckRunHandler {
	return &CheckRunHandler{store: s}
}

// GetCheckRun handles GET /check-runs/:id
func (h *CheckRunHandler) GetCheckRun(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check run ID"})
		return
	}

	run, err := h.store.GetCheckRun(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check run not found"})
		return
	}

	c.JSON(http.StatusOK, run)
}

// ListCheckRuns handles GET /checks/:checkId/runs
func (h *CheckRunHandler) ListCheckRuns(c *gin.Context) {
	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	// Verify check exists
	_, err = h.store.GetCheck(checkID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found"})
		return
	}

	// Get limit from query parameter (default: 50)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	runs, err := h.store.GetCheckRunsByCheck(checkID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list check runs"})
		return
	}

	c.JSON(http.StatusOK, runs)
}
