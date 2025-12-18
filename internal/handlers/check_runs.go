package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/middleware"
	"pulse/internal/models"
	"pulse/internal/store"
)

type CheckRunHandler struct {
	store *store.Store
}

func NewCheckRunHandler(s *store.Store) *CheckRunHandler {
	return &CheckRunHandler{store: s}
}

// GetCheckRun handles GET /projects/:projectId/checks/:checkId/runs/:runId
func (h *CheckRunHandler) GetCheckRun(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	isMember, err := h.store.IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	runID, err := uuid.Parse(c.Param("runId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check run ID"})
		return
	}

	// Get the check run with its associated check
	run, err := h.store.GetCheckRun(runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check run not found"})
		return
	}

	// Verify the check run belongs to the specified check
	if run.CheckID != checkID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check run not found"})
		return
	}

	// Verify the check belongs to the specified project
	// The Check is preloaded in GetCheckRun with its Project
	if run.Check == nil || run.Check.ProjectID != projectID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check run not found"})
		return
	}

	// Convert to response with computed total_time_ms
	response := toCheckRunResponse(run)
	c.JSON(http.StatusOK, response)
}

// ListCheckRuns handles GET /projects/:projectId/checks/:checkId/runs
func (h *CheckRunHandler) ListCheckRuns(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	isMember, err := h.store.IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
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

	// Parse date range parameters
	var startTime, endTime *time.Time
	if startStr := c.Query("start"); startStr != "" {
		if parsedStart, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &parsedStart
		}
	}
	if endStr := c.Query("end"); endStr != "" {
		if parsedEnd, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &parsedEnd
		}
	}

	// Fetch runs (limit+1 to check if there's a next page)
	runs, err := h.store.GetCheckRunsByCheck(checkID, limit, after, before, startTime, endTime)
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

	// Convert runs to responses with computed total_time_ms
	runResponses := make([]CheckRunResponse, len(runs))
	for i := range runs {
		runResponses[i] = toCheckRunResponse(&runs[i])
	}

	response := gin.H{
		"data":        runResponses,
		"prev_cursor": prevCursor,
		"next_cursor": nextCursor,
	}

	c.JSON(http.StatusOK, response)
}

// CheckRunResponse represents a check run with computed fields for JSON serialization.
// It embeds the CheckRun model and adds computed fields.
type CheckRunResponse struct {
	models.CheckRun
	TotalTimeMs *int `json:"total_time_ms,omitempty"` // Computed from RequestStartedAt and ResponseEndedAt
}

// toCheckRunResponse converts a CheckRun model to a CheckRunResponse with computed fields.
func toCheckRunResponse(run *models.CheckRun) CheckRunResponse {
	var totalTimeMs *int

	// Compute total response time from timestamps
	if !run.RequestStartedAt.IsZero() && !run.ResponseEndedAt.IsZero() {
		duration := run.ResponseEndedAt.Sub(run.RequestStartedAt)
		if duration > 0 {
			ms := int(duration.Milliseconds())
			totalTimeMs = &ms
		}
	}

	return CheckRunResponse{
		CheckRun:    *run,
		TotalTimeMs: totalTimeMs,
	}
}

