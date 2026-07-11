package infracstructure

import (
	"context"
	"fmt"
	"server/internal/domain"

	"gorm.io/gorm"
)

// định dạng Gorm
type GormUser struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	Username       string
	HashedPassword string
	Fullname       string
	Role           string
	CreatedAt      int64
	UpdatedAt      int64
}

// Định dạng tên Table
func (GormUser) TableName() string {
	return "users"
}

// khai báo repository implementation
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user GormUser
	err := ur.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return gormToDomainUser(user), nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	gormUser := GormUser{
		Username:       user.Username,
		Role:           user.Role,
		HashedPassword: user.HashedPassword,
		Fullname:       user.Fullname,
		CreatedAt:      user.CreatedAt.Unix(),
		UpdatedAt:      user.UpdateAt.Unix(),
	}

	err := ur.db.WithContext(ctx).Create(&gormUser).Error
	if err != nil {
		fmt.Println("Create fail with err", err)
		return err
	}
	user.ID = gormUser.ID
	return nil
}

func gormToDomainUser(user GormUser) *domain.User {
	return &domain.User{
		ID:             user.ID,
		Username:       user.Username,
		Fullname:       user.Fullname,
		Role:           user.Role,
		HashedPassword: user.HashedPassword,
	}
}
