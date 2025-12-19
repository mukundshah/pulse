package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

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

// CheckListResponse represents a check in the list with metrics
type CheckListResponse struct {
	ID                   uuid.UUID    `json:"id"`
	Type                 string       `json:"type"`
	Name                 string       `json:"name"`
	Interval             string       `json:"interval"`
	LastRun              *time.Time   `json:"last_run,omitempty"`
	LastStatus           string       `json:"last_status"`
	Last24Runs           []RunSummary `json:"last_24_runs"`
	Uptime24h            *string      `json:"uptime_24h,omitempty"`
	Uptime7d             *string      `json:"uptime_7d,omitempty"`
	AvgResponseTime24hMs *string      `json:"avg_response_time_24h_ms,omitempty"`
	P95ResponseTime24hMs *string      `json:"p95_response_time_24h_ms,omitempty"`
}

// RunSummary represents a summary of a check run
type RunSummary struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	TotalTimeMs *int       `json:"total_time_ms,omitempty"`
	Status      *string    `json:"status,omitempty"`
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

	now := time.Now().UTC()
	twentyFourHoursAgo := now.Add(-24 * time.Hour)
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)

	// Single optimized query using CTEs to get all data
	type CheckListRow struct {
		ID                 uuid.UUID  `gorm:"column:id"`
		Type               string     `gorm:"column:type"`
		Name               string     `gorm:"column:name"`
		Interval           string     `gorm:"column:interval"`
		LastRun            *time.Time `gorm:"column:last_run_at"`
		LastStatus         string     `gorm:"column:last_status"`
		Last24Runs         string     `gorm:"column:last_24_runs"` // JSON array
		Uptime24h          *float64   `gorm:"column:uptime_24h"`
		Uptime7d           *float64   `gorm:"column:uptime_7d"`
		AvgResponseTime24h *float64   `gorm:"column:avg_response_time_24h"`
		P95ResponseTime24h *float64   `gorm:"column:p95_response_time_24h"`
	}

	var rows []CheckListRow
	err = h.store.DB().Raw(`
		WITH checks_base AS (
			SELECT
				id,
				type,
				name,
				interval,
				last_run_at,
				last_status
			FROM checks
			WHERE project_id = ? AND deleted_at IS NULL
		),
		ranked_runs AS (
			SELECT
				check_id,
				status,
				request_started_at,
				response_ended_at,
				created_at,
				id,
				ROW_NUMBER() OVER (PARTITION BY check_id ORDER BY created_at DESC, id DESC) as row_num
			FROM check_runs
			WHERE deleted_at IS NULL
		),
		last_24_runs AS (
			SELECT
				check_id,
				json_agg(
					json_build_object(
						'id', id::text,
						'timestamp', created_at,
						'total_time_ms',
						CASE
							WHEN request_started_at IS NOT NULL AND response_ended_at IS NOT NULL
								AND response_ended_at > request_started_at
							THEN EXTRACT(EPOCH FROM (response_ended_at - request_started_at)) * 1000
							ELSE 0
						END,
						'status', status
					) ORDER BY created_at ASC, id ASC
				) as runs
			FROM ranked_runs
			WHERE row_num <= 24
			GROUP BY check_id
		),
		uptime_stats AS (
			SELECT
				check_id,
				-- 24h uptime
				CASE
					WHEN COUNT(*) FILTER (WHERE created_at >= ? AND created_at <= ?) > 0
					THEN (
						COUNT(*) FILTER (
							WHERE created_at >= ? AND created_at <= ?
							AND status IN ('passing', 'degraded')
						)::float /
						NULLIF(COUNT(*) FILTER (WHERE created_at >= ? AND created_at <= ?), 0) * 100.0
					)
					ELSE NULL
				END as uptime_24h,
				-- 7d uptime
				CASE
					WHEN COUNT(*) FILTER (WHERE created_at >= ? AND created_at <= ?) > 0
					THEN (
						COUNT(*) FILTER (
							WHERE created_at >= ? AND created_at <= ?
							AND status IN ('passing', 'degraded')
						)::float /
						NULLIF(COUNT(*) FILTER (WHERE created_at >= ? AND created_at <= ?), 0) * 100.0
					)
					ELSE NULL
				END as uptime_7d
			FROM check_runs
			WHERE deleted_at IS NULL
			GROUP BY check_id
		),
		response_time_stats AS (
			SELECT
				check_id,
				AVG(EXTRACT(EPOCH FROM (response_ended_at - request_started_at)) * 1000) as avg_response_time_24h,
				PERCENTILE_CONT(0.95) WITHIN GROUP (
					ORDER BY EXTRACT(EPOCH FROM (response_ended_at - request_started_at)) * 1000
				) as p95_response_time_24h
			FROM check_runs
			WHERE deleted_at IS NULL
				AND created_at >= ? AND created_at <= ?
				AND request_started_at IS NOT NULL
				AND response_ended_at IS NOT NULL
				AND response_ended_at > request_started_at
			GROUP BY check_id
		)
		SELECT
			cb.id AS id,
			cb.type AS type,
			cb.name AS name,
			cb.interval AS interval,
			cb.last_run_at AS last_run_at,
			cb.last_status AS last_status,
			COALESCE(l24r.runs::text, '[]') AS last_24_runs,
			us.uptime_24h AS uptime_24h,
			us.uptime_7d AS uptime_7d,
			rts.avg_response_time_24h AS avg_response_time_24h,
			rts.p95_response_time_24h AS p95_response_time_24h
		FROM checks_base cb
		LEFT JOIN last_24_runs l24r ON cb.id = l24r.check_id
		LEFT JOIN uptime_stats us ON cb.id = us.check_id
		LEFT JOIN response_time_stats rts ON cb.id = rts.check_id
		ORDER BY cb.name
	`, projectID,
		twentyFourHoursAgo, now, // 24h uptime filters (4 times)
		twentyFourHoursAgo, now,
		twentyFourHoursAgo, now,
		sevenDaysAgo, now, // 7d uptime filters (4 times)
		sevenDaysAgo, now,
		sevenDaysAgo, now,
		twentyFourHoursAgo, now, // response time stats filter
	).Scan(&rows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list checks"})
		return
	}

	// Parse JSON array for last_24_runs and left-pad to 24 items
	responses := make([]CheckListResponse, len(rows))
	for i, row := range rows {
		// Temporary struct for parsing JSON from SQL
		type RunSummaryRaw struct {
			ID          *string    `json:"id"`
			Timestamp   *time.Time `json:"timestamp"`
			TotalTimeMs *float64   `json:"total_time_ms"` // SQL returns float, convert to int
			Status      *string    `json:"status"`
		}

		var last24RunsRaw []RunSummaryRaw
		if row.Last24Runs != "" && row.Last24Runs != "[]" {
			// Parse JSON array
			if err := json.Unmarshal([]byte(row.Last24Runs), &last24RunsRaw); err != nil {
				last24RunsRaw = []RunSummaryRaw{}
			}
		}

		// Convert to RunSummary with proper UUID parsing and float to int conversion
		last24Runs := make([]RunSummary, len(last24RunsRaw))
		for j, raw := range last24RunsRaw {
			var id *uuid.UUID
			if raw.ID != nil {
				if parsedID, err := uuid.Parse(*raw.ID); err == nil {
					id = &parsedID
				}
			}

			// Convert float to int (round down)
			var totalTimeMs *int
			if raw.TotalTimeMs != nil {
				ms := int(math.Floor(*raw.TotalTimeMs))
				totalTimeMs = &ms
			}

			last24Runs[j] = RunSummary{
				ID:          id,
				Timestamp:   raw.Timestamp,
				TotalTimeMs: totalTimeMs,
				Status:      raw.Status,
			}
		}

		// Left-pad with empty datapoints to ensure we always have 24 items
		// Empty datapoints go on the left, actual runs on the right (sorted ascending)
		paddedRuns := make([]RunSummary, 24)
		numRuns := len(last24Runs)
		emptyCount := 24 - numRuns

		// Parse interval to calculate timestamps for padded runs
		interval, err := time.ParseDuration(row.Interval)
		if err != nil {
			// Default to 10 minutes if parsing fails
			interval = 10 * time.Minute
		}

		// Determine base timestamp: use first actual run's timestamp, or now if no runs
		var baseTimestamp time.Time
		if numRuns > 0 && last24Runs[0].Timestamp != nil {
			baseTimestamp = *last24Runs[0].Timestamp
		} else {
			baseTimestamp = time.Now().UTC()
		}

		// Fill left side with empty datapoints, calculating timestamps backwards from base
		for j := 0; j < emptyCount; j++ {
			// Calculate timestamp: base - (emptyCount - j) * interval
			// This ensures the first padded run is the earliest, and they progress forward
			offset := time.Duration(emptyCount-j) * interval
			timestamp := baseTimestamp.Add(-offset)
			paddedRuns[j] = RunSummary{
				TotalTimeMs: &[]int{0}[0],
				Status:      &[]string{"unknown"}[0],
				Timestamp:   &timestamp,
				ID:          &[]uuid.UUID{uuid.Nil}[0],
			}
		}

		// Fill right side with actual runs (already sorted ascending)
		for j := 0; j < numRuns; j++ {
			paddedRuns[emptyCount+j] = last24Runs[j]
		}

		// Format percentages and response times as strings with 0 decimal places (rounded down)
		var uptime24hStr, uptime7dStr, avgMsStr, p95MsStr *string

		if row.Uptime24h != nil {
			val := fmt.Sprintf("%.0f", math.Floor(*row.Uptime24h))
			uptime24hStr = &val
		}
		if row.Uptime7d != nil {
			val := fmt.Sprintf("%.0f", math.Floor(*row.Uptime7d))
			uptime7dStr = &val
		}
		if row.AvgResponseTime24h != nil {
			val := fmt.Sprintf("%.0f", math.Floor(*row.AvgResponseTime24h))
			avgMsStr = &val
		}
		if row.P95ResponseTime24h != nil {
			val := fmt.Sprintf("%.0f", math.Floor(*row.P95ResponseTime24h))
			p95MsStr = &val
		}

		responses[i] = CheckListResponse{
			ID:                   row.ID,
			Type:                 row.Type,
			Name:                 row.Name,
			LastRun:              row.LastRun,
			LastStatus:           row.LastStatus,
			Interval:             row.Interval,
			Last24Runs:           paddedRuns,
			Uptime24h:            uptime24hStr,
			Uptime7d:             uptime7dStr,
			AvgResponseTime24hMs: avgMsStr,
			P95ResponseTime24hMs: p95MsStr,
		}
	}

	c.JSON(http.StatusOK, responses)
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

// GetCheckCountsByStatus handles GET /projects/:projectId/checks/status/counts
func (h *CheckHandler) GetCheckCountsByStatus(c *gin.Context) {
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

	type StatusCount struct {
		Status string `gorm:"column:status" json:"status"`
		Count  int    `gorm:"column:count" json:"count"`
	}

	var counts []StatusCount
	err = h.store.DB().Raw(`
		SELECT
			COALESCE(last_status, 'unknown') as status,
			COUNT(*) as count
		FROM checks
		WHERE project_id = ? AND deleted_at IS NULL
		GROUP BY last_status
		ORDER BY last_status
	`, projectID).Scan(&counts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get check counts"})
		return
	}

	// Initialize all possible statuses with 0 counts
	result := map[string]int{
		"passing":  0,
		"degraded": 0,
		"failing":  0,
		"unknown":  0,
	}

	// Update with actual counts
	for _, count := range counts {
		result[count.Status] = count.Count
	}

	c.JSON(http.StatusOK, result)
}
