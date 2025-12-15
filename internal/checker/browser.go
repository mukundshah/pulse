package checker

import (
	"time"

	"pulse/internal/models"
)

func executeBrowserCheck(_ *models.Check, _ time.Time) Result {
	return Result{
		Status:           models.CheckRunStatusUnknown,
		AssertionResults: emptyJSON(),
		PlaywrightReport: emptyJSON(),
		NetworkTimings:   emptyJSON(),
		Metrics:          mustMarshalJSON(map[string]interface{}{"error": "Browser checks not yet implemented"}),
	}
}
