package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *database.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.Organization, error)
	Update(ctx context.Context, org *database.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.Organization, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *database.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Organization, error) {
	var org database.Organization
	err := r.db.WithContext(ctx).Preload("Members").First(&org, "id = ?", id).Error
	return &org, err
}

func (r *organizationRepository) Update(ctx context.Context, org *database.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.Organization{}, "id = ?", id).Error
}

func (r *organizationRepository) List(ctx context.Context, limit, offset int) ([]database.Organization, error) {
	var orgs []database.Organization
	err := r.db.WithContext(ctx).Preload("Members").Limit(limit).Offset(offset).Find(&orgs).Error
	return orgs, err
}
