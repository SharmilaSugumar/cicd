package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PipelineRepository interface {
	Create(ctx context.Context, pipeline *database.Pipeline) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.Pipeline, error)
	Update(ctx context.Context, pipeline *database.Pipeline) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.Pipeline, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]database.Pipeline, error)
}

type pipelineRepository struct {
	db *gorm.DB
}

func NewPipelineRepository(db *gorm.DB) PipelineRepository {
	return &pipelineRepository{db: db}
}

func (r *pipelineRepository) Create(ctx context.Context, pipeline *database.Pipeline) error {
	return r.db.WithContext(ctx).Create(pipeline).Error
}

func (r *pipelineRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Pipeline, error) {
	var pipeline database.Pipeline
	err := r.db.WithContext(ctx).First(&pipeline, "id = ?", id).Error
	return &pipeline, err
}

func (r *pipelineRepository) Update(ctx context.Context, pipeline *database.Pipeline) error {
	return r.db.WithContext(ctx).Save(pipeline).Error
}

func (r *pipelineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.Pipeline{}, "id = ?", id).Error
}

func (r *pipelineRepository) List(ctx context.Context, limit, offset int) ([]database.Pipeline, error) {
	var pipelines []database.Pipeline
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&pipelines).Error
	return pipelines, err
}

func (r *pipelineRepository) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]database.Pipeline, error) {
	var pipelines []database.Pipeline
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Limit(limit).Offset(offset).Find(&pipelines).Error
	return pipelines, err
}
