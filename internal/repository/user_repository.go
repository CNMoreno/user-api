package repository

import (
	"context"

	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
)

// UserRepository interface of user in BD.
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updateFields *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}
