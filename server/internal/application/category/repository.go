package category

import (
	"context"
	"server/internal/domain"
)

type CategoryRepository interface {
	GetCategoryListByUserID(ctx context.Context, userID uint) ([]domain.Category, error)
	CreateCategory(ctx context.Context, category domain.Category) (*domain.Category, error)

	CheckExist(ctx context.Context, userID uint, categoryName string) (uint, error)
}
