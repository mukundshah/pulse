package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"pulse/internal/clickhouse"
	"pulse/internal/config"
	"pulse/internal/db"
	"pulse/internal/handlers"
	"pulse/internal/redis"
	"pulse/internal/store"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	pgDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.Migrate(pgDB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to ClickHouse (optional, will work without it)
	chClient, err := clickhouse.NewClient(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to ClickHouse: %v", err)
	} else {
		defer chClient.Close()
		ctx := context.Background()
		if err := chClient.InitSchema(ctx); err != nil {
			log.Printf("Warning: Failed to initialize ClickHouse schema: %v", err)
		}
	}

	// Connect to Redis for caching (optional, will work without it)
	redisClient, err := redis.Connect(cfg)
	var s *store.Store
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis, caching disabled: %v", err)
		s = store.New(pgDB)
	} else {
		defer redisClient.Close()
		s = store.NewWithCache(pgDB, redisClient)
	}

	// Setup handlers
	checksHandler := handlers.NewChecksHandler(s)

	// Setup routes
	r := gin.Default()
	r.GET("/checks", checksHandler.ListChecks)
	r.POST("/checks", checksHandler.CreateCheck)
	r.GET("/checks/:id", checksHandler.GetCheck)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/checks")
	})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
