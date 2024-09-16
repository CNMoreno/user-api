package usecase

import (
	"context"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/CNMoreno/cnm-proyect-go/internal/repository"
)

// UserService handles to obtain user repository.
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService obtain new user service.
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser interface for create user.
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	return s.userRepo.CreateUser(ctx, user)
}

// CreateUserBatch interface for create users.
func (s *UserService) CreateUserBatch(ctx context.Context, user *[]domain.User) ([]interface{}, error) {
	return s.userRepo.CreateUserBatch(ctx, user)
}

// GetUserByID interface for get user by ID.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// UpdateUser interface for update user by ID.
func (s *UserService) UpdateUser(ctx context.Context, id string, updateFields *domain.User) (*domain.User, error) {
	return s.userRepo.UpdateUser(ctx, id, updateFields)
}

// DeleteUser interface for delete user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.DeleteUser(ctx, id)
}
