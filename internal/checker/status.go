package checker

import "pulse/internal/models"

func evaluateStatus(check *models.Check, latencyMs int, respStatus int, assertionsPassed bool, execErr error, timeout bool) models.CheckRunStatus {
	if execErr != nil || timeout {
		return models.CheckRunStatusFailing
	}

	if respStatus == 0 {
		return models.CheckRunStatusUnknown
	}

	if respStatus < 200 || respStatus >= 300 {
		return models.CheckRunStatusFailing
	}

	if !assertionsPassed {
		return models.CheckRunStatusFailing
	}

	failedThresholdMs := thresholdToMs(check.FailedThreshold, check.FailedThresholdUnit)
	if failedThresholdMs > 0 && latencyMs >= failedThresholdMs {
		return models.CheckRunStatusFailing
	}

	degradedThresholdMs := thresholdToMs(check.DegradedThreshold, check.DegradedThresholdUnit)
	if degradedThresholdMs > 0 && latencyMs >= degradedThresholdMs {
		return models.CheckRunStatusDegraded
	}

	return models.CheckRunStatusPassing
}

func thresholdToMs(threshold int, unit models.UnitType) int {
	switch unit {
	case models.UnitTypeS:
		return threshold * 1000
	default:
		return threshold
	}
}
