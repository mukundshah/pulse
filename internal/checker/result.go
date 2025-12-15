package checker

import (
	"gorm.io/datatypes"

	"pulse/internal/models"
)

// Result represents the outcome of a check execution.
type Result struct {
	Status           models.CheckRunStatus
	ResponseStatus   int32
	AssertionResults datatypes.JSON
	PlaywrightReport datatypes.JSON
	NetworkTimings   datatypes.JSON
	Metrics          datatypes.JSON
}

func emptyJSON() datatypes.JSON {
	return mustMarshalJSON(map[string]interface{}{})
}
