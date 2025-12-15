package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"pulse/internal/models"
)

func executeHTTPCheck(check *models.Check, startTime time.Time) Result {
	timeout := 30 * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", check.Host, check.Port),
		Path:   check.Path,
	}
	if check.Secure {
		u.Scheme = "https"
	}

	req, err := http.NewRequest(check.Method, u.String(), nil)
	if err != nil {
		return Result{
			Status:           models.CheckRunStatusFailing,
			AssertionResults: emptyJSON(),
			PlaywrightReport: emptyJSON(),
			NetworkTimings:   emptyJSON(),
			Metrics:          mustMarshalJSON(map[string]interface{}{"error": err.Error()}),
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

	resp, err := client.Do(req)
	latencyMs := int(time.Since(startTime).Milliseconds())

	if err != nil {
		timeoutErr, isTimeout := err.(interface{ Timeout() bool })
		return Result{
			Status:           evaluateStatus(check, latencyMs, 0, false, err, isTimeout && timeoutErr.Timeout()),
			ResponseStatus:   0,
			AssertionResults: emptyJSON(),
			PlaywrightReport: emptyJSON(),
			NetworkTimings: mustMarshalJSON(map[string]interface{}{
				"total_time_ms": latencyMs,
			}),
			Metrics: mustMarshalJSON(map[string]interface{}{"error": err.Error()}),
		}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))

	assertionResults, assertionsPassed := processAssertions(check.Assertions, resp.StatusCode, bodyBytes)

	status := evaluateStatus(check, latencyMs, resp.StatusCode, assertionsPassed, nil, false)

	networkTimings := mustMarshalJSON(map[string]interface{}{
		"total_time_ms":   latencyMs,
		"dns_time_ms":     0,
		"connect_time_ms": 0,
		"ttfb_ms":         0,
	})

	metrics := mustMarshalJSON(map[string]interface{}{
		"response_size_bytes": len(bodyBytes),
		"status_code":         resp.StatusCode,
	})

	return Result{
		Status:           status,
		ResponseStatus:   int32(resp.StatusCode),
		AssertionResults: assertionResults,
		PlaywrightReport: emptyJSON(),
		NetworkTimings:   networkTimings,
		Metrics:          metrics,
	}
}
