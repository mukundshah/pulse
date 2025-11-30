package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse/internal/models"
	"pulse/internal/store"
)

type ChecksHandler struct {
	store     *store.Store
	runsStore *store.RunsStore
}

func NewChecksHandler(s *store.Store, rs *store.RunsStore) *ChecksHandler {
	return &ChecksHandler{store: s, runsStore: rs}
}

func (h *ChecksHandler) ListChecks(c *gin.Context) {
	checks, err := h.store.GetAllChecks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, checks)
}

func (h *ChecksHandler) GetCheck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid check ID"})
		return
	}

	check, err := h.store.GetCheck(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, check)
}

func (h *ChecksHandler) CreateCheck(c *gin.Context) {
	var check models.Check
	if err := c.ShouldBindJSON(&check); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if check.ID == uuid.Nil {
		check.ID = uuid.New()
	}
	if check.Method == "" {
		check.Method = "GET"
	}
	if check.ExpectedStatus == 0 {
		check.ExpectedStatus = 200
	}
	if check.TimeoutMs == 0 {
		check.TimeoutMs = 10000
	}
	if check.AlertThreshold == 0 {
		check.AlertThreshold = 3
	}
	if check.IntervalSeconds == 0 {
		check.IntervalSeconds = 60
	}
	if check.LastStatus == "" {
		check.LastStatus = "unknown"
	}
	now := time.Now()
	check.CreatedAt = now
	check.UpdatedAt = now
	if check.NextRunAt == nil {
		nextRun := now.Add(time.Duration(check.IntervalSeconds) * time.Second)
		check.NextRunAt = &nextRun
	}

	if err := h.store.CreateCheck(&check); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, check)
}

func (h *ChecksHandler) GetCheckRuns(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid check ID"})
		return
	}

	// Verify check exists
	_, err = h.store.GetCheck(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Check not found"})
		return
	}

	// Check if ClickHouse is available
	if h.runsStore == nil {
		c.JSON(503, gin.H{"error": "ClickHouse is not available"})
		return
	}

	// Get limit from query parameter (default: 100)
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 1000 {
			c.JSON(400, gin.H{"error": "Invalid limit. Must be between 1 and 1000"})
			return
		}
	}

	// Get runs from ClickHouse
	ctx := context.Background()
	runs, err := h.runsStore.GetCheckRuns(ctx, id, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Create response without check_id in each run
	type RunResponse struct {
		ID         uuid.UUID `json:"id"`
		Status     string    `json:"status"`
		LatencyMs  int64     `json:"latency_ms"`
		StatusCode int32     `json:"status_code"`
		Error      *string   `json:"error,omitempty"`
		RunAt      time.Time `json:"run_at"`
	}

	runResponses := make([]RunResponse, len(runs))
	for i, run := range runs {
		runResponses[i] = RunResponse{
			ID:         run.ID,
			Status:     run.Status,
			LatencyMs:  run.LatencyMs,
			StatusCode: run.StatusCode,
			Error:      run.Error,
			RunAt:      run.RunAt,
		}
	}

	c.JSON(200, gin.H{
		"runs":  runResponses,
		"count": len(runResponses),
	})
}

func (h *ChecksHandler) GetCheckAlerts(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid check ID"})
		return
	}

	// Verify check exists
	_, err = h.store.GetCheck(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Check not found"})
		return
	}

	// Get limit from query parameter (default: 100)
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 1000 {
			c.JSON(400, gin.H{"error": "Invalid limit. Must be between 1 and 1000"})
			return
		}
	}

	// Get alerts from database
	alerts, err := h.store.GetAlerts(id, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

func (h *ChecksHandler) GetCheckWebhooks(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid check ID"})
		return
	}

	// Verify check exists
	_, err = h.store.GetCheck(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Check not found"})
		return
	}

	// Get limit from query parameter (default: 100)
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 1000 {
			c.JSON(400, gin.H{"error": "Invalid limit. Must be between 1 and 1000"})
			return
		}
	}

	// Get webhook attempts from database
	attempts, err := h.store.GetWebhookAttempts(id, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"webhooks": attempts,
		"count":    len(attempts),
	})
}
