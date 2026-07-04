package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *database.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.Project, error)
	Update(ctx context.Context, project *database.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.Project, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]database.Project, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *database.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Project, error) {
	var project database.Project
	err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error
	return &project, err
}

func (r *projectRepository) Update(ctx context.Context, project *database.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.Project{}, "id = ?", id).Error
}

func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]database.Project, error) {
	var projects []database.Project
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&projects).Error
	return projects, err
}

func (r *projectRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]database.Project, error) {
	var projects []database.Project
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Limit(limit).Offset(offset).Find(&projects).Error
	return projects, err
}
