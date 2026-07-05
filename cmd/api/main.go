package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forgeflow/internal/database"
	"forgeflow/internal/handlers"
	"forgeflow/internal/middleware"
	"forgeflow/internal/repositories"
	"forgeflow/internal/services"
	"forgeflow/internal/websocket"
	"forgeflow/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1. Init Logger
	logger.InitLogger()
	defer logger.Sync()
	logger.Log.Info("Starting ForgeFlow API server...")

	// 2. Init DB
	if err := database.ConnectDatabase(); err != nil {
		logger.Log.Fatal("Failed to connect to DB: " + err.Error())
	}
	if err := database.AutoMigrate(); err != nil {
		logger.Log.Fatal("Failed to migrate DB: " + err.Error())
	}

	// 3. Init Repositories
	userRepo := repositories.NewUserRepository(database.DB)
	orgRepo := repositories.NewOrganizationRepository(database.DB)
	projectRepo := repositories.NewProjectRepository(database.DB)
	pipelineRepo := repositories.NewPipelineRepository(database.DB)
	pipelineRunRepo := repositories.NewPipelineRunRepository(database.DB)
	jobRepo := repositories.NewJobRepository(database.DB)
	queueRepo := repositories.NewQueueRepository(database.DB)
	workerRepo := repositories.NewWorkerRepository(database.DB)

	// 4. Init Services
	authSvc := services.NewAuthService(userRepo)
	orgSvc := services.NewOrganizationService(orgRepo)
	projectSvc := services.NewProjectService(projectRepo, orgSvc)
	pipelineSvc := services.NewPipelineService(pipelineRepo, projectRepo, orgSvc)
	pipelineRunSvc := services.NewPipelineRunService(pipelineRunRepo, jobRepo, pipelineSvc)
	jobSvc := services.NewJobService(jobRepo)
	queueSvc := services.NewQueueService(queueRepo, jobRepo)
	workerSvc := services.NewWorkerService(workerRepo, jobRepo)

	// 5. Init Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	orgHandler := handlers.NewOrganizationHandler(orgSvc)
	projectHandler := handlers.NewProjectHandler(projectSvc)
	pipelineHandler := handlers.NewPipelineHandler(pipelineSvc, pipelineRunSvc)
	jobHandler := handlers.NewJobHandler(jobSvc)
	queueHandler := handlers.NewQueueHandler(queueSvc)
	workerHandler := handlers.NewWorkerHandler(workerSvc)
	metricsHandler := handlers.NewMetricsHandler(database.DB)

	// 6. Setup Gin
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 7. Middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())
	r.Use(cors.Default())
	// 8. Health Routes (exempt from rate limit)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/ready", func(c *gin.Context) {
		sqlDB, err := database.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "database unready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	r.Use(middleware.RateLimit(100))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	wsHub := websocket.NewHub()
	go wsHub.Run()

	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(wsHub, c.Writer, c.Request)
	})

	v1 := r.Group("/api/v1")
	{
		authHandler.RegisterRoutes(v1.Group("/auth"))

		// Public worker routes
		workerHandler.RegisterRoutes(v1.Group("/workers"))

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authSvc))
		{
			orgHandler.RegisterRoutes(protected.Group("/organizations"))
			projectHandler.RegisterRoutes(protected.Group("/projects"))
			pipelineHandler.RegisterRoutes(protected.Group("/pipelines"))
			jobHandler.RegisterRoutes(protected.Group("/jobs"))
			queueHandler.RegisterRoutes(protected.Group("/queues"))
			metricsHandler.RegisterRoutes(protected.Group("/metrics"))
		}
	}

	// 9. Server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown: " + err.Error())
	}
	logger.Log.Info("Server exiting")
}
