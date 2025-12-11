package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gorm.io/datatypes"

	"pulse/internal/models"
)

// Result represents the result of executing a check
type Result struct {
	Status           models.CheckRunStatus
	ResponseStatus   int32
	AssertionResults datatypes.JSON
	PlaywrightReport *datatypes.JSON // optional
	NetworkTimings   datatypes.JSON
	Metrics          datatypes.JSON
}

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
			Status:  models.CheckRunStatusError,
			Metrics: mustMarshalJSON(map[string]interface{}{"error": "unknown check type"}),
		}
	}
}

func executeHTTPCheck(check *models.Check, startTime time.Time) Result {
	// Create HTTP client with timeout (default 10 seconds)
	timeout := 30 * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", check.Host, check.Port),
	}
	if check.Secure {
		url.Scheme = "https"
	}

	// Create request
	req, err := http.NewRequest(check.Method, url.String(), nil)
	if err != nil {
		return Result{
			Status:  models.CheckRunStatusError,
			Metrics: mustMarshalJSON(map[string]interface{}{"error": err.Error()}),
		}
	}

	// Set headers if provided
	if check.Headers != nil {
		var headers map[string]interface{}
		if err := json.Unmarshal(check.Headers, &headers); err == nil {
			for k, v := range headers {
				if str, ok := v.(string); ok {
					req.Header.Set(k, str)
				}
			}
		}
	}

	// Execute request
	resp, err := client.Do(req)
	latencyMs := int(time.Since(startTime).Milliseconds())

	if err != nil {
		// Check if it's a timeout
		if timeoutErr, ok := err.(interface{ Timeout() bool }); ok && timeoutErr.Timeout() {
			return Result{
				Status:         models.CheckRunStatusTimeout,
				ResponseStatus: 0,
				NetworkTimings: mustMarshalJSON(map[string]interface{}{
					"total_time_ms": latencyMs,
				}),
				Metrics: mustMarshalJSON(map[string]interface{}{"error": "request timeout"}),
			}
		}
		return Result{
			Status:         models.CheckRunStatusError,
			ResponseStatus: 0,
			NetworkTimings: mustMarshalJSON(map[string]interface{}{
				"total_time_ms": latencyMs,
			}),
			Metrics: mustMarshalJSON(map[string]interface{}{"error": err.Error()}),
		}
	}
	defer resp.Body.Close()

	// Read response body (limited to prevent memory issues)
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB limit

	// Determine status based on response (default to success for 2xx status codes)
	status := models.CheckRunStatusSuccess
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = models.CheckRunStatusFail
	}

	// Process assertions if provided
	assertionResults := processAssertions(check.Assertions, resp.StatusCode, bodyBytes)

	// Build network timings
	networkTimings := mustMarshalJSON(map[string]interface{}{
		"total_time_ms":   latencyMs,
		"dns_time_ms":     0, // Could be extracted from http.Transport if needed
		"connect_time_ms": 0,
		"ttfb_ms":         0, // Time to first byte
	})

	// Build metrics
	metrics := mustMarshalJSON(map[string]interface{}{
		"response_size_bytes": len(bodyBytes),
		"status_code":         resp.StatusCode,
	})

	return Result{
		Status:           status,
		ResponseStatus:   int32(resp.StatusCode),
		AssertionResults: assertionResults,
		NetworkTimings:   networkTimings,
		Metrics:          metrics,
	}
}

func executeTCPCheck(check *models.Check, startTime time.Time) Result {
	// TODO: Implement TCP check
	return Result{
		Status:  models.CheckRunStatusError,
		Metrics: mustMarshalJSON(map[string]interface{}{"error": "TCP checks not yet implemented"}),
	}
}

func executeDNSCheck(check *models.Check, startTime time.Time) Result {
	// TODO: Implement DNS check
	return Result{
		Status:  models.CheckRunStatusError,
		Metrics: mustMarshalJSON(map[string]interface{}{"error": "DNS checks not yet implemented"}),
	}
}

func executeBrowserCheck(check *models.Check, startTime time.Time) Result {
	// TODO: Implement browser check with Playwright
	return Result{
		Status:  models.CheckRunStatusError,
		Metrics: mustMarshalJSON(map[string]interface{}{"error": "Browser checks not yet implemented"}),
	}
}

func executeHeartbeatCheck(check *models.Check, startTime time.Time) Result {
	// Heartbeat check just verifies the system is running
	return Result{
		Status:         models.CheckRunStatusSuccess,
		ResponseStatus: 200,
		NetworkTimings: mustMarshalJSON(map[string]interface{}{
			"total_time_ms": int(time.Since(startTime).Milliseconds()),
		}),
		Metrics: mustMarshalJSON(map[string]interface{}{
			"type": "heartbeat",
		}),
	}
}

func processAssertions(assertions datatypes.JSON, statusCode int, body []byte) datatypes.JSON {
	if assertions == nil || len(assertions) == 0 {
		return mustMarshalJSON(map[string]interface{}{})
	}

	var assertionRules []map[string]interface{}
	if err := json.Unmarshal(assertions, &assertionRules); err != nil {
		return mustMarshalJSON(map[string]interface{}{"error": "invalid assertions format"})
	}

	results := make(map[string]interface{})
	for i, rule := range assertionRules {
		ruleType, _ := rule["type"].(string)
		passed := false

		switch ruleType {
		case "status_code":
			expected, _ := rule["value"].(float64)
			passed = statusCode == int(expected)
		case "body_contains":
			expected, _ := rule["value"].(string)
			passed = bytes.Contains(body, []byte(expected))
		case "body_not_contains":
			expected, _ := rule["value"].(string)
			passed = !bytes.Contains(body, []byte(expected))
		case "response_time":
			// This would need to be passed from the check execution
			// For now, just mark as not implemented
			passed = false
		}

		results[fmt.Sprintf("assertion_%d", i)] = map[string]interface{}{
			"type":    ruleType,
			"passed":  passed,
			"message": rule["message"],
		}
	}

	// Determine overall assertion result
	allPassed := true
	for _, v := range results {
		if assertion, ok := v.(map[string]interface{}); ok {
			if passed, ok := assertion["passed"].(bool); ok && !passed {
				allPassed = false
				break
			}
		}
	}

	return mustMarshalJSON(map[string]interface{}{
		"all_passed": allPassed,
		"results":    results,
	})
}

func mustMarshalJSON(v interface{}) datatypes.JSON {
	data, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON("{}")
	}
	return datatypes.JSON(data)
}
