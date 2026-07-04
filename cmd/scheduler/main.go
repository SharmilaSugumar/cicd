package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"

	"forgeflow/internal/database"
	"forgeflow/internal/repositories"
	"forgeflow/internal/scheduler"
	"forgeflow/pkg/logger"
)

func main() {
	// 1. Init Logger
	logger.InitLogger()
	defer logger.Sync()
	logger.Log.Info("Starting ForgeFlow Scheduler...")

	// 2. Init DB
	if err := database.ConnectDatabase(); err != nil {
		logger.Log.Fatal("Failed to connect to DB: " + err.Error())
	}

	// 3. Init Repositories
	jobRepo := repositories.NewJobRepository(database.DB)
	queueRepo := repositories.NewQueueRepository(database.DB)

	// 4. Init and Start Scheduler
	sched := scheduler.NewScheduler(jobRepo, queueRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched.Start(ctx)

	// 5. Wait for termination
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down Scheduler...")
	sched.Stop()
	logger.Log.Info("Scheduler exiting cleanly")
}
