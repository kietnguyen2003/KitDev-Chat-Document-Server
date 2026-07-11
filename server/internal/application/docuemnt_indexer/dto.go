package docuemntindexer

type IndexDocumentRequest struct {
	DocumentID uint   `json:"document_id"`
	Bucket     string `json:"bucket"`
	Object     string `json:"object"`
	UserID     uint   `json:"user_id"`
	Category   string `json:"category"`
	Type       string `json:"type"`
}

type DeleteDocumentRequest struct {
	DocumentID uint `json:"document_id"`
}
