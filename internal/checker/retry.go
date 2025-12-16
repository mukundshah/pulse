package checker

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"pulse/internal/models"
)

func executeOnce(check *models.Check) Result {
	switch check.Type {
	case models.CheckTypeHTTP:
		return ExecuteHTTPCheck(context.Background(), check)
	// case models.CheckTypeTCP:
	// 	return executeTCPCheck(check)
	// case models.CheckTypeDNS:
	// 	return executeDNSCheck(check)
	// case models.CheckTypeBrowser:
	// 	return executeBrowserCheck(check)
	// case models.CheckTypeHeartbeat:
	// 	return executeHeartbeatCheck(check)
	default:
		now := time.Now().UTC()
		err := fmt.Errorf("unsupported check type: %s", check.Type)
		return Result{
			Status:            models.CheckRunStatusFailing,
			FailureReason:     failureReasonPtr(models.FailureUnknown),
			ResponseStatus:    nil,
			RequestStartedAt:  now,
			FirstByteAt:       time.Time{},
			ResponseEndedAt:   time.Time{},
			ConnectionReused:  false,
			IPVersion:         "",
			IPAddress:         "",
			ResponseSizeBytes: 0,
			AssertionResults:  emptyJSONArray(),
			PlaywrightReport:  emptyJSONObject(),
			NetworkTimings:    emptyJSONObject(),
			Error:             err,
		}
	}
}

func computeAttempts(check *models.Check) int {
	if check.Retries == models.RetryTypeNone {
		return 1
	}

	if check.RetriesCount != nil && *check.RetriesCount > 0 {
		return 1 + *check.RetriesCount
	}

	return 1
}

func computeRetryDelay(check *models.Check, attempt int, previousDelay time.Duration) time.Duration {
	if check.Retries == models.RetryTypeNone {
		return 0
	}

	baseDelay := check.RetriesDelayDuration()
	maxDelay := check.RetriesMaxDelayDuration()
	delay := baseDelay

	switch check.Retries {
	case models.RetryTypeFixed:
		delay = baseDelay
	case models.RetryTypeLinear:
		delay = baseDelay * time.Duration(attempt+1)
	case models.RetryTypeExponential:
		factor := 2.0
		if check.RetriesFactor != nil && *check.RetriesFactor > 0 {
			factor = *check.RetriesFactor
		}
		delay = baseDelay * time.Duration(math.Pow(factor, float64(attempt)))
	}

	delay = applyJitter(delay, previousDelay, check)

	if maxDelay > 0 && delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

func applyJitter(delay time.Duration, previousDelay time.Duration, check *models.Check) time.Duration {
	if delay <= 0 {
		return 0
	}

	jitterType := models.RetryJitterTypeNone
	if check.RetriesJitter != nil {
		jitterType = *check.RetriesJitter
	}

	jitterFactor := 1.0
	if check.RetriesJitterFactor != nil && *check.RetriesJitterFactor > 0 {
		jitterFactor = *check.RetriesJitterFactor
	}

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	switch jitterType {
	case models.RetryJitterTypeFull:
		return time.Duration(r.Float64() * float64(delay) * jitterFactor)
	case models.RetryJitterTypeEqual:
		min := float64(delay) * 0.5 * jitterFactor
		max := float64(delay) * jitterFactor
		return time.Duration(min + r.Float64()*(max-min))
	case models.RetryJitterTypeDecorrelated:
		min := float64(delay)
		max := float64(delay)
		if previousDelay > 0 {
			max = float64(previousDelay) * 3
		}
		max *= jitterFactor
		if max < min {
			max = min
		}
		return time.Duration(min + r.Float64()*(max-min))
	default:
		return delay
	}
}
