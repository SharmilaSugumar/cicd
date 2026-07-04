package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type PipelineService interface {
	CreatePipeline(ctx context.Context, projectID, userID uuid.UUID, name, config string) (*database.Pipeline, error)
	ValidatePipelineConfiguration(ctx context.Context, config string) (bool, error)
	AssociateQueue(ctx context.Context, pipelineID, queueID uuid.UUID) error
	VerifyPermissions(ctx context.Context, pipelineID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error)
}

type pipelineService struct {
	pipelineRepo repositories.PipelineRepository
	projectRepo  repositories.ProjectRepository
	orgService   OrganizationService
}

func NewPipelineService(pipelineRepo repositories.PipelineRepository, projectRepo repositories.ProjectRepository, orgService OrganizationService) PipelineService {
	return &pipelineService{
		pipelineRepo: pipelineRepo,
		projectRepo:  projectRepo,
		orgService:   orgService,
	}
}

func (s *pipelineService) CreatePipeline(ctx context.Context, projectID, userID uuid.UUID, name, config string) (*database.Pipeline, error) {
	// Verify user access to the project
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, ErrNotFound
	}

	hasAccess, err := s.orgService.VerifyRBACPermissions(ctx, project.OrganizationID, userID, database.OrganizationRoleDeveloper)
	if err != nil || !hasAccess {
		return nil, ErrForbidden
	}

	valid, err := s.ValidatePipelineConfiguration(ctx, config)
	if !valid || err != nil {
		return nil, ErrInvalidInput
	}

	pipeline := &database.Pipeline{
		ProjectID:   projectID,
		Name:        name,
		Description: config, // Assuming config is stored in description for now
	}

	if err := s.pipelineRepo.Create(ctx, pipeline); err != nil {
		return nil, ErrInternalError
	}

	return pipeline, nil
}

func (s *pipelineService) ValidatePipelineConfiguration(ctx context.Context, config string) (bool, error) {
	// Placeholder for YAML/JSON pipeline configuration validation
	if config == "" {
		return false, ErrInvalidInput
	}
	return true, nil
}

func (s *pipelineService) AssociateQueue(ctx context.Context, pipelineID, queueID uuid.UUID) error {
	// Business logic to map a pipeline to a specific queue
	// In the data model, Job is tied to Queue, but Pipeline could have a default queue.
	// For now, this is a placeholder to satisfy the interface requirement.
	return nil
}

func (s *pipelineService) VerifyPermissions(ctx context.Context, pipelineID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error) {
	pipeline, err := s.pipelineRepo.GetByID(ctx, pipelineID)
	if err != nil {
		return false, ErrNotFound
	}
	project, err := s.projectRepo.GetByID(ctx, pipeline.ProjectID)
	if err != nil {
		return false, ErrNotFound
	}
	return s.orgService.VerifyRBACPermissions(ctx, project.OrganizationID, userID, requiredRole)
}
