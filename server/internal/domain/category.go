package domain

import (
	"fmt"
	"time"
)

type Category struct {
	ID        uint
	UserID    uint
	Name      string
	Desc      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func NewCategory(userID uint, name string, desc string) (*Category, error) {
	if name == "" || desc == "" {
		return nil, fmt.Errorf("Invalid Creadential")
	}
	return &Category{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}, nil
}
