package handlers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pulse/internal/clickhouse"
	"pulse/internal/db"
	"pulse/internal/metrics"
	"pulse/internal/redis"
)

type HealthHandler struct {
	pgDB        *gorm.DB
	redisClient *redis.Client
	chClient    *clickhouse.Client
	startTime   time.Time
}

func NewHealthHandler(pgDB *gorm.DB, redisClient *redis.Client, chClient *clickhouse.Client) *HealthHandler {
	return &HealthHandler{
		pgDB:        pgDB,
		redisClient: redisClient,
		chClient:    chClient,
		startTime:   time.Now(),
	}
}

type HealthResponse struct {
	Status        string            `json:"status"`
	UptimeSeconds int64             `json:"uptime_seconds"`
	Dependencies  map[string]string `json:"dependencies"`
	Version       string            `json:"version"`
	BuildTime     string            `json:"build_time"`
}

// Health returns the health status of the service and its dependencies
func (h *HealthHandler) Health(c *gin.Context) {
	dependencies := make(map[string]string)
	overallStatus := "ok"

	// Check PostgreSQL
	if err := db.HealthCheck(h.pgDB); err != nil {
		dependencies["postgres"] = "down"
		overallStatus = "degraded"
	} else {
		dependencies["postgres"] = "ok"
	}

	// Check Redis
	if h.redisClient != nil {
		if err := h.redisClient.HealthCheck(); err != nil {
			dependencies["redis"] = "down"
			overallStatus = "degraded"
		} else {
			dependencies["redis"] = "ok"
		}
	} else {
		dependencies["redis"] = "unavailable"
	}

	// Check ClickHouse
	if h.chClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := h.chClient.HealthCheck(ctx); err != nil {
			dependencies["clickhouse"] = "down"
			overallStatus = "degraded"
		} else {
			dependencies["clickhouse"] = "ok"
		}
	} else {
		dependencies["clickhouse"] = "unavailable"
	}

	uptime := int64(time.Since(h.startTime).Seconds())

	response := HealthResponse{
		Status:        overallStatus,
		UptimeSeconds: uptime,
		Dependencies:  dependencies,
		Version:       "v0.1.0",
		BuildTime:     h.startTime.Format(time.RFC3339),
	}

	c.JSON(200, response)
}

// Ready returns 200 OK only when all critical dependencies are reachable
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check PostgreSQL
	if err := db.HealthCheck(h.pgDB); err != nil {
		c.JSON(503, gin.H{"status": "not ready", "error": "postgres unreachable"})
		return
	}

	// Check Redis (required for worker)
	if h.redisClient != nil {
		if err := h.redisClient.HealthCheck(); err != nil {
			c.JSON(503, gin.H{"status": "not ready", "error": "redis unreachable"})
			return
		}
	}

	c.JSON(200, gin.H{"status": "ready"})
}

type MetricsResponse struct {
	ChecksExecutedTotal int64   `json:"checks_executed_total"`
	ChecksFailedTotal   int64   `json:"checks_failed_total"`
	AlertsSentTotal     int64   `json:"alerts_sent_total"`
	AvgLatencyMs        float64 `json:"avg_latency_ms"`
	Worker              struct {
		ActiveJobs int64 `json:"active_jobs"`
		QueueDepth int64 `json:"queue_depth"`
		Uptime     int64 `json:"uptime"`
	} `json:"worker"`
}

// Metrics returns system metrics
func (h *HealthHandler) Metrics(c *gin.Context) {
	m := metrics.GetMetrics()

	var avgLatency float64
	if h.chClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		latency, err := h.chClient.GetAverageLatency(ctx)
		if err == nil {
			avgLatency = latency
		}
	}

	var queueDepth int64
	if h.redisClient != nil {
		depth, err := h.redisClient.GetQueueDepth()
		if err == nil {
			queueDepth = depth
		}
	}

	response := MetricsResponse{
		ChecksExecutedTotal: m.ChecksExecuted,
		ChecksFailedTotal:   m.ChecksFailed,
		AlertsSentTotal:     m.AlertsSent,
		AvgLatencyMs:        avgLatency,
	}
	response.Worker.ActiveJobs = m.ActiveJobs
	response.Worker.QueueDepth = queueDepth
	response.Worker.Uptime = metrics.GetUptime()

	c.JSON(200, response)
}
