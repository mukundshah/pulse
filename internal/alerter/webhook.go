package alerter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pulse/internal/models"
	"time"
)

type WebhookPayload struct {
	CheckID   string    `json:"check_id"`
	CheckName string    `json:"check_name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func SendWebhook(webhookURL string, check *models.Check, alertType string, result *models.CheckRun) error {
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

	// Retry logic: 3 attempts with exponential backoff
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			if attempt < maxRetries-1 {
				delay := baseDelay * time.Duration(1<<uint(attempt))
				time.Sleep(delay)
				continue
			}
			return fmt.Errorf("failed to send webhook after %d attempts: %w", maxRetries, err)
		}

		// Read response body to ensure connection is closed
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(attempt))
			time.Sleep(delay)
			continue
		}

		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return fmt.Errorf("failed to send webhook after %d attempts", maxRetries)
}
