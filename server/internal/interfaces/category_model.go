package interfaces

import "time"

type CreateCateReq struct {
	Name string `json:"name"`
	Des  string `json:"desc"`
}

type GetCateRes struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}
