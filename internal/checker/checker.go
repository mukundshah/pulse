package checker

import (
	"time"

	"pulse/internal/models"
)

// Execute runs a check and returns the result
func Execute(check *models.Check) Result {
	attempts := computeAttempts(check)
	var result Result
	var previousDelay time.Duration

	for attempt := 0; attempt < attempts; attempt++ {
		result = executeOnce(check)
		if result.Status == models.CheckRunStatusPassing {
			return result
		}

		if attempt == attempts-1 {
			break
		}

		delay := computeRetryDelay(check, attempt, previousDelay)
		previousDelay = delay
		if delay > 0 {
			time.Sleep(delay)
		}
	}

	return result
}
