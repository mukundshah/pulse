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
	"pulse/internal/handlers"
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

	// Run migrations
	if err := db.Migrate(pgDB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to ClickHouse
	chClient, err := clickhouse.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer chClient.Close()

	// Initialize schema
	ctx := context.Background()
	if err := chClient.InitSchema(ctx); err != nil {
		log.Fatalf("Failed to initialize ClickHouse schema: %v", err)
	}

	// Connect to Redis
	redisClient, err := redis.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create stores with cache
	s := store.NewWithCache(pgDB, redisClient)
	rs := store.NewRunsStore(chClient)

	// Create alerter
	a := alerter.New(s, rs)

	// Create scheduler
	sched := scheduler.New(s, redisClient)
	sched.Start()
	defer sched.Stop()

	// Create and start workers
	workerCount := 3
	w := worker.New(s, rs, redisClient, a, workerCount)
	w.Start()
	defer w.Stop()

	// Setup health/metrics HTTP server
	healthHandler := handlers.NewHealthHandler(pgDB, redisClient, chClient)
	healthRouter := gin.Default()
	healthRouter.GET("/health", healthHandler.Health)
	healthRouter.GET("/ready", healthHandler.Ready)
	healthRouter.GET("/metrics", healthHandler.Metrics)

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
