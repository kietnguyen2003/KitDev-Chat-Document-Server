package document

import (
	"context"
	"io"
	"server/internal/domain"
)

type MinIOStorageDocument interface {
	UploadReader(ctx context.Context, bucket string, object string, reader io.Reader, size int64, contentType string) (string, error)
	DeleteDocument(ctx context.Context, object string) error
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, req domain.Document) (uint, error)
	DeletDocument(ctx context.Context, object string) (uint, error)

	UpdateStatus(ctx context.Context, docID uint, status string) error
	GetDocumentListByUserAndCatID(ctx context.Context, userID, categoryID uint) ([]domain.Document, error)
}