// GetCheckUptime handles GET /projects/:projectId/checks/:checkId/uptime
// Supports both period presets and explicit datetime ranges:
//   - Query params: period (today, 1hr, 3hr, 24hr, 7d, 30d) OR start/end (RFC3339 datetime)
//   - If both are provided, start/end takes precedence
func (h *CheckRunHandler) GetCheckUptime(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	isMember, err := h.store.IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	now := time.Now().UTC()
	var startTime, endTime time.Time
	var timeBucket string
	var period string

	// Check if explicit datetime range is provided
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr != "" && endStr != "" {
		// Parse explicit datetime range
		startTime, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format. Use RFC3339 (e.g., 2024-01-01T00:00:00Z)"})
			return
		}

		endTime, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format. Use RFC3339 (e.g., 2024-01-01T23:59:59Z)"})
			return
		}

		if startTime.After(endTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
			return
		}

		// Auto-determine time bucket based on range
		timeBucket = store.DetermineTimeBucket(startTime, endTime)
		period = "custom"
	} else {
		// Use period preset
		periodParam := c.DefaultQuery("period", "24hr")
		period = periodParam

		switch periodParam {
		case "today":
			// Start of today to now
			startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			endTime = now
			timeBucket = "hour"
		case "1hr":
			startTime = now.Add(-1 * time.Hour)
			endTime = now
			timeBucket = "minute"
		case "3hr":
			startTime = now.Add(-3 * time.Hour)
			endTime = now
			timeBucket = "minute"
		case "24hr":
			startTime = now.Add(-24 * time.Hour)
			endTime = now
			timeBucket = "hour"
		case "7d":
			startTime = now.Add(-7 * 24 * time.Hour)
			endTime = now
			timeBucket = "hour"
		case "30d":
			startTime = now.Add(-30 * 24 * time.Hour)
			endTime = now
			timeBucket = "day"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period. Must be one of: today, 1hr, 3hr, 24hr, 7d, 30d, or provide start/end datetime range"})
			return
		}
	}

	// Get the actual data range to determine if we should narrow the time bucket
	actualStart, actualEnd, err := h.store.GetCheckRunsDataRange(checkID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data range"})
		return
	}

	// If we have data, check if the actual data span is much smaller than the requested range
	// If so, adjust the time bucket based on the actual data range
	if actualStart != nil && actualEnd != nil {
		requestedDuration := endTime.Sub(startTime)
		actualDuration := actualEnd.Sub(*actualStart)

		// If actual data spans less than 50% of the requested range, use a more granular bucket
		// based on the actual data range. This handles cases like requesting 7d but only having
		// data for today - we'll use hour/minute buckets instead of day buckets.
		if actualDuration > 0 && requestedDuration > 0 {
			// Calculate percentage: (actualDuration / requestedDuration) * 100
			// Use float64 to avoid integer division issues
			percentage := float64(actualDuration) / float64(requestedDuration) * 100.0
			if percentage < 50.0 {
				adjustedBucket := store.DetermineTimeBucket(*actualStart, *actualEnd)
				// Only narrow down (never widen) - e.g., if we had "day" but data is only 3 hours, use "hour"
				if shouldNarrowBucket(timeBucket, adjustedBucket) {
					timeBucket = adjustedBucket
				}
			}
		}
	}

	uptimeData, err := h.store.GetCheckUptimeData(checkID, startTime, endTime, timeBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch uptime data"})
		return
	}

	response := gin.H{
		"data":        uptimeData,
		"period":      period,
		"start_time":  startTime.Format(time.RFC3339),
		"end_time":    endTime.Format(time.RFC3339),
		"time_bucket": timeBucket,
	}

	c.JSON(http.StatusOK, response)
}

// shouldNarrowBucket determines if we should narrow the time bucket based on actual data
// Returns true if adjustedBucket is more granular than currentBucket
func shouldNarrowBucket(currentBucket, adjustedBucket string) bool {
	bucketGranularity := map[string]int{
		"second": 0,
		"minute": 1,
		"hour":   2,
		"day":    3,
		"week":   4,
	}

	currentGranularity, ok1 := bucketGranularity[currentBucket]
	adjustedGranularity, ok2 := bucketGranularity[adjustedBucket]

	if !ok1 || !ok2 {
		return false
	}

	// Narrow if adjusted is more granular (lower number)
	return adjustedGranularity < currentGranularity
}

// GetCheckTimings handles GET /projects/:projectId/checks/:checkId/timings
// Supports both period presets and explicit datetime ranges:
//   - Query params: period (today, 1hr, 3hr, 24hr, 7d, 30d) OR start/end (RFC3339 datetime)
//   - If both are provided, start/end takes precedence
//
// Returns a list of timing data points from all check runs in the time range
func (h *CheckRunHandler) GetCheckTimings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	isMember, err := h.store.IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	checkID, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	now := time.Now().UTC()
	var startTime, endTime time.Time
	var period string

	// Check if explicit datetime range is provided
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr != "" && endStr != "" {
		// Parse explicit datetime range
		startTime, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format. Use RFC3339 (e.g., 2024-01-01T00:00:00Z)"})
			return
		}

		endTime, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format. Use RFC3339 (e.g., 2024-01-01T23:59:59Z)"})
			return
		}

		if startTime.After(endTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
			return
		}

		period = "custom"
	} else {
		// Use period preset
		periodParam := c.DefaultQuery("period", "24hr")
		period = periodParam

		switch periodParam {
		case "today":
			// Start of today to now
			startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			endTime = now
		case "1hr":
			startTime = now.Add(-1 * time.Hour)
			endTime = now
		case "3hr":
			startTime = now.Add(-3 * time.Hour)
			endTime = now
		case "24hr":
			startTime = now.Add(-24 * time.Hour)
			endTime = now
		case "7d":
			startTime = now.Add(-7 * 24 * time.Hour)
			endTime = now
		case "30d":
			startTime = now.Add(-30 * 24 * time.Hour)
			endTime = now
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period. Must be one of: today, 1hr, 3hr, 24hr, 7d, 30d, or provide start/end datetime range"})
			return
		}
	}

	timingsData, err := h.store.GetCheckTimingsData(checkID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timings data"})
		return
	}

	response := gin.H{
		"data":       timingsData,
		"period":     period,
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
