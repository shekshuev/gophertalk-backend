package models

import "time"

type CreatePostDTO struct {
	Text   string `json:"text" validate:"required,min=0,max=280"`
	UserID uint64
}

type ReadPostUserDTO struct {
	ID        uint64 `json:"id"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ReadPostDTO struct {
	ID         uint64           `json:"id"`
	Text       string           `json:"text"`
	RepostOfID *uint64          `json:"repost_of_id,omitempty"`
	User       *ReadPostUserDTO `json:"user,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	LikesCount uint             `json:"likes_count"`
	ViewsCount uint             `json:"views_count"`
}
