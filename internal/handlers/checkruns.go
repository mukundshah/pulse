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

	// Parse cursor parameters
	var after, before *uuid.UUID
	if afterStr := c.Query("after"); afterStr != "" {
		if parsedAfter, err := uuid.Parse(afterStr); err == nil {
			after = &parsedAfter
		}
	}
	if beforeStr := c.Query("before"); beforeStr != "" {
		if parsedBefore, err := uuid.Parse(beforeStr); err == nil {
			before = &parsedBefore
		}
	}

	// Fetch runs (limit+1 to check if there's a next page)
	runs, err := h.store.GetCheckRunsByCheck(checkID, limit, after, before)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list check runs"})
		return
	}

	// Determine pagination cursors
	var prevCursor, nextCursor *string

	// Check if we have more items (fetched limit+1)
	hasNext := len(runs) > limit
	if hasNext {
		// Remove the extra item
		runs = runs[:limit]
	}

	// Set next cursor (oldest item in current page) - for getting older items
	// Set when there are more older items available
	if len(runs) > 0 && hasNext {
		nextCursorStr := runs[len(runs)-1].ID.String()
		nextCursor = &nextCursorStr
	}

	// Set prev cursor (newest item in current page) - for getting newer items
	// Set when we're not at the beginning (we used an "after" or "before" cursor)
	if (after != nil || before != nil) && len(runs) > 0 {
		prevCursorStr := runs[0].ID.String()
		prevCursor = &prevCursorStr
	}

	response := gin.H{
		"data":        runs,
		"prev_cursor": prevCursor,
		"next_cursor": nextCursor,
	}

	c.JSON(http.StatusOK, response)
}
