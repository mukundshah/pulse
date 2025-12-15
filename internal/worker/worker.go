package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"pulse/internal/alerter"
	"pulse/internal/checker"
	"pulse/internal/metrics"
	"pulse/internal/models"
	"pulse/internal/redis"
	"pulse/internal/store"
)

type Worker struct {
	store       *store.Store
	redis       *redis.Client
	alerter     *alerter.Alerter
	workerCount int
	regionID    uuid.UUID
	quit        chan struct{}
	wg          sync.WaitGroup
}

func New(s *store.Store, r *redis.Client, a *alerter.Alerter, workerCount int, regionID uuid.UUID) *Worker {
	return &Worker{
		store:       s,
		redis:       r,
		alerter:     a,
		workerCount: workerCount,
		regionID:    regionID,
		quit:        make(chan struct{}),
	}
}

func (w *Worker) Start() {
	log.Printf("Starting %d workers for region %s...", w.workerCount, w.regionID)
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.work(i)
	}
}

func (w *Worker) Stop() {
	log.Println("Stopping workers...")
	close(w.quit)
	w.wg.Wait()
	log.Println("Workers stopped.")
}

func (w *Worker) work(id int) {
	defer w.wg.Done()
	ctx := context.Background()

	for {
		select {
		case <-w.quit:
			return
		default:
			// Dequeue job with 5 second timeout
			checkID, err := w.redis.DequeueJob(5 * time.Second)
			if err != nil {
				// No job available, continue
				continue
			}

			w.processCheck(ctx, checkID, id)
		}
	}
}

func (w *Worker) processCheck(ctx context.Context, checkID uuid.UUID, workerID int) {
	metrics.IncrementActiveJobs()
	defer metrics.DecrementActiveJobs()

	// Load check from database
	check, err := w.store.GetCheck(checkID)
	if err != nil {
		log.Printf("Worker %d: Error loading check %s: %v", workerID, checkID, err)
		return
	}

	// Verify check has this region enabled
	hasRegion := false
	for _, region := range check.Regions {
		if region.ID == w.regionID {
			hasRegion = true
			break
		}
	}
	if !hasRegion {
		log.Printf("Worker %d: Check %s does not have region %s enabled, skipping", workerID, checkID, w.regionID)
		return
	}

	// Execute check
	result := checker.Execute(check)
	metrics.IncrementChecksExecuted()

	// Track failures
	if result.Status != models.CheckRunStatusPassing {
		metrics.IncrementChecksFailed()
	}

	// Convert checker.Result to models.CheckRun
	checkRun := &models.CheckRun{
		Status:           result.Status,
		TotalTimeMs:      result.TotalTimeMs,
		ResponseStatus:   result.ResponseStatus,
		AssertionResults: result.AssertionResults,
		PlaywrightReport: result.PlaywrightReport,
		NetworkTimings:   result.NetworkTimings,
		RegionID:         w.regionID,
		CheckID:          check.ID,
	}

	if result.Error != nil {
		checkRun.Remarks = result.Error.Error()
	}

	if err := w.store.CreateCheckRun(checkRun); err != nil {
		log.Printf("Worker %d: Error saving check run for %s: %v", workerID, checkID, err)
	}

	// Process alerts
	w.alerter.ProcessCheckResult(check, result)

	// Parse interval and update check status
	interval, err := time.ParseDuration(check.Interval)
	if err != nil {
		log.Printf("Worker %d: Error parsing interval for check %s: %v", workerID, checkID, err)
		interval = 10 * time.Minute // Default to 10 minutes if parsing fails
	}
	nextRun := time.Now().Add(interval)
	if err := w.store.UpdateCheckStatus(checkID, nextRun, result.Status); err != nil {
		log.Printf("Worker %d: Error updating check status for %s: %v", workerID, checkID, err)
	}

	log.Printf("Worker %d: Check %s executed: %s", workerID, check.Name, result.Status)
}
