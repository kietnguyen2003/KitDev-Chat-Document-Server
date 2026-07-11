package auth

import (
	"context"
	"server/internal/domain"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}
