package models

import "time"

type CreatePostDTO struct {
	Text   string `json:"text" validate:"required,min=0,max=280"`
	UserID uint64
}

type ReadPostDTO struct {
	ID         uint64    `json:"id"`
	Text       string    `json:"text"`
	RepostOfID *uint64   `json:"repost_of_id,omitempty"`
	UserID     uint64    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}
