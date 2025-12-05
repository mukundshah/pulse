package alerter

import (
	"pulse/internal/checker"
	"pulse/internal/models"
	"pulse/internal/store"
)

type Alerter struct {
	store *store.Store
}

func New(s *store.Store) *Alerter {
	return &Alerter{
		store: s,
	}
}

// ProcessCheckResult processes the check result
// TODO: Implement alerting logic
func (a *Alerter) ProcessCheckResult(check *models.Check, result checker.Result) {
	// TODO: Implement alerting logic
}
