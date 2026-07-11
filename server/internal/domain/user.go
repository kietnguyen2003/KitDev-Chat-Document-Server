package domain

import (
	"errors"
	"time"
)

type User struct {
	ID             uint
	Username       string
	HashedPassword string
	Fullname       string
	Role           string
	CreatedAt      time.Time
	UpdateAt       time.Time
}

func NewUser(username, hashedPassword, fullname string) (*User, error) {
	if username == "" || hashedPassword == "" || fullname == "" {
		return nil, errors.New("Invalid Credentials")
	}

	return &User{
		Username:       username,
		HashedPassword: hashedPassword,
		Fullname:       fullname,
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdateAt:       time.Now(),
	}, nil
}
