package domain

import "time"

type UserStorage struct {
	ID         uint
	UserID     uint
	CurrenSize uint64
	LimitSize  uint64
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

func NewStorage(userID uint) *UserStorage {
	if userID == 0 {
		return nil
	}
	return &UserStorage{
		UserID:     userID,
		CurrenSize: 0,
		LimitSize:  5 * 1024 * 1024,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
}
