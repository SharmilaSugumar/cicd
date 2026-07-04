package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QueueRepository interface {
	Create(ctx context.Context, queue *database.Queue) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.Queue, error)
	Update(ctx context.Context, queue *database.Queue) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.Queue, error)
	Pause(ctx context.Context, id uuid.UUID) error
	Resume(ctx context.Context, id uuid.UUID) error
}

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) Create(ctx context.Context, queue *database.Queue) error {
	return r.db.WithContext(ctx).Create(queue).Error
}

func (r *queueRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Queue, error) {
	var queue database.Queue
	err := r.db.WithContext(ctx).First(&queue, "id = ?", id).Error
	return &queue, err
}

func (r *queueRepository) Update(ctx context.Context, queue *database.Queue) error {
	return r.db.WithContext(ctx).Save(queue).Error
}

func (r *queueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.Queue{}, "id = ?", id).Error
}

func (r *queueRepository) List(ctx context.Context, limit, offset int) ([]database.Queue, error) {
	var queues []database.Queue
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&queues).Error
	return queues, err
}

func (r *queueRepository) Pause(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&database.Queue{}).Where("id = ?", id).Update("status", database.QueueStatusPaused).Error
}

func (r *queueRepository) Resume(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&database.Queue{}).Where("id = ?", id).Update("status", database.QueueStatusActive).Error
}
