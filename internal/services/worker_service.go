package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type WorkerService interface {
	RegisterWorker(ctx context.Context, name string) (*database.Worker, error)
	ReceiveHeartbeat(ctx context.Context, workerID uuid.UUID) error
	MarkOfflineWorkers(ctx context.Context) error
	AssignJob(ctx context.Context, workerID uuid.UUID, queueID *uuid.UUID) (*database.Job, error)
}

type workerService struct {
	workerRepo repositories.WorkerRepository
	jobRepo    repositories.JobRepository
}

func NewWorkerService(workerRepo repositories.WorkerRepository, jobRepo repositories.JobRepository) WorkerService {
	return &workerService{
		workerRepo: workerRepo,
		jobRepo:    jobRepo,
	}
}

func (s *workerService) RegisterWorker(ctx context.Context, name string) (*database.Worker, error) {
	worker := &database.Worker{
		Name: name,
	}

	if err := s.workerRepo.Register(ctx, worker); err != nil {
		return nil, ErrInternalError
	}

	return worker, nil
}

func (s *workerService) ReceiveHeartbeat(ctx context.Context, workerID uuid.UUID) error {
	_, err := s.workerRepo.GetByID(ctx, workerID)
	if err != nil {
		return ErrNotFound
	}

	if err := s.workerRepo.UpdateHeartbeat(ctx, workerID); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s *workerService) MarkOfflineWorkers(ctx context.Context) error {
	// Fetch all workers and check last heartbeat time, mark offline if stale
	// Placeholder for cron job logic
	workers, err := s.workerRepo.List(ctx, 1000, 0)
	if err != nil {
		return ErrInternalError
	}

	for _, w := range workers {
		if w.Status == database.WorkerStatusOnline {
			// Basic mock logic: Mark offline directly for the sake of interface completion
			_ = s.workerRepo.MarkOffline(ctx, w.ID)
		}
	}
	return nil
}

func (s *workerService) AssignJob(ctx context.Context, workerID uuid.UUID, queueID *uuid.UUID) (*database.Job, error) {
	_, err := s.workerRepo.GetByID(ctx, workerID)
	if err != nil {
		return nil, ErrNotFound
	}

	job, err := s.jobRepo.ClaimNextJob(ctx, queueID)
	if err != nil {
		return nil, ErrInternalError
	}
	if job == nil {
		return nil, ErrNotFound
	}

	return job, nil
}
