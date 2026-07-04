package scheduler

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"
	"forgeflow/pkg/logger"
	"os"
	"time"

	"go.uber.org/zap"
)

type Scheduler interface {
	Start(ctx context.Context)
	Stop()
}

type schedulerImpl struct {
	jobRepo   repositories.JobRepository
	queueRepo repositories.QueueRepository
	stopChan  chan struct{}
}

func NewScheduler(jobRepo repositories.JobRepository, queueRepo repositories.QueueRepository) Scheduler {
	return &schedulerImpl{
		jobRepo:   jobRepo,
		queueRepo: queueRepo,
		stopChan:  make(chan struct{}),
	}
}

func (s *schedulerImpl) Start(ctx context.Context) {
	logger.Log.Info("Scheduler daemon starting...")

	intervalStr := os.Getenv("SCHEDULER_INTERVAL")
	if intervalStr == "" {
		intervalStr = "5s"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Scheduler context cancelled")
				return
			case <-s.stopChan:
				logger.Log.Info("Scheduler stopped gracefully")
				return
			case <-ticker.C:
				s.dispatchJobs(ctx)
			}
		}
	}()
}

func (s *schedulerImpl) Stop() {
	close(s.stopChan)
}

func (s *schedulerImpl) dispatchJobs(ctx context.Context) {
	// Find CREATED jobs where dependencies are met
	// Update their status to QUEUED

	// Note: In a production app, we would use a more complex query
	// joining on dependencies to find runnable jobs efficiently.
	// For now, we simulate fetching eligible jobs.

	jobs, err := s.jobRepo.List(ctx, 100, 0)
	if err != nil {
		logger.Log.Error("Failed to fetch jobs in scheduler", zap.Error(err))
		return
	}

	for _, job := range jobs {
		if job.Status == database.JobStatusCreated {
			// Validate dependencies
			dependenciesMet, err := s.jobRepo.AreDependenciesCompleted(ctx, job.ID)
			if err != nil {
				logger.Log.Error("Failed to check dependencies", zap.Error(err))
				continue
			}

			if dependenciesMet {
				// We can queue this job
				if err := s.jobRepo.UpdateStatus(ctx, job.ID, database.JobStatusQueued); err != nil {
					logger.Log.Error("Failed to queue job", zap.String("job_id", job.ID.String()), zap.Error(err))
				} else {
					logger.Log.Info("Job queued successfully", zap.String("job_id", job.ID.String()))
				}
			}
		}
	}
}
