package main

import (
	"context"
	"log"
	"path/filepath"

	scalargo "github.com/bdpiprava/scalar-go"
	scalarModel "github.com/bdpiprava/scalar-go/model"
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

	// Connect to ClickHouse (optional, will work without it)
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

	// Setup routes
	r := gin.Default()

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(s)
	checkHandler := handlers.NewCheckHandler(s)
	checkRunHandler := handlers.NewCheckRunHandler(s)
	tagHandler := handlers.NewTagHandler(s)

	r.GET("/docs/v1", (func(c *gin.Context) {
		html, err := scalargo.NewV2(
			scalargo.WithSpecDir(filepath.Join(cfg.APISpecDir, "v1")),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Pulse API"),
				scalargo.WithKeyValue("description", "Pulse API"),
			),
			scalargo.WithSpecModifier(func(spec *scalarModel.Spec) *scalarModel.Spec {
				localhost := "Localhost"
				spec.Servers = []scalarModel.Server{
					{URL: "http://localhost:8080/api/v1", Description: &localhost},
				}
				return spec
			}),
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Data(200, "text/html", []byte(html))
	}))

	r.GET("/health", (func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}))

	// API routes
	api := r.Group("/api/v1")
	{
		// Tag routes (no conflicts)
		api.POST("/tags", tagHandler.CreateTag)
		api.GET("/tags", tagHandler.ListTags)

		// Project routes - specific routes first to avoid conflicts
		api.POST("/projects", projectHandler.CreateProject)
		api.GET("/projects", projectHandler.ListProjects)

		// Project sub-resources (must come before /projects/:projectId)
		api.POST("/projects/:projectId/checks", checkHandler.CreateCheck)
		api.GET("/projects/:projectId/checks", checkHandler.ListChecks)
		api.POST("/projects/:projectId/tags/:tagId", tagHandler.AddTagToProject)
		api.DELETE("/projects/:projectId/tags/:tagId", tagHandler.RemoveTagFromProject)

		// Project CRUD (generic routes come after specific ones)
		api.GET("/projects/:projectId", projectHandler.GetProject)
		api.PUT("/projects/:projectId", projectHandler.UpdateProject)
		api.DELETE("/projects/:projectId", projectHandler.DeleteProject)

		// Check routes - specific routes first
		api.GET("/checks/:checkId/runs", checkRunHandler.ListCheckRuns)
		api.POST("/checks/:checkId/tags/:tagId", tagHandler.AddTagToCheck)
		api.DELETE("/checks/:checkId/tags/:tagId", tagHandler.RemoveTagFromCheck)

		// Check CRUD (generic routes come after specific ones)
		api.GET("/checks/:checkId", checkHandler.GetCheck)
		api.PUT("/checks/:checkId", checkHandler.UpdateCheck)
		api.DELETE("/checks/:checkId", checkHandler.DeleteCheck)

		// CheckRun routes
		api.GET("/check-runs/:id", checkRunHandler.GetCheckRun)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
