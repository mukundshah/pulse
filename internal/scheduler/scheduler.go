package scheduler

import (
	"log"
	"sync"
	"time"

	"pulse/internal/redis"
	"pulse/internal/store"
)

type Scheduler struct {
	store *store.Store
	redis *redis.Client
	quit  chan struct{}
	wg    sync.WaitGroup
}

func New(s *store.Store, r *redis.Client) *Scheduler {
	return &Scheduler{
		store: s,
		redis: r,
		quit:  make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	log.Println("Starting scheduler...")
	s.wg.Add(1)
	go s.poller()
}

func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	close(s.quit)
	s.wg.Wait()
	log.Println("Scheduler stopped.")
}

func (s *Scheduler) poller() {
	defer s.wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.quit:
			return
		case <-ticker.C:
			checks, err := s.store.GetDueChecks()
			if err != nil {
				log.Printf("Error getting due checks: %v", err)
				continue
			}

			for _, check := range checks {
				// Enqueue check to Redis
				if err := s.redis.EnqueueJob(check.ID); err != nil {
					log.Printf("Error enqueueing check %s: %v", check.ID, err)
					continue
				}

				// Update next_run_at
				nextRun := time.Now().Add(time.Duration(check.IntervalSeconds) * time.Second)
				if err := s.store.UpdateCheckStatus(check.ID, nextRun, check.ConsecutiveFails, check.LastStatus); err != nil {
					log.Printf("Error updating check status for %s: %v", check.ID, err)
				}

				log.Printf("Enqueued check: %s (next run: %v)", check.Name, nextRun)
			}
		}
	}
}
