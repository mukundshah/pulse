package checker

import (
	"time"

	"pulse/internal/models"
)

// Execute runs a check and returns the result
func Execute(check *models.Check) Result {
	startTime := time.Now()

	switch check.Type {
	case models.CheckTypeHTTP:
		return executeHTTPCheck(check, startTime)
	case models.CheckTypeTCP:
		return executeTCPCheck(check, startTime)
	case models.CheckTypeDNS:
		return executeDNSCheck(check, startTime)
	case models.CheckTypeBrowser:
		return executeBrowserCheck(check, startTime)
	case models.CheckTypeHeartbeat:
		return executeHeartbeatCheck(check, startTime)
	default:
		return Result{
			Status:           models.CheckRunStatusUnknown,
			AssertionResults: emptyJSON(),
			PlaywrightReport: emptyJSON(),
			NetworkTimings:   emptyJSON(),
			Metrics:          mustMarshalJSON(map[string]interface{}{"error": "unknown check type"}),
		}
	}
}
