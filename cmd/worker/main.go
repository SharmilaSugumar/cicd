package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forgeflow/internal/database"
	"forgeflow/internal/execution"
	"forgeflow/internal/repositories"
	"forgeflow/internal/workers"
	"forgeflow/pkg/logger"

	"github.com/google/uuid"
)

func main() {
	// 1. Init Logger
	logger.InitLogger()
	defer logger.Sync()
	logger.Log.Info("Starting ForgeFlow Worker...")

	// 2. Init DB
	if err := database.ConnectDatabase(); err != nil {
		logger.Log.Fatal("Failed to connect to DB: " + err.Error())
	}

	if err := database.AutoMigrate(); err != nil {
		logger.Log.Fatal("Failed to migrate DB: " + err.Error())
	}

	// 3. Init Repositories
	workerRepo := repositories.NewWorkerRepository(database.DB)
	jobRepo := repositories.NewJobRepository(database.DB)

	// 4. Init Execution Engine
	engine := execution.NewEngine()

	// 5. Setup Worker Config
	name := os.Getenv("WORKER_NAME")
	if name == "" {
		name = "default-worker"
	}

	var queueID *uuid.UUID
	if q := os.Getenv("WORKER_QUEUE_ID"); q != "" {
		parsed, err := uuid.Parse(q)
		if err == nil {
			queueID = &parsed
		} else {
			logger.Log.Warn("Invalid WORKER_QUEUE_ID format, proceeding with all queues")
		}
	}

	cfg := workers.WorkerConfig{
		Name:              name,
		QueueID:           queueID,
		PollInterval:      5 * time.Second,
		HeartbeatInterval: 10 * time.Second,
	}

	// 6. Init and Start Worker
	worker := workers.NewWorker(cfg, workerRepo, jobRepo, database.DB, engine)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := worker.Start(ctx); err != nil {
		logger.Log.Fatal("Failed to start worker: " + err.Error())
	}

	// 7. Wait for Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down worker...")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()

	if err := worker.Stop(stopCtx); err != nil {
		logger.Log.Fatal("Worker forced to shutdown: " + err.Error())
	}

	logger.Log.Info("Worker exited gracefully")
}
