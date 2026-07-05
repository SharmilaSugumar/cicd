package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PipelineRunRepository interface {
	Create(ctx context.Context, run *database.PipelineRun) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.PipelineRun, error)
	Update(ctx context.Context, run *database.PipelineRun) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.PipelineRun, error)
	ListByPipelineID(ctx context.Context, pipelineID uuid.UUID) ([]database.PipelineRun, error)
	GetRunDetails(ctx context.Context, id uuid.UUID) (*database.PipelineRun, error)
}

type pipelineRunRepository struct {
	db *gorm.DB
}

func NewPipelineRunRepository(db *gorm.DB) PipelineRunRepository {
	return &pipelineRunRepository{db: db}
}

func (r *pipelineRunRepository) Create(ctx context.Context, run *database.PipelineRun) error {
	return r.db.WithContext(ctx).Create(run).Error
}

func (r *pipelineRunRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.PipelineRun, error) {
	var run database.PipelineRun
	err := r.db.WithContext(ctx).First(&run, "id = ?", id).Error
	return &run, err
}

func (r *pipelineRunRepository) Update(ctx context.Context, run *database.PipelineRun) error {
	return r.db.WithContext(ctx).Save(run).Error
}

func (r *pipelineRunRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.PipelineRun{}, "id = ?", id).Error
}

func (r *pipelineRunRepository) List(ctx context.Context, limit, offset int) ([]database.PipelineRun, error) {
	var runs []database.PipelineRun
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&runs).Error
	return runs, err
}

func (r *pipelineRunRepository) ListByPipelineID(ctx context.Context, pipelineID uuid.UUID) ([]database.PipelineRun, error) {
	var runs []database.PipelineRun
	err := r.db.WithContext(ctx).Where("pipeline_id = ?", pipelineID).Order("created_at desc").Find(&runs).Error
	return runs, err
}

func (r *pipelineRunRepository) GetRunDetails(ctx context.Context, id uuid.UUID) (*database.PipelineRun, error) {
	var run database.PipelineRun
	err := r.db.WithContext(ctx).Preload("Jobs.Logs").Preload("Jobs.DLQ").Preload("Jobs").First(&run, "id = ?", id).Error
	return &run, err
}
