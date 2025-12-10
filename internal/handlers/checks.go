package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"pulse/internal/models"
	"pulse/internal/store"
)

type CheckHandler struct {
	store *store.Store
}

func NewCheckHandler(s *store.Store) *CheckHandler {
	return &CheckHandler{store: s}
}

// CreateCheck handles POST /projects/:projectId/checks
func (h *CheckHandler) CreateCheck(c *gin.Context) {
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
		Name             string         `json:"name" binding:"required"`
		Type             string         `json:"type" binding:"required"`
		Host             string         `json:"host" binding:"required"`
		Port             *int           `json:"port"`
		Secure           *bool          `json:"secure"`
		Method           string         `json:"method"`
		Headers          datatypes.JSON `json:"headers"`
		QueryParams      datatypes.JSON `json:"query_params"`
		Body             datatypes.JSON `json:"body"`
		IPVersion        string         `json:"ip_version"`
		SSLVerification  *bool          `json:"ssl_verification"`
		FollowRedirects  *bool          `json:"follow_redirects"`
		PlaywrightScript *string        `json:"playwright_script,omitempty"`
		Assertions       datatypes.JSON `json:"assertions"`
		ExpectedStatus   int            `json:"expected_status"`
		ShouldFail       bool           `json:"should_fail"`
		PreScript        *string        `json:"pre_script,omitempty"`
		PostScript       *string        `json:"post_script,omitempty"`
		TimeoutMs        int            `json:"timeout_ms"`
		IntervalSeconds  int            `json:"interval_seconds" binding:"required"`
		AlertThreshold   int            `json:"alert_threshold"`
		IsEnabled        bool           `json:"is_enabled"`
		IsMuted          bool           `json:"is_muted"`
		TagIDs           []uuid.UUID    `json:"tag_ids,omitempty"`
		RegionIDs        []uuid.UUID    `json:"region_ids,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkType := models.CheckType(req.Type)
	if checkType != models.CheckTypeHTTP && checkType != models.CheckTypeTCP &&
		checkType != models.CheckTypeDNS && checkType != models.CheckTypeBrowser &&
		checkType != models.CheckTypeHeartbeat {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check type"})
		return
	}

	check := &models.Check{
		Name:             req.Name,
		Type:             checkType,
		Host:             req.Host,
		Method:           req.Method,
		Headers:          req.Headers,
		QueryParams:      req.QueryParams,
		Body:             req.Body,
		IPVersion:        models.IPVersionType(req.IPVersion),
		PlaywrightScript: req.PlaywrightScript,
		Assertions:       req.Assertions,
		ExpectedStatus:   req.ExpectedStatus,
		ShouldFail:       req.ShouldFail,
		PreScript:        req.PreScript,
		PostScript:       req.PostScript,
		TimeoutMs:        req.TimeoutMs,
		IntervalSeconds:  req.IntervalSeconds,
		AlertThreshold:   req.AlertThreshold,
		IsEnabled:        req.IsEnabled,
		IsMuted:          req.IsMuted,
		ProjectID:        projectID,
	}

	// Set defaults
	if check.Method == "" {
		check.Method = "GET"
	}

	// Handle Port and Secure
	if req.Port != nil {
		check.Port = *req.Port
	} else {
		// Default port based on secure flag
		if req.Secure != nil && *req.Secure {
			check.Port = 443
		} else {
			check.Port = 80
		}
	}

	if req.Secure != nil {
		check.Secure = *req.Secure
	} else {
		// Default secure based on port
		check.Secure = check.Port == 443
	}

	if req.SSLVerification != nil {
		check.SSLVerification = *req.SSLVerification
	} else {
		check.SSLVerification = true // default
	}

	if req.FollowRedirects != nil {
		check.FollowRedirects = *req.FollowRedirects
	} else {
		check.FollowRedirects = true // default
	}

	if check.IPVersion == "" {
		check.IPVersion = models.IPVersionTypeIPv4
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
	if !req.IsEnabled {
		check.IsEnabled = true
	}

	if err := h.store.CreateCheck(check); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create check"})
		return
	}

	// Add tags if provided
	if len(req.TagIDs) > 0 {
		for _, tagID := range req.TagIDs {
			_ = h.store.AddTagToCheck(check.ID, tagID)
		}
	}

	// Add regions if provided
	if len(req.RegionIDs) > 0 {
		// Note: Region association would need to be implemented in store
		// For now, we'll just create the check
	}

	// Reload check with associations
	check, err = h.store.GetCheck(check.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load check"})
		return
	}

	c.JSON(http.StatusCreated, check)
}

// GetCheck handles GET /checks/:checkId
func (h *CheckHandler) GetCheck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	check, err := h.store.GetCheck(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found"})
		return
	}

	c.JSON(http.StatusOK, check)
}

// ListChecks handles GET /projects/:projectId/checks
func (h *CheckHandler) ListChecks(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	checks, err := h.store.GetChecksByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list checks"})
		return
	}

	c.JSON(http.StatusOK, checks)
}

// UpdateCheck handles PUT /checks/:checkId
func (h *CheckHandler) UpdateCheck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	check, err := h.store.GetCheck(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found"})
		return
	}

	var req struct {
		Name             string         `json:"name"`
		Type             string         `json:"type"`
		Host             string         `json:"host"`
		Port             *int           `json:"port"`
		Secure           *bool          `json:"secure"`
		Method           string         `json:"method"`
		Headers          datatypes.JSON `json:"headers"`
		QueryParams      datatypes.JSON `json:"query_params"`
		Body             datatypes.JSON `json:"body"`
		IPVersion        string         `json:"ip_version"`
		SSLVerification  *bool          `json:"ssl_verification"`
		FollowRedirects  *bool          `json:"follow_redirects"`
		PlaywrightScript *string        `json:"playwright_script,omitempty"`
		Assertions       datatypes.JSON `json:"assertions"`
		ExpectedStatus   int            `json:"expected_status"`
		ShouldFail       bool           `json:"should_fail"`
		PreScript        *string        `json:"pre_script,omitempty"`
		PostScript       *string        `json:"post_script,omitempty"`
		TimeoutMs        int            `json:"timeout_ms"`
		IntervalSeconds  int            `json:"interval_seconds"`
		AlertThreshold   int            `json:"alert_threshold"`
		IsEnabled        *bool          `json:"is_enabled"`
		IsMuted          *bool          `json:"is_muted"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		check.Name = req.Name
	}
	if req.Type != "" {
		checkType := models.CheckType(req.Type)
		check.Type = checkType
	}
	if req.Host != "" {
		check.Host = req.Host
	}
	if req.Port != nil {
		check.Port = *req.Port
		// Auto-set secure based on port if secure not explicitly provided
		if req.Secure == nil {
			check.Secure = *req.Port == 443
		} else {
			check.Secure = *req.Secure
		}
	}
	if req.Secure != nil && req.Port == nil {
		check.Secure = *req.Secure
	}
	if req.Method != "" {
		check.Method = req.Method
	}
	if req.Headers != nil {
		check.Headers = req.Headers
	}
	if req.QueryParams != nil {
		check.QueryParams = req.QueryParams
	}
	if req.Body != nil {
		check.Body = req.Body
	}
	if req.IPVersion != "" {
		check.IPVersion = models.IPVersionType(req.IPVersion)
	}
	if req.SSLVerification != nil {
		check.SSLVerification = *req.SSLVerification
	}
	if req.FollowRedirects != nil {
		check.FollowRedirects = *req.FollowRedirects
	}
	if req.PlaywrightScript != nil {
		check.PlaywrightScript = req.PlaywrightScript
	}
	if req.Assertions != nil {
		check.Assertions = req.Assertions
	}
	if req.ExpectedStatus != 0 {
		check.ExpectedStatus = req.ExpectedStatus
	}
	check.ShouldFail = req.ShouldFail
	if req.PreScript != nil {
		check.PreScript = req.PreScript
	}
	if req.PostScript != nil {
		check.PostScript = req.PostScript
	}
	if req.TimeoutMs != 0 {
		check.TimeoutMs = req.TimeoutMs
	}
	if req.IntervalSeconds != 0 {
		check.IntervalSeconds = req.IntervalSeconds
	}
	if req.AlertThreshold != 0 {
		check.AlertThreshold = req.AlertThreshold
	}
	if req.IsEnabled != nil {
		check.IsEnabled = *req.IsEnabled
	}
	if req.IsMuted != nil {
		check.IsMuted = *req.IsMuted
	}

	if err := h.store.UpdateCheck(check); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update check"})
		return
	}

	// Reload check with associations
	check, err = h.store.GetCheck(check.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load check"})
		return
	}

	c.JSON(http.StatusOK, check)
}

// DeleteCheck handles DELETE /checks/:checkId
func (h *CheckHandler) DeleteCheck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("checkId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check ID"})
		return
	}

	if err := h.store.DeleteCheck(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete check"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check deleted successfully"})
}
