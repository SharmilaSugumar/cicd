package services

import (
	"context"
	"encoding/json"
	"forgeflow/internal/database"
	"forgeflow/internal/execution"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type PipelineRunService interface {
	CreatePipelineRun(ctx context.Context, pipelineID, userID uuid.UUID) (*database.PipelineRun, error)
	GenerateExecutionPlan(ctx context.Context, run *database.PipelineRun) ([]database.Job, error)
	CreateJobs(ctx context.Context, jobs []database.Job) error
	InitializePipelineStatus(ctx context.Context, runID uuid.UUID) error
	ListRunsByPipeline(ctx context.Context, pipelineID uuid.UUID) ([]database.PipelineRun, error)
	GetRunDetails(ctx context.Context, runID uuid.UUID) (*database.PipelineRun, error)
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

type pipelineConfig struct {
	SourceType string `json:"source_type"`
	RepoURL    string `json:"repo_url"`
}

func (s *pipelineRunService) GenerateExecutionPlan(ctx context.Context, run *database.PipelineRun) ([]database.Job, error) {
	pipeline, err := s.pipelineSvc.GetPipeline(ctx, run.PipelineID)
	if err != nil {
		return nil, ErrNotFound
	}

	var cfg pipelineConfig
	if pipeline.YamlConfig != "" {
		_ = json.Unmarshal([]byte(pipeline.YamlConfig), &cfg)
	}

	repoUrl := cfg.RepoURL

	payloadObj := execution.JobPayload{
		RepoURL: repoUrl,
		Branch:  "",
		// We leave Dependencies, BuildCmds, TestCmds empty so the worker auto-detects them
	}

	payloadBytes, _ := json.Marshal(payloadObj)

	jobs := []database.Job{
		{
			PipelineRunID: run.ID,
			Name:          "build-and-test",
			Status:        database.JobStatusCreated,
			Priority:      10,
			Payload:       string(payloadBytes),
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

func (s *pipelineRunService) ListRunsByPipeline(ctx context.Context, pipelineID uuid.UUID) ([]database.PipelineRun, error) {
	return s.runRepo.ListByPipelineID(ctx, pipelineID)
}

func (s *pipelineRunService) GetRunDetails(ctx context.Context, runID uuid.UUID) (*database.PipelineRun, error) {
	return s.runRepo.GetRunDetails(ctx, runID)
}
