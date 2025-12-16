package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"pulse/internal/models"
	"pulse/internal/redis"
	"pulse/internal/store"
)

const (
	defaultPollInterval    = 10 * time.Second
	defaultShutdownTimeout = 30 * time.Second
	defaultMaxRetries      = 3
	defaultRetryDelay      = 5 * time.Second
)

var (
	// ErrSchedulerStopped is returned when operations are attempted on a stopped scheduler.
	ErrSchedulerStopped = errors.New("scheduler is stopped")
	// ErrShutdownTimeout is returned when graceful shutdown times out.
	ErrShutdownTimeout = errors.New("shutdown timeout exceeded")
)

// Config holds the scheduler configuration.
type Config struct {
	PollInterval    time.Duration
	ShutdownTimeout time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
	Logger          *slog.Logger
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		PollInterval:    defaultPollInterval,
		ShutdownTimeout: defaultShutdownTimeout,
		MaxRetries:      defaultMaxRetries,
		RetryDelay:      defaultRetryDelay,
		Logger:          slog.Default(),
	}
}

// Scheduler manages the periodic execution of checks.
type Scheduler struct {
	store      *store.Store
	redis      *redis.Client
	regionCode string
	config     *Config
	logger     *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Metrics
	mu              sync.RWMutex
	checksEnqueued  int64
	checksFailed    int64
	lastPollTime    time.Time
	lastPollSuccess bool
}

// New creates a new Scheduler with the given dependencies.
func New(s *store.Store, r *redis.Client, regionCode string, config *Config) *Scheduler {
	if config == nil {
		config = DefaultConfig()
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		store:      s,
		redis:      r,
		regionCode: regionCode,
		config:     config,
		logger:     config.Logger.With("component", "scheduler", "region", regionCode),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins the scheduler's polling loop.
func (s *Scheduler) Start() error {
	s.logger.Info("starting scheduler")

	s.wg.Add(1)
	go s.poller()

	return nil
}

// Stop gracefully shuts down the scheduler.
// It waits for in-flight operations to complete or times out.
func (s *Scheduler) Stop() error {
	s.logger.Info("stopping scheduler")

	// Signal cancellation
	s.cancel()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("scheduler stopped successfully")
		return nil
	case <-time.After(s.config.ShutdownTimeout):
		s.logger.Warn("scheduler shutdown timeout exceeded")
		return ErrShutdownTimeout
	}
}

// poller runs the main polling loop.
func (s *Scheduler) poller() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.config.PollInterval)
	defer ticker.Stop()

	// Run immediately on start
	s.pollOnce()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("poller received shutdown signal")
			return
		case <-ticker.C:
			s.pollOnce()
		}
	}
}

// pollOnce performs a single poll cycle.
func (s *Scheduler) pollOnce() {
	start := time.Now().UTC()
	s.logger.Debug("starting poll cycle")

	checks, err := s.getDueChecks()
	if err != nil {
		s.logger.Error("failed to get due checks", "error", err, "duration", time.Since(start))
		return
	}

	if len(checks) == 0 {
		s.logger.Debug("no checks due", "duration", time.Since(start))
		return
	}

	s.logger.Info("processing due checks", "count", len(checks), "region", s.regionCode)

	// Process checks with limited concurrency
	s.processChecks(checks)

	s.logger.Debug("poll cycle completed",
		"duration", time.Since(start),
		"checks_processed", len(checks))
}

// getDueChecks retrieves checks that are due for execution with retry logic.
func (s *Scheduler) getDueChecks() ([]models.Check, error) {
	for attempt := 0; attempt < s.config.MaxRetries; attempt++ {
		checks, err := s.store.GetDueChecks(s.regionCode)
		if err == nil {
			return checks, nil
		}

		if attempt < s.config.MaxRetries-1 {
			s.logger.Warn("failed to get due checks, retrying",
				"attempt", attempt+1,
				"max_retries", s.config.MaxRetries,
				"error", err)

			select {
			case <-s.ctx.Done():
				return nil, s.ctx.Err()
			case <-time.After(s.config.RetryDelay):
				continue
			}
		}
	}

	return nil, fmt.Errorf("failed after %d attempts", s.config.MaxRetries)
}

// processChecks processes a batch of checks concurrently with bounded parallelism.
func (s *Scheduler) processChecks(checks []models.Check) {
	// Use a semaphore pattern to limit concurrency
	const maxConcurrent = 10
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup

	for _, check := range checks {
		// Check for shutdown
		select {
		case <-s.ctx.Done():
			s.logger.Info("stopping check processing due to shutdown")
			wg.Wait()
			return
		default:
		}

		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(c models.Check) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			s.processCheck(&c)
		}(check)
	}

	wg.Wait()
}

// processCheck handles a single check: enqueuing it and updating its next run time.
func (s *Scheduler) processCheck(check *models.Check) {
	logger := s.logger.With(
		"check_id", check.ID,
		"check_name", check.Name,
		"interval", check.Interval,
	)

	// Enqueue the check
	if err := s.enqueueCheck(check); err != nil {
		logger.Error("failed to enqueue check", "error", err)
		return
	}

	// Update next run time
	if err := s.updateNextRun(check); err != nil {
		logger.Error("failed to update next run time", "error", err)
		// Don't return here - the check was enqueued successfully
	}

	logger.Info("check enqueued successfully", "next_run", check.NextRunAt)
}

// enqueueCheck adds a check to the Redis queue with retry logic.
func (s *Scheduler) enqueueCheck(check *models.Check) error {
	var err error

	for attempt := 0; attempt < s.config.MaxRetries; attempt++ {
		err = s.redis.EnqueueJob(check.ID)
		if err == nil {
			return nil
		}

		if attempt < s.config.MaxRetries-1 {
			s.logger.Warn("failed to enqueue check, retrying",
				"check_id", check.ID,
				"attempt", attempt+1,
				"error", err)

			select {
			case <-s.ctx.Done():
				return s.ctx.Err()
			case <-time.After(s.config.RetryDelay):
				continue
			}
		}
	}

	return fmt.Errorf("failed to enqueue after %d attempts: %w", s.config.MaxRetries, err)
}

// updateNextRun calculates and updates the next run time for a check.
func (s *Scheduler) updateNextRun(check *models.Check) error {
	interval := check.IntervalDuration()

	nextRun := time.Now().UTC().Add(interval)

	if err := s.store.UpdateCheckStatus(check.ID, nextRun, check.LastStatus); err != nil {
		return fmt.Errorf("failed to update check status: %w", err)
	}

	// Update check object for logging
	check.NextRunAt = &nextRun

	return nil
}
