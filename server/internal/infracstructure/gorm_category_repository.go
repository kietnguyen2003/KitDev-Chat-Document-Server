package infracstructure

import (
	"context"
	"errors"
	"server/internal/domain"
	"time"

	"gorm.io/gorm"
)

type GormCategory struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	UserID uint   `gorm:"not null;uniqueIndex:idx_user_category_name"`
	Name   string `gorm:"not null;uniqueIndex:idx_user_category_name"`
	Desc   string `gorm:"not null"`

	User GormUser `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GormCategory) TableName() string {
	return "categories"
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) GetCategoryListByUserID(ctx context.Context, userID uint) ([]domain.Category, error) {
	var gormCategoryList []GormCategory

	err := cr.db.WithContext(ctx).Where("user_id = ?", userID).Find(&gormCategoryList).Error
	if err != nil {
		return nil, err
	}
	var categoryList []domain.Category
	for _, category := range gormCategoryList {
		categoryList = append(categoryList, *gormCategoryToDomain(category))
	}
	return categoryList, nil
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, category domain.Category) (*domain.Category, error) {
	gormCategory := GormCategory{
		UserID:    category.UserID,
		Name:      category.Name,
		Desc:      category.Desc,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdateAt,
	}

	err := cr.db.WithContext(ctx).Create(&gormCategory).Error
	if err != nil {
		return nil, err
	}

	return gormCategoryToDomain(gormCategory), nil
}

func (cr *CategoryRepository) CheckExist(
	ctx context.Context,
	userID uint,
	categoryName string,
) (uint, error) {
	var category GormCategory

	err := cr.db.WithContext(ctx).
		Where("user_id = ? AND name = ?", userID, categoryName).
		First(&category).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return category.ID, nil
}

func gormCategoryToDomain(category GormCategory) *domain.Category {
	return &domain.Category{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdateAt:  category.UpdatedAt,
	}
}
