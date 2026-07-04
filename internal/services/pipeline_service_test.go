package services

import (
	"context"
	"forgeflow/internal/database"
	"testing"

	"github.com/google/uuid"
)

type mockPipelineRepository struct {
	CreateFunc        func(ctx context.Context, pipeline *database.Pipeline) error
	GetByIDFunc       func(ctx context.Context, id uuid.UUID) (*database.Pipeline, error)
	UpdateFunc        func(ctx context.Context, pipeline *database.Pipeline) error
	DeleteFunc        func(ctx context.Context, id uuid.UUID) error
	ListFunc          func(ctx context.Context, limit, offset int) ([]database.Pipeline, error)
	ListByProjectFunc func(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]database.Pipeline, error)
}

func (m *mockPipelineRepository) Create(ctx context.Context, pipeline *database.Pipeline) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, pipeline)
	}
	return nil
}
func (m *mockPipelineRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Pipeline, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockPipelineRepository) Update(ctx context.Context, pipeline *database.Pipeline) error {
	return nil
}
func (m *mockPipelineRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockPipelineRepository) List(ctx context.Context, limit, offset int) ([]database.Pipeline, error) {
	return nil, nil
}
func (m *mockPipelineRepository) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]database.Pipeline, error) {
	return nil, nil
}

type mockProjectRepository struct {
	GetByIDFunc            func(ctx context.Context, id uuid.UUID) (*database.Project, error)
	CreateFunc             func(ctx context.Context, project *database.Project) error
	UpdateFunc             func(ctx context.Context, project *database.Project) error
	DeleteFunc             func(ctx context.Context, id uuid.UUID) error
	ListFunc               func(ctx context.Context, limit, offset int) ([]database.Project, error)
	ListByOrganizationFunc func(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]database.Project, error)
}

func (m *mockProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockProjectRepository) Create(ctx context.Context, project *database.Project) error {
	return nil
}
func (m *mockProjectRepository) Update(ctx context.Context, project *database.Project) error {
	return nil
}
func (m *mockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockProjectRepository) List(ctx context.Context, limit, offset int) ([]database.Project, error) {
	return nil, nil
}
func (m *mockProjectRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]database.Project, error) {
	return nil, nil
}

type mockOrganizationService struct {
	VerifyRBACPermissionsFunc    func(ctx context.Context, orgID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error)
	CreateOrganizationFunc       func(ctx context.Context, name string, ownerID uuid.UUID) (*database.Organization, error)
	AddMemberFunc                func(ctx context.Context, orgID, userID uuid.UUID, role database.OrganizationRole) error
	RemoveMemberFunc             func(ctx context.Context, orgID, userID uuid.UUID) error
	ListOrganizationsForUserFunc func(ctx context.Context, userID uuid.UUID) ([]database.Organization, error)
}

func (m *mockOrganizationService) VerifyRBACPermissions(ctx context.Context, orgID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error) {
	if m.VerifyRBACPermissionsFunc != nil {
		return m.VerifyRBACPermissionsFunc(ctx, orgID, userID, requiredRole)
	}
	return false, nil
}
func (m *mockOrganizationService) CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*database.Organization, error) {
	return nil, nil
}
func (m *mockOrganizationService) AddMember(ctx context.Context, orgID, userID uuid.UUID, role database.OrganizationRole) error {
	return nil
}
func (m *mockOrganizationService) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	return nil
}
func (m *mockOrganizationService) ListOrganizationsForUser(ctx context.Context, userID uuid.UUID) ([]database.Organization, error) {
	return nil, nil
}

func TestCreatePipeline_Success(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	mockProjectRepo := &mockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*database.Project, error) {
			return &database.Project{OrganizationID: uuid.New()}, nil
		},
	}
	mockOrgService := &mockOrganizationService{
		VerifyRBACPermissionsFunc: func(ctx context.Context, orgID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error) {
			return true, nil
		},
	}
	mockPipelineRepo := &mockPipelineRepository{
		CreateFunc: func(ctx context.Context, pipeline *database.Pipeline) error {
			return nil
		},
	}

	service := NewPipelineService(mockPipelineRepo, mockProjectRepo, mockOrgService)
	pipeline, err := service.CreatePipeline(context.Background(), projectID, userID, "Test Pipeline", "valid-config")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pipeline == nil {
		t.Fatal("expected pipeline to be returned")
	}
}

func TestCreatePipeline_Forbidden(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	mockProjectRepo := &mockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*database.Project, error) {
			return &database.Project{OrganizationID: uuid.New()}, nil
		},
	}
	mockOrgService := &mockOrganizationService{
		VerifyRBACPermissionsFunc: func(ctx context.Context, orgID, userID uuid.UUID, requiredRole database.OrganizationRole) (bool, error) {
			return false, nil // no access
		},
	}

	service := NewPipelineService(&mockPipelineRepository{}, mockProjectRepo, mockOrgService)
	_, err := service.CreatePipeline(context.Background(), projectID, userID, "Test Pipeline", "valid-config")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
