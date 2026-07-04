package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkerRepository interface {
	Create(ctx context.Context, worker *database.Worker) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.Worker, error)
	Update(ctx context.Context, worker *database.Worker) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.Worker, error)
	Register(ctx context.Context, worker *database.Worker) error
	UpdateHeartbeat(ctx context.Context, workerID uuid.UUID) error
	MarkOffline(ctx context.Context, id uuid.UUID) error
}

type workerRepository struct {
	db *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) WorkerRepository {
	return &workerRepository{db: db}
}

func (r *workerRepository) Create(ctx context.Context, worker *database.Worker) error {
	return r.db.WithContext(ctx).Create(worker).Error
}

func (r *workerRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Worker, error) {
	var worker database.Worker
	err := r.db.WithContext(ctx).First(&worker, "id = ?", id).Error
	return &worker, err
}

func (r *workerRepository) Update(ctx context.Context, worker *database.Worker) error {
	return r.db.WithContext(ctx).Save(worker).Error
}

func (r *workerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.Worker{}, "id = ?", id).Error
}

func (r *workerRepository) List(ctx context.Context, limit, offset int) ([]database.Worker, error) {
	var workers []database.Worker
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&workers).Error
	return workers, err
}

func (r *workerRepository) Register(ctx context.Context, worker *database.Worker) error {
	worker.Status = database.WorkerStatusOnline
	return r.db.WithContext(ctx).Create(worker).Error
}

func (r *workerRepository) UpdateHeartbeat(ctx context.Context, workerID uuid.UUID) error {
	heartbeat := database.WorkerHeartbeat{
		WorkerID: workerID,
	}
	return r.db.WithContext(ctx).Create(&heartbeat).Error
}

func (r *workerRepository) MarkOffline(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&database.Worker{}).Where("id = ?", id).Update("status", database.WorkerStatusOffline).Error
}
