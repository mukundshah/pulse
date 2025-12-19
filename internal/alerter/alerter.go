package alerter

import (
	"log"

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

// ProcessCheckResult processes the check result and creates an alert if the status has changed
// check is the check that was executed, run is the check run that was just created
func (a *Alerter) ProcessCheckResult(check *models.Check, run *models.CheckRun) {
	if check.LastStatus == models.CheckRunStatusUnknown {
		return
	}

	if check.LastStatus == run.Status {
		return
	}

	alert := &models.Alert{
		Status:    run.Status,
		RunID:     run.ID,
		RegionID:  run.RegionID,
		ProjectID: check.ProjectID,
		CheckID:   check.ID,
	}

	if err := a.store.CreateAlert(alert); err != nil {
		log.Printf("Error creating alert for check %s: %v", check.ID, err)
		return
	}

	log.Printf("Alert created for check %s: status changed from %s to %s", check.Name, check.LastStatus, run.Status)
}
