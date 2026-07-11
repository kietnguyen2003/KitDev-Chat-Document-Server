package document

import (
	"context"
	"fmt"
	"server/internal/application/auth"
	"server/internal/application/category"
	docuemntindexer "server/internal/application/docuemnt_indexer"
	"server/internal/application/storage"
	"server/internal/domain"
)

type DocumentService struct {
	userRepo        auth.UserRepository
	cateRepo        category.CategoryRepository
	documentStoRepo MinIOStorageDocument
	storageRepo     storage.StorageRepository
	ragRepo         docuemntindexer.DocumentIndexer
	documentRepo    DocumentRepository
}

func NewDocumentService(
	documentStoRepo MinIOStorageDocument,
	userRepo auth.UserRepository,
	cateRepo category.CategoryRepository,
	storageRepo storage.StorageRepository,
	ragRepo docuemntindexer.DocumentIndexer,
	documentRepo DocumentRepository,

) *DocumentService {
	return &DocumentService{
		documentStoRepo: documentStoRepo,
		userRepo:        userRepo,
		cateRepo:        cateRepo,
		storageRepo:     storageRepo,
		ragRepo:         ragRepo,
		documentRepo:    documentRepo,
	}
}

func (ds *DocumentService) UploadFile(ctx context.Context, req UploadDocDTO) (string, error) {
	// lay user
	user, err := ds.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("User dont exist")
	}
	// kiem tra user co category chua
	cateID, err := ds.cateRepo.CheckExist(ctx, user.ID, req.Category)
	if err != nil {
		return "", err
	}
	if cateID == 0 {
		return "", fmt.Errorf("Cate doesnt exist")
	}
	// Tru vao storage
	uploaded, err := ds.storageRepo.AddStorage(ctx, user.ID, uint64(req.Size))
	if err != nil {
		return "", err
	}
	if uploaded == false {
		return "", err
	}

	// upload vo minIO
	objectName := fmt.Sprintf("users/%d/categories/%s/documents/%s", user.ID, req.Category, req.DocName)

	key, err := ds.documentStoRepo.UploadReader(ctx, req.Bucket, objectName, req.Reader, req.Size, req.ContentType)
	if err != nil {
		_, err := ds.storageRepo.SubStorage(ctx, user.ID, uint64(req.Size))
		return "", err
	}

	fmt.Println("Content-Type: ", req.ContentType)

	docID, err := ds.documentRepo.CreateDocument(ctx, domain.Document{
		UserID:     user.ID,
		ObjectName: objectName,
		Name:       req.DocName,
		Size:       req.Size,
		CateID:     cateID,
	})
	if err != nil {
		return "", err
	}

	go func() {
		err := ds.ragRepo.IndexDocument(
			context.Background(), docuemntindexer.IndexDocumentRequest{
				DocumentID: docID,
				Bucket:     "documents",
				Object:     objectName,
				UserID:     user.ID,
				Category:   req.Category,
				Type:       req.ContentType,
			},
		)
		if err != nil {
			fmt.Println("Error in request to rag: ", err)
		}
	}()

	return key, nil
}

func (ds *DocumentService) UpdateStatus(ctx context.Context, username string, docID uint, status string) error {
	user, err := ds.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("User not found")
	}

	err = ds.documentRepo.UpdateStatus(ctx, docID, status)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DocumentService) GetDocumentList(ctx context.Context, username, category string) ([]domain.Document, error) {
	user, err := ds.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("User not found")
	}

	cateID, err := ds.cateRepo.CheckExist(ctx, user.ID, category)
	if err != nil {
		return nil, err
	}

	if cateID <= 0 {
		return nil, fmt.Errorf("Category not found in this user")
	}

	res, err := ds.documentRepo.GetDocumentListByUserAndCatID(ctx, user.ID, cateID)

	return res, err
}

func (ds *DocumentService) DeleteDocument(ctx context.Context, username, object string, size uint64) (bool, error) {
	user, err := ds.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, fmt.Errorf("User doesnt exist")
	}
	// tru storage
	subStorage, err := ds.storageRepo.SubStorage(ctx, user.ID, size)
	if err != nil || subStorage == false {
		return false, err
	}
	// xoa document trong minio
	err = ds.documentStoRepo.DeleteDocument(ctx, object)
	if err != nil {
		ds.storageRepo.AddStorage(ctx, user.ID, size)
		return false, err
	}

	// xoa document trong gorm
	docID, err := ds.documentRepo.DeletDocument(ctx, object)
	if err != nil {
		ds.storageRepo.AddStorage(ctx, user.ID, size)
		return false, err
	}

	// delete trong server rag
	err = ds.ragRepo.DeleteDocument(ctx, docuemntindexer.DeleteDocumentRequest{DocumentID: docID})
	if err != nil {
		return false, err
	}
	return true, nil
}
