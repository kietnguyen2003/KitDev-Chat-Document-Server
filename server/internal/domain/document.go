package domain

type Document struct {
	ID         uint
	UserID     uint
	CateID     uint
	Name       string
	ObjectName string
	CreatedAt  string
	UpdatedAt  string
	Size       int64
	Status     string
}

type FileInfo struct {
	ObjectKey string
	Size      int64
}
