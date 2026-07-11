package infracstructure

import (
	"context"
	"errors"
	"fmt"
	"server/internal/domain"
	"time"

	"gorm.io/gorm"
)

type GormDocument struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	UserID     uint
	CategoryID uint
	Name       string
	Object     string `gorm:"primaryKey;autoIncrement"`
	Size       int64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (GormDocument) TableName() string {
	return "documents"
}

type GormDocumentRepository struct {
	db *gorm.DB
}

func NewGormDocumentRepository(db *gorm.DB) *GormDocumentRepository {
	return &GormDocumentRepository{db: db}
}

func (gd *GormDocumentRepository) CreateDocument(ctx context.Context, req domain.Document) (uint, error) {
	gormDocument := &GormDocument{
		UserID:     req.UserID,
		CategoryID: req.CateID,
		Object:     req.ObjectName,
		Name:       req.Name,
		Size:       req.Size,
		Status:     "indexing",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := gd.db.WithContext(ctx).Create(gormDocument).Error; err != nil {
		return 0, err
	}

	return gormDocument.ID, nil
}

func (gd *GormDocumentRepository) UpdateStatus(ctx context.Context, docID uint, status string) error {
	var gormDocument GormDocument

	if err := gd.db.WithContext(ctx).Where("id = ?", docID).First(&gormDocument).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("Document not found")
		}
		return err
	}
	gormDocument.Status = status

	if err := gd.db.WithContext(ctx).Model(&gormDocument).Where("id = ?", docID).Update("status", status).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (gd *GormDocumentRepository) GetDocumentListByUserAndCatID(ctx context.Context, userID, categoryID uint) ([]domain.Document, error) {
	var gormDocumentList []GormDocument
	if err := gd.db.WithContext(ctx).Where("user_id = ?", userID).Where("category_id = ?", categoryID).Find(&gormDocumentList).Error; err != nil {
		return nil, err
	}
	var domainDocumentList []domain.Document

	for _, doc := range gormDocumentList {
		document := domain.Document{
			ID:         doc.ID,
			UserID:     userID,
			Name:       doc.Name,
			ObjectName: doc.Object,
			Size:       doc.Size,
			Status:     doc.Status,
		}

		domainDocumentList = append(domainDocumentList, document)
	}
	return domainDocumentList, nil
}

func (gd *GormDocumentRepository) DeletDocument(ctx context.Context, object string) (uint, error) {
	var gormDocument GormDocument

	if err := gd.db.WithContext(ctx).Where("object = ?", object).Delete(&gormDocument).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return 0, fmt.Errorf("Document not found")
		}
		return 0, err
	}
	return gormDocument.ID, nil
}
