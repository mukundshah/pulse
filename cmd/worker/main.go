package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"

	"pulse/internal/alerter"
	"pulse/internal/clickhouse"
	"pulse/internal/config"
	"pulse/internal/db"
	"pulse/internal/redis"
	"pulse/internal/scheduler"
	"pulse/internal/store"
	"pulse/internal/worker"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	pgDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to ClickHouse (optional, for health checks)
	var chClient *clickhouse.Client
	chClient, err = clickhouse.NewClient(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to ClickHouse: %v", err)
	} else {
		defer chClient.Close()
		ctx := context.Background()
		if err := chClient.InitSchema(ctx); err != nil {
			log.Printf("Warning: Failed to initialize ClickHouse schema: %v", err)
		}
	}

	// Connect to Redis
	redisClient, err := redis.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create stores with cache
	s := store.NewWithCache(pgDB, redisClient)

	// Get region from config
	if cfg.RegionCode == "" {
		log.Fatalf("REGION_CODE is required")
	}
	region, err := s.GetRegionByCode(cfg.RegionCode)
	if err != nil {
		log.Fatalf("Failed to get region by code %s: %v", cfg.RegionCode, err)
	}
	log.Printf("Worker running in region: %s (%s)", region.Name, region.Code)

	// Create alerter
	a := alerter.New(s)

	// Create scheduler with region code
	sched := scheduler.New(s, redisClient, cfg.RegionCode, scheduler.DefaultConfig())
	sched.Start()
	defer sched.Stop()

	// Create and start workers
	workerCount := 3
	w := worker.New(s, redisClient, a, workerCount, region.ID)
	w.Start()
	defer w.Stop()

	// Setup health/metrics HTTP server
	// chClient can be nil, health handler should handle it
	healthRouter := gin.Default()
	healthRouter.GET("/health", (func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}))

	// Calculate health port as API port + 1
	apiPort, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatalf("Invalid PORT configuration: %v", err)
	}
	healthPort := strconv.Itoa(apiPort + 1)

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Worker health/metrics server starting on port %s", healthPort)
		if err := healthRouter.Run(":" + healthPort); err != nil && err != http.ErrServerClosed {
			log.Printf("Health server failed: %v", err)
		}
	}()

	log.Println("Worker process started. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
