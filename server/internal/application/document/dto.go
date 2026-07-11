package document

import (
	"io"
)

type UploadDocDTO struct {
	Username    string
	Category    string
	Bucket      string
	DocName     string
	Reader      io.Reader
	Size        int64
	ContentType string
}
