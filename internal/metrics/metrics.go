package metrics

import (
	"sync/atomic"
	"time"
)

// Metrics tracks system-wide metrics using atomic operations for thread-safety
type Metrics struct {
	ChecksExecuted int64
	ChecksFailed   int64
	AlertsSent     int64
	ActiveJobs     int64
	StartTime      time.Time
}

var globalMetrics = &Metrics{
	StartTime: time.Now(),
}

// GetMetrics returns the current metrics snapshot
func GetMetrics() *Metrics {
	return &Metrics{
		ChecksExecuted: atomic.LoadInt64(&globalMetrics.ChecksExecuted),
		ChecksFailed:   atomic.LoadInt64(&globalMetrics.ChecksFailed),
		AlertsSent:     atomic.LoadInt64(&globalMetrics.AlertsSent),
		ActiveJobs:     atomic.LoadInt64(&globalMetrics.ActiveJobs),
		StartTime:      globalMetrics.StartTime,
	}
}

// IncrementChecksExecuted increments the checks executed counter
func IncrementChecksExecuted() {
	atomic.AddInt64(&globalMetrics.ChecksExecuted, 1)
}

// IncrementChecksFailed increments the checks failed counter
func IncrementChecksFailed() {
	atomic.AddInt64(&globalMetrics.ChecksFailed, 1)
}

// IncrementAlertsSent increments the alerts sent counter
func IncrementAlertsSent() {
	atomic.AddInt64(&globalMetrics.AlertsSent, 1)
}

// IncrementActiveJobs increments the active jobs counter
func IncrementActiveJobs() {
	atomic.AddInt64(&globalMetrics.ActiveJobs, 1)
}

// DecrementActiveJobs decrements the active jobs counter
func DecrementActiveJobs() {
	atomic.AddInt64(&globalMetrics.ActiveJobs, -1)
}

// GetUptime returns the uptime in seconds
func GetUptime() int64 {
	return int64(time.Since(globalMetrics.StartTime).Seconds())
}
