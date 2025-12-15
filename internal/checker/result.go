package checker

import (
	"gorm.io/datatypes"

	"pulse/internal/models"
)

// Result represents the outcome of a check execution.
type Result struct {
	Status           models.CheckRunStatus
	ResponseStatus   int32
	TotalTimeMs      int
	AssertionResults datatypes.JSON
	PlaywrightReport datatypes.JSON
	NetworkTimings   datatypes.JSON
	Error            error
}

func emptyJSONArray() datatypes.JSON {
	return mustMarshalJSON([]interface{}{})
}

func emptyJSONObject() datatypes.JSON {
	return mustMarshalJSON(map[string]interface{}{})
}
