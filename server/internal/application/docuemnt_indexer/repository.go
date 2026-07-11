package docuemntindexer

import "context"

type DocumentIndexer interface {
	IndexDocument(ctx context.Context, req IndexDocumentRequest) error
	DeleteDocument(ctx context.Context, id DeleteDocumentRequest) error
}
