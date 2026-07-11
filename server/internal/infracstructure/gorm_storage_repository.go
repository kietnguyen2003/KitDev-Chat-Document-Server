package infracstructure

import (
	"context"
	"fmt"
	"server/internal/domain"
	"time"

	"gorm.io/gorm"
)

type GormStorage struct {
	ID     uint `gorm:"primaryKey;autoIncrement"`
	UserID uint `gorm:"not null;uniqueIndex"`

	User GormUser `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CurrentSize uint64
	LimitSize   uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (GormStorage) TableName() string {
	return "storages"
}

type StorageRepository struct {
	db *gorm.DB
}

func NewStorageRepository(db *gorm.DB) *StorageRepository {
	return &StorageRepository{db: db}
}

func (sr *StorageRepository) Create(ctx context.Context, storage domain.UserStorage) (*domain.UserStorage, error) {
	gormStorage := &GormStorage{
		UserID:      storage.UserID,
		LimitSize:   storage.LimitSize,
		CurrentSize: storage.CurrenSize,
		CreatedAt:   storage.CreatedAt,
		UpdatedAt:   storage.UpdatedAt,
	}
	if err := sr.db.Create(gormStorage).Error; err != nil {
		return nil, err
	}
	return &domain.UserStorage{
		ID:         gormStorage.ID,
		UserID:     gormStorage.UserID,
		CurrenSize: gormStorage.CurrentSize,
		LimitSize:  gormStorage.LimitSize,
		CreatedAt:  gormStorage.CreatedAt,
		UpdatedAt:  gormStorage.UpdatedAt,
	}, nil
}

func (sr *StorageRepository) GetStorage(ctx context.Context, userID uint) (*domain.UserStorage, error) {
	var gormStorage GormStorage

	if err := sr.db.Where("user_id = ?", userID).First(&gormStorage).Error; err != nil {
		return nil, err
	}

	return &domain.UserStorage{
		ID:         gormStorage.ID,
		UserID:     gormStorage.UserID,
		CurrenSize: gormStorage.CurrentSize,
		LimitSize:  gormStorage.LimitSize,
		CreatedAt:  gormStorage.CreatedAt,
		UpdatedAt:  gormStorage.UpdatedAt,
	}, nil
}

func (sr *StorageRepository) AddStorage(ctx context.Context, userID uint, size uint64) (bool, error) {
	var gormStorage GormStorage
	if err := sr.db.Where("user_id = ?", userID).First(&gormStorage).Error; err != nil {
		return false, err
	}
	newSize := gormStorage.CurrentSize + size
	if newSize > gormStorage.LimitSize {
		return false, fmt.Errorf("Khong du dung luong")
	}
	if err := sr.db.Model(&gormStorage).Where("user_id = ?", userID).Update("current_size", newSize).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (sr *StorageRepository) SubStorage(ctx context.Context, userID uint, size uint64) (bool, error) {
	var gormStorage GormStorage
	if err := sr.db.Where("user_id = ?", userID).First(&gormStorage).Error; err != nil {
		return false, err
	}
	if gormStorage.CurrentSize < size {
		return false, fmt.Errorf("invalid storage size")
	}
	newSize := gormStorage.CurrentSize - size

	if err := sr.db.Model(&gormStorage).Where("user_id = ?", userID).Update("current_size", newSize).Error; err != nil {
		return false, err
	}
	return true, nil
}
