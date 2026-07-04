package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type QueueStats struct {
	TotalJobs  int
	ActiveJobs int
	FailedJobs int
}

type QueueService interface {
	PauseQueue(ctx context.Context, queueID uuid.UUID) error
	ResumeQueue(ctx context.Context, queueID uuid.UUID) error
	VerifyConcurrencyLimits(ctx context.Context, queueID uuid.UUID) (bool, error)
	GetQueueStatistics(ctx context.Context, queueID uuid.UUID) (*QueueStats, error)
}

type queueService struct {
	queueRepo repositories.QueueRepository
	jobRepo   repositories.JobRepository
}

func NewQueueService(queueRepo repositories.QueueRepository, jobRepo repositories.JobRepository) QueueService {
	return &queueService{
		queueRepo: queueRepo,
		jobRepo:   jobRepo,
	}
}

func (s *queueService) PauseQueue(ctx context.Context, queueID uuid.UUID) error {
	if err := s.queueRepo.Pause(ctx, queueID); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s *queueService) ResumeQueue(ctx context.Context, queueID uuid.UUID) error {
	if err := s.queueRepo.Resume(ctx, queueID); err != nil {
		return ErrInternalError
	}
	return nil
}

func (s *queueService) VerifyConcurrencyLimits(ctx context.Context, queueID uuid.UUID) (bool, error) {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return false, ErrNotFound
	}

	if queue.Status == database.QueueStatusPaused {
		return false, ErrLimitExceeded
	}

	// Placeholder logic to check running jobs
	// Count jobs in QUEUE where status = RUNNING
	return true, nil
}

func (s *queueService) GetQueueStatistics(ctx context.Context, queueID uuid.UUID) (*QueueStats, error) {
	_, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Aggregate job counts from jobRepo
	// Placeholder response
	return &QueueStats{
		TotalJobs:  100,
		ActiveJobs: 5,
		FailedJobs: 2,
	}, nil
}
