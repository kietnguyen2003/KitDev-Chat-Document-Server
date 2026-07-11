package storage

import (
	"context"
	"server/internal/domain"
)

type StorageRepository interface {
	Create(ctx context.Context, storage domain.UserStorage) (*domain.UserStorage, error)
	GetStorage(ctx context.Context, userID uint) (*domain.UserStorage, error)
	AddStorage(ctx context.Context, userID uint, size uint64) (bool, error)
	SubStorage(ctx context.Context, userID uint, size uint64) (bool, error)
}
