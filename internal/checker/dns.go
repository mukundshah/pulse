package checker

import (
	"time"

	"pulse/internal/models"
)

func executeDNSCheck(_ *models.Check, _ time.Time) Result {
	return Result{
		Status:           models.CheckRunStatusUnknown,
		AssertionResults: emptyJSON(),
		PlaywrightReport: emptyJSON(),
		NetworkTimings:   emptyJSON(),
		Metrics:          mustMarshalJSON(map[string]interface{}{"error": "DNS checks not yet implemented"}),
	}
}
