package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"pulse/internal/middleware"
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

	var req struct {
		Name                  string         `json:"name" binding:"required"`
		Type                  string         `json:"type" binding:"required"`
		Host                  string         `json:"host" binding:"required"`
		Port                  *int           `json:"port"`
		Secure                *bool          `json:"secure"`
		Method                string         `json:"method"`
		Path                  string         `json:"path"`
		QueryParams           datatypes.JSON `json:"query_params"`
		Headers               datatypes.JSON `json:"headers"`
		Body                  datatypes.JSON `json:"body"`
		IPVersion             string         `json:"ip_version"`
		SkipSSLVerification   *bool          `json:"skip_ssl_verification"`
		FollowRedirects       *bool          `json:"follow_redirects"`
		PlaywrightScript      *string        `json:"playwright_script,omitempty"`
		Assertions            datatypes.JSON `json:"assertions"`
		PreScript             *string        `json:"pre_script,omitempty"`
		PostScript            *string        `json:"post_script,omitempty"`
		Interval              string         `json:"interval" binding:"required"`
		DegradedThreshold     int            `json:"degraded_threshold"`
		DegradedThresholdUnit string         `json:"degraded_threshold_unit"`
		FailedThreshold       int            `json:"failed_threshold"`
		FailedThresholdUnit   string         `json:"failed_threshold_unit"`
		Retries               string         `json:"retries"`
		RetriesCount          *int           `json:"retries_count,omitempty"`
		RetriesDelay          *int           `json:"retries_delay,omitempty"`
		RetriesDelayUnit      *string        `json:"retries_delay_unit,omitempty"`
		RetriesFactor         *float64       `json:"retries_factor,omitempty"`
		RetriesJitter         *string        `json:"retries_jitter,omitempty"`
		RetriesJitterFactor   *float64       `json:"retries_jitter_factor,omitempty"`
		RetriesMaxDelay       *int           `json:"retries_max_delay,omitempty"`
		RetriesMaxDelayUnit   *string        `json:"retries_max_delay_unit,omitempty"`
		RetriesTimeout        *int           `json:"retries_timeout,omitempty"`
		RetriesTimeoutUnit    *string        `json:"retries_timeout_unit,omitempty"`
		IsEnabled             bool           `json:"is_enabled"`
		IsMuted               bool           `json:"is_muted"`
		ShouldFail            bool           `json:"should_fail"`
		TagIDs                []uuid.UUID    `json:"tag_ids,omitempty"`
		RegionIDs             []uuid.UUID    `json:"region_ids,omitempty"`
		DNSRecordType         *string        `json:"dns_record_type,omitempty"`
		DNSResolver           *string        `json:"dns_resolver,omitempty"`
		DNSResolverPort       *int           `json:"dns_resolver_port,omitempty"`
		DNSResolverProtocol   *string        `json:"dns_resolver_protocol,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that at least one region is required
	if len(req.RegionIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one region is required"})
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
		Name:                  req.Name,
		Type:                  checkType,
		Host:                  req.Host,
		Method:                req.Method,
		Path:                  req.Path,
		QueryParams:           req.QueryParams,
		Headers:               req.Headers,
		Body:                  req.Body,
		IPVersion:             models.IPVersionType(req.IPVersion),
		PlaywrightScript:      req.PlaywrightScript,
		Assertions:            req.Assertions,
		PreScript:             req.PreScript,
		PostScript:            req.PostScript,
		Interval:              req.Interval,
		DegradedThreshold:     req.DegradedThreshold,
		DegradedThresholdUnit: models.UnitType(req.DegradedThresholdUnit),
		FailedThreshold:       req.FailedThreshold,
		FailedThresholdUnit:   models.UnitType(req.FailedThresholdUnit),
		Retries:               models.RetryType(req.Retries),
		RetriesCount:          req.RetriesCount,
		RetriesDelay:          req.RetriesDelay,
		RetriesDelayUnit:      (*models.UnitType)(req.RetriesDelayUnit),
		RetriesFactor:         req.RetriesFactor,
		RetriesJitter:         (*models.RetryJitterType)(req.RetriesJitter),
		RetriesJitterFactor:   req.RetriesJitterFactor,
		RetriesMaxDelay:       req.RetriesMaxDelay,
		RetriesMaxDelayUnit:   (*models.UnitType)(req.RetriesMaxDelayUnit),
		RetriesTimeout:        req.RetriesTimeout,
		RetriesTimeoutUnit:    (*models.UnitType)(req.RetriesTimeoutUnit),
		IsEnabled:             req.IsEnabled,
		IsMuted:               req.IsMuted,
		ShouldFail:            req.ShouldFail,
		ProjectID:             projectID,
	}

	// Handle DNS fields
	if req.DNSRecordType != nil {
		check.DNSRecordType = (*models.DNSRecordType)(req.DNSRecordType)
	}
	if req.DNSResolver != nil {
		check.DNSResolver = req.DNSResolver
	}
	if req.DNSResolverPort != nil {
		check.DNSResolverPort = req.DNSResolverPort
	}
	if req.DNSResolverProtocol != nil {
		check.DNSResolverProtocol = (*models.DNSResolverProtocolType)(req.DNSResolverProtocol)
	}

	// Set defaults
	if check.Method == "" {
		check.Method = "GET"
	}
	if check.Path == "" {
		check.Path = "/"
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

	if req.SkipSSLVerification != nil {
		check.SkipSSLVerification = *req.SkipSSLVerification
	} else {
		check.SkipSSLVerification = false // default
	}

	if req.FollowRedirects != nil {
		check.FollowRedirects = *req.FollowRedirects
	} else {
		check.FollowRedirects = true // default
	}

	if check.IPVersion == "" {
		check.IPVersion = models.IPVersionTypeIPv4
	}
	if check.Interval == "" {
		check.Interval = "10m"
	}
	if check.DegradedThreshold == 0 {
		check.DegradedThreshold = 3000
	}
	if check.DegradedThresholdUnit == "" {
		check.DegradedThresholdUnit = models.UnitTypeMs
	}
	if check.FailedThreshold == 0 {
		check.FailedThreshold = 5000
	}
	if check.FailedThresholdUnit == "" {
		check.FailedThresholdUnit = models.UnitTypeMs
	}
	if check.Retries == "" {
		check.Retries = models.RetryTypeNone
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

	// Add regions
	for _, regionID := range req.RegionIDs {
		if err := h.store.AddRegionToCheck(check.ID, regionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate region with check"})
			return
		}
	}

	// Reload check with associations
	check, err = h.store.GetCheck(check.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load check"})
		return
	}

	c.JSON(http.StatusCreated, check)
}

// GetCheck handles GET /projects/:projectId/checks/:checkId
func (h *CheckHandler) GetCheck(c *gin.Context) {
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

	check, err := h.store.GetCheck(checkID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found"})
		return
	}

	c.JSON(http.StatusOK, check)
}

// ListChecks handles GET /projects/:projectId/checks
func (h *CheckHandler) ListChecks(c *gin.Context) {
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

	checks, err := h.store.GetChecksByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list checks"})
		return
	}

	c.JSON(http.StatusOK, checks)
}

// UpdateCheck handles PUT /projects/:projectId/checks/:checkId
func (h *CheckHandler) UpdateCheck(c *gin.Context) {
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

	check, err := h.store.GetCheck(checkID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check not found"})
		return
	}

	var req struct {
		Name                  string         `json:"name"`
		Type                  string         `json:"type"`
		Host                  string         `json:"host"`
		Port                  *int           `json:"port"`
		Secure                *bool          `json:"secure"`
		Method                string         `json:"method"`
		Path                  string         `json:"path"`
		QueryParams           datatypes.JSON `json:"query_params"`
		Headers               datatypes.JSON `json:"headers"`
		Body                  datatypes.JSON `json:"body"`
		IPVersion             string         `json:"ip_version"`
		SkipSSLVerification   *bool          `json:"skip_ssl_verification"`
		FollowRedirects       *bool          `json:"follow_redirects"`
		PlaywrightScript      *string        `json:"playwright_script,omitempty"`
		Assertions            datatypes.JSON `json:"assertions"`
		PreScript             *string        `json:"pre_script,omitempty"`
		PostScript            *string        `json:"post_script,omitempty"`
		Interval              string         `json:"interval"`
		DegradedThreshold     *int           `json:"degraded_threshold"`
		DegradedThresholdUnit *string        `json:"degraded_threshold_unit"`
		FailedThreshold       *int           `json:"failed_threshold"`
		FailedThresholdUnit   *string        `json:"failed_threshold_unit"`
		Retries               *string        `json:"retries"`
		RetriesCount          *int           `json:"retries_count,omitempty"`
		RetriesDelay          *int           `json:"retries_delay,omitempty"`
		RetriesDelayUnit      *string        `json:"retries_delay_unit,omitempty"`
		RetriesFactor         *float64       `json:"retries_factor,omitempty"`
		RetriesJitter         *string        `json:"retries_jitter,omitempty"`
		RetriesJitterFactor   *float64       `json:"retries_jitter_factor,omitempty"`
		RetriesMaxDelay       *int           `json:"retries_max_delay,omitempty"`
		RetriesMaxDelayUnit   *string        `json:"retries_max_delay_unit,omitempty"`
		RetriesTimeout        *int           `json:"retries_timeout,omitempty"`
		RetriesTimeoutUnit    *string        `json:"retries_timeout_unit,omitempty"`
		IsEnabled             *bool          `json:"is_enabled"`
		IsMuted               *bool          `json:"is_muted"`
		ShouldFail            *bool          `json:"should_fail"`
		DNSRecordType         *string        `json:"dns_record_type,omitempty"`
		DNSResolver           *string        `json:"dns_resolver,omitempty"`
		DNSResolverPort       *int           `json:"dns_resolver_port,omitempty"`
		DNSResolverProtocol   *string        `json:"dns_resolver_protocol,omitempty"`
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
	if req.Path != "" {
		check.Path = req.Path
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
	if req.SkipSSLVerification != nil {
		check.SkipSSLVerification = *req.SkipSSLVerification
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
	if req.PreScript != nil {
		check.PreScript = req.PreScript
	}
	if req.PostScript != nil {
		check.PostScript = req.PostScript
	}
	if req.Interval != "" {
		check.Interval = req.Interval
	}
	if req.DegradedThreshold != nil {
		check.DegradedThreshold = *req.DegradedThreshold
	}
	if req.DegradedThresholdUnit != nil {
		check.DegradedThresholdUnit = models.UnitType(*req.DegradedThresholdUnit)
	}
	if req.FailedThreshold != nil {
		check.FailedThreshold = *req.FailedThreshold
	}
	if req.FailedThresholdUnit != nil {
		check.FailedThresholdUnit = models.UnitType(*req.FailedThresholdUnit)
	}
	if req.Retries != nil {
		check.Retries = models.RetryType(*req.Retries)
	}
	if req.RetriesCount != nil {
		check.RetriesCount = req.RetriesCount
	}
	if req.RetriesDelay != nil {
		check.RetriesDelay = req.RetriesDelay
	}
	if req.RetriesDelayUnit != nil {
		check.RetriesDelayUnit = (*models.UnitType)(req.RetriesDelayUnit)
	}
	if req.RetriesFactor != nil {
		check.RetriesFactor = req.RetriesFactor
	}
	if req.RetriesJitter != nil {
		check.RetriesJitter = (*models.RetryJitterType)(req.RetriesJitter)
	}
	if req.RetriesJitterFactor != nil {
		check.RetriesJitterFactor = req.RetriesJitterFactor
	}
	if req.RetriesMaxDelay != nil {
		check.RetriesMaxDelay = req.RetriesMaxDelay
	}
	if req.RetriesMaxDelayUnit != nil {
		check.RetriesMaxDelayUnit = (*models.UnitType)(req.RetriesMaxDelayUnit)
	}
	if req.RetriesTimeout != nil {
		check.RetriesTimeout = req.RetriesTimeout
	}
	if req.RetriesTimeoutUnit != nil {
		check.RetriesTimeoutUnit = (*models.UnitType)(req.RetriesTimeoutUnit)
	}
	if req.IsEnabled != nil {
		check.IsEnabled = *req.IsEnabled
	}
	if req.IsMuted != nil {
		check.IsMuted = *req.IsMuted
	}
	if req.ShouldFail != nil {
		check.ShouldFail = *req.ShouldFail
	}
	if req.DNSRecordType != nil {
		check.DNSRecordType = (*models.DNSRecordType)(req.DNSRecordType)
	}
	if req.DNSResolver != nil {
		check.DNSResolver = req.DNSResolver
	}
	if req.DNSResolverPort != nil {
		check.DNSResolverPort = req.DNSResolverPort
	}
	if req.DNSResolverProtocol != nil {
		check.DNSResolverProtocol = (*models.DNSResolverProtocolType)(req.DNSResolverProtocol)
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

// DeleteCheck handles DELETE /projects/:projectId/checks/:checkId
func (h *CheckHandler) DeleteCheck(c *gin.Context) {
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

	if err := h.store.DeleteCheck(checkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete check"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check deleted successfully"})
}
