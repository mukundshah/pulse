package checker

import (
	"time"

	"gorm.io/datatypes"

	"pulse/internal/models"
)

// Result represents the outcome of a check execution.
type Result struct {
	Status         models.CheckRunStatus
	FailureReason  *models.FailureReason
	ResponseStatus *int32

	// Timestamps (authoritative)
	RequestStartedAt time.Time
	FirstByteAt      time.Time
	ResponseEndedAt  time.Time

	// Metadata
	ConnectionReused  bool
	IPVersion         string
	IPAddress         string
	ResponseSizeBytes int64

	// JSON fields
	AssertionResults datatypes.JSON
	PlaywrightReport datatypes.JSON
	NetworkTimings   datatypes.JSON
	Response         datatypes.JSON

	Error error
}

func emptyJSONArray() datatypes.JSON {
	return mustMarshalJSON([]interface{}{})
}

func emptyJSONObject() datatypes.JSON {
	return mustMarshalJSON(map[string]interface{}{})
}
