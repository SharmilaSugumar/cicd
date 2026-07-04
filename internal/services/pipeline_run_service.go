package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type PipelineRunService interface {
	CreatePipelineRun(ctx context.Context, pipelineID, userID uuid.UUID) (*database.PipelineRun, error)
	GenerateExecutionPlan(ctx context.Context, run *database.PipelineRun) ([]database.Job, error)
	CreateJobs(ctx context.Context, jobs []database.Job) error
	InitializePipelineStatus(ctx context.Context, runID uuid.UUID) error
}

type pipelineRunService struct {
	runRepo     repositories.PipelineRunRepository
	jobRepo     repositories.JobRepository
	pipelineSvc PipelineService
}

func NewPipelineRunService(runRepo repositories.PipelineRunRepository, jobRepo repositories.JobRepository, pipelineSvc PipelineService) PipelineRunService {
	return &pipelineRunService{
		runRepo:     runRepo,
		jobRepo:     jobRepo,
		pipelineSvc: pipelineSvc,
	}
}

func (s *pipelineRunService) CreatePipelineRun(ctx context.Context, pipelineID, userID uuid.UUID) (*database.PipelineRun, error) {
	hasAccess, err := s.pipelineSvc.VerifyPermissions(ctx, pipelineID, userID, database.OrganizationRoleDeveloper)
	if err != nil || !hasAccess {
		return nil, ErrForbidden
	}

	run := &database.PipelineRun{
		PipelineID: pipelineID,
		Status:     database.PipelineRunStatusPending,
	}

	if err := s.runRepo.Create(ctx, run); err != nil {
		return nil, ErrInternalError
	}

	return run, nil
}

func (s *pipelineRunService) GenerateExecutionPlan(ctx context.Context, run *database.PipelineRun) ([]database.Job, error) {
	// Parse pipeline YAML and generate DAG of jobs.
	// This is a placeholder since business logic/scheduling logic shouldn't be fully implemented yet.
	jobs := []database.Job{
		{
			PipelineRunID: run.ID,
			Name:          "build",
			Status:        database.JobStatusCreated,
			Priority:      10,
		},
	}
	return jobs, nil
}

func (s *pipelineRunService) CreateJobs(ctx context.Context, jobs []database.Job) error {
	for _, j := range jobs {
		if err := s.jobRepo.Create(ctx, &j); err != nil {
			return ErrInternalError
		}
	}
	return nil
}

func (s *pipelineRunService) InitializePipelineStatus(ctx context.Context, runID uuid.UUID) error {
	run, err := s.runRepo.GetByID(ctx, runID)
	if err != nil {
		return ErrNotFound
	}

	run.Status = database.PipelineRunStatusRunning
	if err := s.runRepo.Update(ctx, run); err != nil {
		return ErrInternalError
	}
	return nil
}
