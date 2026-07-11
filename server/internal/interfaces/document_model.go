package interfaces

type UploadFileRes struct {
	ObjectName string `json:"object_name"`
}

type GetDocumentListRes struct {
	Doc_ID     uint   `json:"document_id"`
	UserID     uint   `json:"user_id"`
	Name       string `json:"document_name"`
	ObjectName string `json:"object_name"`
	Size       int64  `json:"size"`
	Status     string `jsson:"status"`
}
