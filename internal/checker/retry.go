package checker

import (
	"math"
	"math/rand"
	"time"

	"pulse/internal/models"
)

func executeOnce(check *models.Check) Result {
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

	baseDelay := durationFrom(check.RetriesDelay, check.RetriesDelayUnit)
	maxDelay := durationFrom(check.RetriesMaxDelay, check.RetriesMaxDelayUnit)
	delay := baseDelay

	switch check.Retries {
	case models.RetryTypeFixed:
		delay = baseDelay
	case models.RetryTypeLinear:
		delay = time.Duration(int64(baseDelay) * int64(attempt+1))
	case models.RetryTypeExponential:
		factor := 2.0
		if check.RetriesFactor != nil && *check.RetriesFactor > 0 {
			factor = *check.RetriesFactor
		}
		delay = time.Duration(float64(baseDelay) * math.Pow(factor, float64(attempt)))
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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

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

func durationFrom(value *int, unit *models.UnitType) time.Duration {
	if value == nil {
		return 0
	}

	multiplier := time.Millisecond
	if unit != nil && *unit == models.UnitTypeS {
		multiplier = time.Second
	}

	return time.Duration(*value) * multiplier
}
