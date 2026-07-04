package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type JobService interface {
	UpdateStatus(ctx context.Context, jobID uuid.UUID, newStatus database.JobStatus) error
	RetryFailedJob(ctx context.Context, jobID uuid.UUID) error
	VerifyDependencyCompletion(ctx context.Context, jobID uuid.UUID) (bool, error)
	TransitionJobState(ctx context.Context, jobID uuid.UUID, newStatus database.JobStatus) error
	PublishJobEvent(ctx context.Context, jobID uuid.UUID, eventType string) error
}

type jobService struct {
	jobRepo repositories.JobRepository
}

func NewJobService(jobRepo repositories.JobRepository) JobService {
	return &jobService{jobRepo: jobRepo}
}

func (s *jobService) UpdateStatus(ctx context.Context, jobID uuid.UUID, newStatus database.JobStatus) error {
	if err := s.jobRepo.UpdateStatus(ctx, jobID, newStatus); err != nil {
		return ErrInternalError
	}
	return s.PublishJobEvent(ctx, jobID, "status_updated")
}

func (s *jobService) RetryFailedJob(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return ErrNotFound
	}
	if job.Status != database.JobStatusFailed {
		return ErrInvalidStateTransition
	}
	return s.TransitionJobState(ctx, jobID, database.JobStatusRetrying)
}

func (s *jobService) VerifyDependencyCompletion(ctx context.Context, jobID uuid.UUID) (bool, error) {
	_, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return false, ErrNotFound
	}

	return s.jobRepo.AreDependenciesCompleted(ctx, jobID)
}

func (s *jobService) TransitionJobState(ctx context.Context, jobID uuid.UUID, newStatus database.JobStatus) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return ErrNotFound
	}

	validTransitions := map[database.JobStatus][]database.JobStatus{
		database.JobStatusCreated:  {database.JobStatusQueued},
		database.JobStatusQueued:   {database.JobStatusClaimed},
		database.JobStatusClaimed:  {database.JobStatusRunning, database.JobStatusFailed},
		database.JobStatusRunning:  {database.JobStatusCompleted, database.JobStatusFailed},
		database.JobStatusFailed:   {database.JobStatusRetrying, database.JobStatusDeadLetter},
		database.JobStatusRetrying: {database.JobStatusQueued},
	}

	allowed := false
	for _, status := range validTransitions[job.Status] {
		if status == newStatus {
			allowed = true
			break
		}
	}

	if !allowed {
		return ErrInvalidStateTransition
	}

	return s.UpdateStatus(ctx, jobID, newStatus)
}

func (s *jobService) PublishJobEvent(ctx context.Context, jobID uuid.UUID, eventType string) error {
	// Placeholder for WebSocket or message broker publishing
	// Example: websocket.Broadcast(...)
	return nil
}
