package models

import "time"

type CreatePostDTO struct {
	Text      string  `json:"text" validate:"required,min=0,max=280"`
	ReplyToID *uint64 `json:"reply_to_id,omitempty" validate:"omitempty,gt=0"`
	UserID    uint64
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
	ReplyToID  *uint64          `json:"reply_to_id,omitempty"`
	User       *ReadPostUserDTO `json:"user,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	LikesCount uint             `json:"likes_count"`
	ViewsCount uint             `json:"views_count"`
	UserLiked  bool             `json:"user_liked"`
	UserViewed bool             `json:"user_viewed"`
}

type FilterPostDTO struct {
	Search    string `json:"search,omitempty"`
	UserID    uint64
	ReplyToID uint64 `json:"reply_to_id,omitempty"`
	Limit     uint64 `json:"limit" validate:"required,min=0,max=100"`
	Offset    uint64 `json:"offset" validate:"required,gte=0"`
}
