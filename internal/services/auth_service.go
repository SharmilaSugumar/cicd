package services

import (
	"context"
	"forgeflow/internal/database"
	"forgeflow/internal/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password, name string) (*database.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	VerifyPassword(hash, password string) bool
	GenerateJWT(user *database.User) (string, error)
	ValidateJWT(tokenStr string) (*jwt.RegisteredClaims, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

func (s *authService) RegisterUser(ctx context.Context, email, password, name string) (*database.User, error) {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrDuplicateResource
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}

	user := &database.User{
		Email:        email,
		PasswordHash: string(hashed),
		Name:         name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, ErrInternalError
	}

	return user, nil
}

func (s *authService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrUnauthorized
	}

	if !s.VerifyPassword(user.PasswordHash, password) {
		return "", ErrUnauthorized
	}

	return s.GenerateJWT(user)
}

func (s *authService) VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) GenerateJWT(user *database.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrUnauthorized
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrUnauthorized
}
