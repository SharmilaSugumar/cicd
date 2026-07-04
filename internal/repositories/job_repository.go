package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobRepository interface {
	Create(ctx context.Context, job *database.Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.Job, error)
	Update(ctx context.Context, job *database.Job) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.Job, error)
	ListByPipelineRun(ctx context.Context, runID uuid.UUID, limit, offset int) ([]database.Job, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status database.JobStatus) error
	ClaimNextJob(ctx context.Context, queueID *uuid.UUID) (*database.Job, error)
	GetDependencies(ctx context.Context, jobID uuid.UUID) ([]database.JobDependency, error)
	AreDependenciesCompleted(ctx context.Context, jobID uuid.UUID) (bool, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *database.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Job, error) {
	var job database.Job
	err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error
	return &job, err
}

func (r *jobRepository) Update(ctx context.Context, job *database.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.Job{}, "id = ?", id).Error
}

func (r *jobRepository) List(ctx context.Context, limit, offset int) ([]database.Job, error) {
	var jobs []database.Job
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) ListByPipelineRun(ctx context.Context, runID uuid.UUID, limit, offset int) ([]database.Job, error) {
	var jobs []database.Job
	err := r.db.WithContext(ctx).Where("pipeline_run_id = ?", runID).Limit(limit).Offset(offset).Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status database.JobStatus) error {
	return r.db.WithContext(ctx).Model(&database.Job{}).Where("id = ?", id).Update("status", status).Error
}

func (r *jobRepository) ClaimNextJob(ctx context.Context, queueID *uuid.UUID) (*database.Job, error) {
	var job database.Job

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", database.JobStatusQueued).
			Order("priority desc, created_at asc")

		if queueID != nil {
			query = query.Where("queue_id = ?", queueID)
		}

		if err := query.First(&job).Error; err != nil {
			return err
		}

		job.Status = database.JobStatusClaimed
		return tx.Save(&job).Error
	})

	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetDependencies(ctx context.Context, jobID uuid.UUID) ([]database.JobDependency, error) {
	var deps []database.JobDependency
	err := r.db.WithContext(ctx).Where("child_job_id = ?", jobID).Find(&deps).Error
	return deps, err
}

func (r *jobRepository) AreDependenciesCompleted(ctx context.Context, jobID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("job_dependencies").
		Joins("JOIN jobs ON jobs.id = job_dependencies.parent_job_id").
		Where("job_dependencies.child_job_id = ? AND jobs.status != ?", jobID, database.JobStatusCompleted).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
