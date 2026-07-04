package repositories

import (
	"context"
	"forgeflow/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *database.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.User, error)
	GetByEmail(ctx context.Context, email string) (*database.User, error)
	Update(ctx context.Context, user *database.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]database.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *database.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.User, error) {
	var user database.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *database.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.User{}, "id = ?", id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]database.User, error) {
	var users []database.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
