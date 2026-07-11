package category

import (
	"context"
	"errors"
	"fmt"
	"server/internal/application/auth"
	"server/internal/domain"

	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepository CategoryRepository
	userRepository     auth.UserRepository
}

func NewCategoryService(categoryRepository CategoryRepository, userRepository auth.UserRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
		userRepository:     userRepository,
	}
}

func (cs *CategoryService) GetCategoryListByID(ctx context.Context, username string) ([]domain.Category, error) {
	user, err := cs.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("User not found")
		}
	}

	res, err := cs.categoryRepository.GetCategoryListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cs *CategoryService) CreateCategory(ctx context.Context, username string, nameCate string, descCate string) (*CreateCateRes, error) {
	user, err := cs.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("User not found")
		}
	}

	category, err := domain.NewCategory(user.ID, nameCate, descCate)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("Ban da co folder nay roi")
		}
		return nil, err
	}

	res, err := cs.categoryRepository.CreateCategory(ctx, *category)
	if err != nil {
		return nil, err
	}
	return &CreateCateRes{
		res.ID,
		res.Name,
		res.Desc,
	}, nil
}
