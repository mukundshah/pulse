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
	"pulse/internal/redis"
	"pulse/internal/store"
)

type Worker struct {
	store       *store.Store
	runsStore   *store.RunsStore
	redis       *redis.Client
	alerter     *alerter.Alerter
	workerCount int
	quit        chan struct{}
	wg          sync.WaitGroup
}

func New(s *store.Store, rs *store.RunsStore, r *redis.Client, a *alerter.Alerter, workerCount int) *Worker {
	return &Worker{
		store:       s,
		runsStore:   rs,
		redis:       r,
		alerter:     a,
		workerCount: workerCount,
		quit:        make(chan struct{}),
	}
}

func (w *Worker) Start() {
	log.Printf("Starting %d workers...", w.workerCount)
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

	// Execute check
	result := checker.Execute(check)
	metrics.IncrementChecksExecuted()

	// Track failures
	if result.Status != "success" {
		metrics.IncrementChecksFailed()
	}

	// Save result to ClickHouse
	if err := w.runsStore.CreateCheckRun(ctx, result); err != nil {
		log.Printf("Worker %d: Error saving check run for %s: %v", workerID, checkID, err)
	}

	// Process alerts and get new failure count
	newFailures := w.alerter.ProcessCheckResult(check, result)

	// Update check status
	nextRun := time.Now().Add(time.Duration(check.IntervalSeconds) * time.Second)
	lastStatus := result.Status
	if err := w.store.UpdateCheckStatus(checkID, nextRun, newFailures, lastStatus); err != nil {
		log.Printf("Worker %d: Error updating check status for %s: %v", workerID, checkID, err)
	}

	log.Printf("Worker %d: Check %s executed: %s (Latency: %dms, Failures: %d)",
		workerID, check.Name, result.Status, result.LatencyMs, newFailures)
}
