package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pulse/internal/models"

	"github.com/google/uuid"
)

// Execute runs an HTTP check and returns the result
func Execute(check *models.Check) *models.CheckRun {
	start := time.Now()

	// Default timeout if not specified
	timeout := 10 * time.Second
	if check.TimeoutMs > 0 {
		timeout = time.Duration(check.TimeoutMs) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, check.Method, check.URL, nil)
	if err != nil {
		return &models.CheckRun{
			ID:        uuid.New(),
			CheckID:   check.ID,
			Status:    "error",
			LatencyMs: time.Since(start).Milliseconds(),
			Error:     stringPtr(err.Error()),
			RunAt:     start,
		}
	}

	// Add headers if present
	if check.Headers != nil && len(check.Headers) > 0 {
		var headers map[string]interface{}
		headersBytes, err := check.Headers.MarshalJSON()
		if err == nil {
			if err := json.Unmarshal(headersBytes, &headers); err == nil {
				for k, v := range headers {
					if str, ok := v.(string); ok {
						req.Header.Set(k, str)
					}
				}
			}
		}
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		status := "fail"
		errMsg := err.Error()
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			status = "timeout"
		}
		return &models.CheckRun{
			ID:        uuid.New(),
			CheckID:   check.ID,
			Status:    status,
			LatencyMs: latency.Milliseconds(),
			Error:     stringPtr(errMsg),
			RunAt:     start,
		}
	}
	defer resp.Body.Close()

	// Read body if we need to check content
	var bodyMatched = true
	if check.BodyContains != nil && *check.BodyContains != "" {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return &models.CheckRun{
				ID:         uuid.New(),
				CheckID:    check.ID,
				Status:     "error",
				LatencyMs:  latency.Milliseconds(),
				StatusCode: int32(resp.StatusCode),
				Error:      stringPtr("failed to read body: " + err.Error()),
				RunAt:      start,
			}
		}
		if !strings.Contains(string(bodyBytes), *check.BodyContains) {
			bodyMatched = false
		}
	}

	status := "success"
	var errMsg *string

	if resp.StatusCode != check.ExpectedStatus {
		status = "fail"
		msg := fmt.Sprintf("unexpected status code: got %d, expected %d", resp.StatusCode, check.ExpectedStatus)
		errMsg = &msg
	} else if !bodyMatched {
		status = "fail"
		msg := "body content mismatch"
		errMsg = &msg
	}

	return &models.CheckRun{
		ID:         uuid.New(),
		CheckID:    check.ID,
		Status:     status,
		LatencyMs:  latency.Milliseconds(),
		StatusCode: int32(resp.StatusCode),
		Error:      errMsg,
		RunAt:      start,
	}
}

func stringPtr(s string) *string {
	return &s
}
