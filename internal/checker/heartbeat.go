package checker

import (
	"time"

	"pulse/internal/models"
)

func executeHeartbeatCheck(_ *models.Check, startTime time.Time) Result {
	return Result{
		Status:           models.CheckRunStatusPassing,
		ResponseStatus:   200,
		AssertionResults: emptyJSON(),
		PlaywrightReport: emptyJSON(),
		NetworkTimings: mustMarshalJSON(map[string]interface{}{
			"total_time_ms": int(time.Since(startTime).Milliseconds()),
		}),
		Metrics: mustMarshalJSON(map[string]interface{}{
			"type": "heartbeat",
		}),
	}
}
