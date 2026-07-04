package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(ctx context.Context, name, description string, orgID, userID uuid.UUID) (*database.Project, error)
	ListProjectsByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]database.Project, error)
	VerifyOrganizationMembership(ctx context.Context, orgID, userID uuid.UUID) (bool, error)
	SoftDeleteProject(ctx context.Context, projectID, userID uuid.UUID) error
}

type projectService struct {
	projectRepo repositories.ProjectRepository
	orgService  OrganizationService
}

func NewProjectService(projectRepo repositories.ProjectRepository, orgService OrganizationService) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		orgService:  orgService,
	}
}

func (s *projectService) CreateProject(ctx context.Context, name, description string, orgID, userID uuid.UUID) (*database.Project, error) {
	// Verify user is at least a maintainer to create a project
	hasAccess, err := s.orgService.VerifyRBACPermissions(ctx, orgID, userID, database.OrganizationRoleMaintainer)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrForbidden
	}

	project := &database.Project{
		Name:           name,
		Description:    description,
		OrganizationID: orgID,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, ErrInternalError
	}
	return project, nil
}

func (s *projectService) ListProjectsByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]database.Project, error) {
	projects, err := s.projectRepo.ListByOrganization(ctx, orgID, limit, offset)
	if err != nil {
		return nil, ErrInternalError
	}
	return projects, nil
}

func (s *projectService) VerifyOrganizationMembership(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
	return s.orgService.VerifyRBACPermissions(ctx, orgID, userID, database.OrganizationRoleViewer)
}

func (s *projectService) SoftDeleteProject(ctx context.Context, projectID, userID uuid.UUID) error {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return ErrNotFound
	}

	// Verify user is owner or maintainer
	hasAccess, err := s.orgService.VerifyRBACPermissions(ctx, project.OrganizationID, userID, database.OrganizationRoleMaintainer)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrForbidden
	}

	if err := s.projectRepo.Delete(ctx, projectID); err != nil {
		return ErrInternalError
	}
	return nil
}
