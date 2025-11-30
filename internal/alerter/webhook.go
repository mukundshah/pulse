package alerter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"pulse/internal/models"
	"pulse/internal/store"
)

type WebhookPayload struct {
	CheckID   string    `json:"check_id"`
	CheckName string    `json:"check_name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func SendWebhook(s *store.Store, alertID uuid.UUID, webhookURL string, check *models.Check, alertType string, result *models.CheckRun) error {
	var errorMsg string
	if result.Error != nil {
		errorMsg = *result.Error
	}

	payload := WebhookPayload{
		CheckID:   check.ID.String(),
		CheckName: check.Name,
		Type:      alertType,
		Status:    result.Status,
		Error:     errorMsg,
		Timestamp: result.RunAt,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	requestBody := string(jsonData)
	requestHeaders := make(map[string]string)
	requestHeaders["Content-Type"] = "application/json"

	// Convert headers to JSON
	headersJSON, _ := json.Marshal(requestHeaders)

	// Retry logic: 3 attempts with exponential backoff
	maxRetries := 3
	baseDelay := 1 * time.Second
	timeout := 5 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		startTime := time.Now()
		attemptRecord := &models.WebhookAttempt{
			AlertID:        &alertID,
			CheckID:        check.ID,
			URL:            webhookURL,
			RequestBody:    requestBody,
			RequestHeaders: datatypes.JSON(headersJSON),
			RetryNumber:    attempt,
			Status:         "failed",
		}

		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
		if err != nil {
			errorStr := fmt.Sprintf("failed to create request: %v", err)
			attemptRecord.Error = &errorStr
			attemptRecord.Timeout = false
			_ = s.CreateWebhookAttempt(attemptRecord)
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		latency := time.Since(startTime).Milliseconds()
		attemptRecord.LatencyMs = &latency

		if err != nil {
			errorStr := err.Error()
			attemptRecord.Error = &errorStr
			attemptRecord.Timeout = strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded")
			lastErr = err

			_ = s.CreateWebhookAttempt(attemptRecord)

			if attempt < maxRetries-1 {
				delay := baseDelay * time.Duration(1<<uint(attempt))
				time.Sleep(delay)
				continue
			}
			return fmt.Errorf("failed to send webhook after %d attempts: %w", maxRetries, err)
		}

		// Read response body
		responseBodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		responseBody := string(responseBodyBytes)
		attemptRecord.ResponseBody = &responseBody
		attemptRecord.ResponseCode = &resp.StatusCode

		// Convert response headers to JSON
		responseHeaders := make(map[string]string)
		for k, v := range resp.Header {
			if len(v) > 0 {
				responseHeaders[k] = v[0]
			}
		}
		responseHeadersJSON, _ := json.Marshal(responseHeaders)
		attemptRecord.ResponseHeaders = datatypes.JSON(responseHeadersJSON)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			attemptRecord.Status = "success"
			attemptRecord.Error = nil
			_ = s.CreateWebhookAttempt(attemptRecord)
			return nil
		}

		errorStr := fmt.Sprintf("webhook returned status %d", resp.StatusCode)
		attemptRecord.Error = &errorStr
		lastErr = fmt.Errorf("webhook returned status %d", resp.StatusCode)

		_ = s.CreateWebhookAttempt(attemptRecord)

		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(attempt))
			time.Sleep(delay)
			continue
		}
	}

	return fmt.Errorf("failed to send webhook after %d attempts: %v", maxRetries, lastErr)
}
