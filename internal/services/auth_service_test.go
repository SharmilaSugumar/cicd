package services

import (
	"context"
	"errors"
	"forgeflow/internal/database"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	GetByEmailFunc func(ctx context.Context, email string) (*database.User, error)
	CreateFunc     func(ctx context.Context, user *database.User) error
	GetByIDFunc    func(ctx context.Context, id uuid.UUID) (*database.User, error)
	UpdateFunc     func(ctx context.Context, user *database.User) error
	DeleteFunc     func(ctx context.Context, id uuid.UUID) error
	ListFunc       func(ctx context.Context, limit, offset int) ([]database.User, error)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*database.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}
func (m *mockUserRepository) Create(ctx context.Context, user *database.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}
func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.User, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *mockUserRepository) Update(ctx context.Context, user *database.User) error {
	return m.UpdateFunc(ctx, user)
}
func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}
func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]database.User, error) {
	return m.ListFunc(ctx, limit, offset)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*database.User, error) {
			return nil, errors.New("not found") // indicates no duplicate
		},
		CreateFunc: func(ctx context.Context, user *database.User) error {
			user.ID = uuid.New()
			return nil
		},
	}
	service := NewAuthService(mockRepo)

	user, err := service.RegisterUser(context.Background(), "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be returned")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %v", user.Email)
	}
}

func TestRegisterUser_Duplicate(t *testing.T) {
	mockRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*database.User, error) {
			return &database.User{}, nil // user found
		},
	}
	service := NewAuthService(mockRepo)

	_, err := service.RegisterUser(context.Background(), "test@example.com", "password123", "Test User")
	if err != ErrDuplicateResource {
		t.Fatalf("expected ErrDuplicateResource, got %v", err)
	}
}

func TestLoginUser_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*database.User, error) {
			return &database.User{
				BaseModel:    database.BaseModel{ID: uuid.New()},
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}
	service := NewAuthService(mockRepo)

	token, err := service.LoginUser(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a valid token")
	}
}

func TestJWT_GenerateAndValidate(t *testing.T) {
	service := NewAuthService(nil) // repo not needed for JWT

	user := &database.User{
		BaseModel: database.BaseModel{ID: uuid.New()},
	}

	token, err := service.GenerateJWT(user)
	if err != nil {
		t.Fatalf("expected no error generating JWT, got %v", err)
	}

	claims, err := service.ValidateJWT(token)
	if err != nil {
		t.Fatalf("expected no error validating JWT, got %v", err)
	}
	if claims.Subject != user.ID.String() {
		t.Errorf("expected subject %s, got %s", user.ID.String(), claims.Subject)
	}
}
