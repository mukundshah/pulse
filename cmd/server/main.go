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
	"pulse/internal/middleware"
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
	authHandler := handlers.NewAuthHandler(s, cfg)
	invitesHandler := handlers.NewInvitesHandler(s)
	membersHandler := handlers.NewMembersHandler(s)

	r.GET("/docs/:version", (func(c *gin.Context) {
		version := c.Param("version")

		if version != "v1" {
			c.JSON(404, gin.H{
				"error": "Unsupported version",
			})
			return
		}

		html, err := scalargo.NewV2(
			scalargo.WithSpecDir(filepath.Join(cfg.APISpecDir, version)),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Pulse API"),
				scalargo.WithKeyValue("description", "Pulse API"),
			),
			scalargo.WithSpecModifier(func(spec *scalarModel.Spec) *scalarModel.Spec {
				localhost := "Localhost"
				spec.Servers = []scalarModel.Server{
					{URL: "http://localhost:8080/api", Description: &localhost},
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
		// Auth routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/password/reset", authHandler.ForgotPassword)
		api.POST("/auth/password/reset/done", authHandler.ResetPassword)
	}

	protected := api.Group("/").Use(middleware.AuthMiddleware(cfg))

	{
		// Project routes - specific routes first to avoid conflicts
		protected.POST("/projects", projectHandler.CreateProject)
		protected.GET("/projects", projectHandler.ListProjects)

		// Project sub-resources (must come before /projects/:projectId)
		protected.POST("/projects/:projectId/checks", checkHandler.CreateCheck)
		protected.GET("/projects/:projectId/checks", checkHandler.ListChecks)
		protected.POST("/projects/:projectId/tags", tagHandler.CreateTag)
		protected.GET("/projects/:projectId/tags", tagHandler.ListTags)
		protected.POST("/projects/:projectId/tags/:tagId", tagHandler.AddTagToProject)
		protected.DELETE("/projects/:projectId/tags/:tagId", tagHandler.RemoveTagFromProject)
		protected.POST("/projects/:projectId/invites", invitesHandler.CreateInvite)
		protected.GET("/projects/:projectId/invites", invitesHandler.ListInvites)
		protected.GET("/projects/:projectId/members", membersHandler.ListMembers)
		protected.PUT("/projects/:projectId/members/:userId", membersHandler.UpdateMemberRole)
		protected.DELETE("/projects/:projectId/members/:userId", membersHandler.RemoveMember)

		// Project CRUD (generic routes come after specific ones)
		protected.GET("/projects/:projectId", projectHandler.GetProject)
		protected.PUT("/projects/:projectId", projectHandler.UpdateProject)
		protected.DELETE("/projects/:projectId", projectHandler.DeleteProject)

		// Check routes - specific routes first
		protected.GET("/checks/:checkId/runs", checkRunHandler.ListCheckRuns)
		protected.POST("/checks/:checkId/tags/:tagId", tagHandler.AddTagToCheck)
		protected.DELETE("/checks/:checkId/tags/:tagId", tagHandler.RemoveTagFromCheck)

		// Check CRUD (generic routes come after specific ones)
		protected.GET("/checks/:checkId", checkHandler.GetCheck)
		protected.PUT("/checks/:checkId", checkHandler.UpdateCheck)
		protected.DELETE("/checks/:checkId", checkHandler.DeleteCheck)

		// CheckRun routes
		protected.GET("/check-runs/:id", checkRunHandler.GetCheckRun)

		// Invite routes
		protected.POST("/invites/accept", invitesHandler.AcceptInvite)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
