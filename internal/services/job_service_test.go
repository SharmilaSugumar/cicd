package services

import (
	"context"
	"forgeflow/internal/database"
	"testing"

	"github.com/google/uuid"
)

type mockJobRepository struct {
	GetByIDFunc           func(ctx context.Context, id uuid.UUID) (*database.Job, error)
	UpdateStatusFunc      func(ctx context.Context, id uuid.UUID, status database.JobStatus) error
	CreateFunc            func(ctx context.Context, job *database.Job) error
	UpdateFunc            func(ctx context.Context, job *database.Job) error
	DeleteFunc            func(ctx context.Context, id uuid.UUID) error
	ListFunc              func(ctx context.Context, limit, offset int) ([]database.Job, error)
	ListByPipelineRunFunc func(ctx context.Context, runID uuid.UUID, limit, offset int) ([]database.Job, error)
	ClaimNextJobFunc      func(ctx context.Context, queueID *uuid.UUID) (*database.Job, error)
	AreDependenciesCompletedFunc func(ctx context.Context, jobID uuid.UUID) (bool, error)
	GetDependenciesFunc   func(ctx context.Context, jobID uuid.UUID) ([]database.JobDependency, error)
}

func (m *mockJobRepository) AreDependenciesCompleted(ctx context.Context, jobID uuid.UUID) (bool, error) {
	if m.AreDependenciesCompletedFunc != nil {
		return m.AreDependenciesCompletedFunc(ctx, jobID)
	}
	return true, nil
}

func (m *mockJobRepository) GetDependencies(ctx context.Context, jobID uuid.UUID) ([]database.JobDependency, error) {
	if m.GetDependenciesFunc != nil {
		return m.GetDependenciesFunc(ctx, jobID)
	}
	return nil, nil
}

func (m *mockJobRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Job, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockJobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status database.JobStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}
func (m *mockJobRepository) Create(ctx context.Context, job *database.Job) error { return nil }
func (m *mockJobRepository) Update(ctx context.Context, job *database.Job) error { return nil }
func (m *mockJobRepository) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *mockJobRepository) List(ctx context.Context, limit, offset int) ([]database.Job, error) {
	return nil, nil
}
func (m *mockJobRepository) ListByPipelineRun(ctx context.Context, runID uuid.UUID, limit, offset int) ([]database.Job, error) {
	return nil, nil
}
func (m *mockJobRepository) ClaimNextJob(ctx context.Context, queueID *uuid.UUID) (*database.Job, error) {
	return nil, nil
}

func TestUpdateStatus(t *testing.T) {
	mockRepo := &mockJobRepository{
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status database.JobStatus) error {
			return nil
		},
	}
	service := NewJobService(mockRepo)
	err := service.UpdateStatus(context.Background(), uuid.New(), database.JobStatusRunning)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRetryFailedJob_Success(t *testing.T) {
	jobID := uuid.New()
	mockRepo := &mockJobRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*database.Job, error) {
			return &database.Job{Status: database.JobStatusFailed}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status database.JobStatus) error {
			return nil
		},
	}
	service := NewJobService(mockRepo)
	err := service.RetryFailedJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRetryFailedJob_InvalidState(t *testing.T) {
	jobID := uuid.New()
	mockRepo := &mockJobRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*database.Job, error) {
			return &database.Job{Status: database.JobStatusRunning}, nil
		},
	}
	service := NewJobService(mockRepo)
	err := service.RetryFailedJob(context.Background(), jobID)
	if err != ErrInvalidStateTransition {
		t.Fatalf("expected ErrInvalidStateTransition, got %v", err)
	}
}

func TestTransitionJobState_Valid(t *testing.T) {
	jobID := uuid.New()
	mockRepo := &mockJobRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*database.Job, error) {
			return &database.Job{Status: database.JobStatusQueued}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status database.JobStatus) error {
			return nil
		},
	}
	service := NewJobService(mockRepo)
	err := service.TransitionJobState(context.Background(), jobID, database.JobStatusClaimed)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTransitionJobState_Invalid(t *testing.T) {
	jobID := uuid.New()
	mockRepo := &mockJobRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*database.Job, error) {
			return &database.Job{Status: database.JobStatusQueued}, nil
		},
	}
	service := NewJobService(mockRepo)
	err := service.TransitionJobState(context.Background(), jobID, database.JobStatusCompleted)
	if err != ErrInvalidStateTransition {
		t.Fatalf("expected ErrInvalidStateTransition, got %v", err)
	}
}
