package alerter

import (
	"log"
	"pulse/internal/metrics"
	"pulse/internal/models"
	"pulse/internal/store"
)

type Alerter struct {
	store      *store.Store
	runsStore  *store.RunsStore
	webhookURL string
}

func New(s *store.Store, rs *store.RunsStore) *Alerter {
	return &Alerter{
		store:     s,
		runsStore: rs,
	}
}

// ProcessCheckResult evaluates the result and triggers alerts if necessary.
// It returns the new consecutive failure count.
func (a *Alerter) ProcessCheckResult(check *models.Check, result *models.CheckRun) int {
	newFailures := check.ConsecutiveFails

	if result.Status != "success" {
		newFailures++
		if newFailures == check.AlertThreshold {
			a.sendAlert(check, "failure", result)
		}
	} else {
		if newFailures >= check.AlertThreshold {
			a.sendAlert(check, "recovery", result)
		}
		newFailures = 0
	}

	return newFailures
}

func (a *Alerter) sendAlert(check *models.Check, alertType string, result *models.CheckRun) {
	var errorMsg string
	if result.Error != nil {
		errorMsg = *result.Error
	}

	log.Printf("ALERT [%s]: Check %s (%s) - Status: %s, Error: %s", alertType, check.Name, check.URL, result.Status, errorMsg)

	alert := &models.Alert{
		CheckID: check.ID,
		Type:    alertType,
		Payload: errorMsg,
		SentAt:  result.RunAt,
	}

	// Persist alert
	if err := a.store.CreateAlert(alert); err != nil {
		log.Printf("Error saving alert: %v", err)
	}

	// Send webhook if configured
	if check.WebhookURL != nil && *check.WebhookURL != "" {
		if err := SendWebhook(*check.WebhookURL, check, alertType, result); err != nil {
			log.Printf("Error sending webhook: %v", err)
		} else {
			metrics.IncrementAlertsSent()
		}
	} else {
		// Still count alerts even if no webhook configured (alert was persisted)
		metrics.IncrementAlertsSent()
	}
}
