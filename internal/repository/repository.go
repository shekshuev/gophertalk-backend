package repository

import (
	"fmt"

	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type DatabaseChecker interface {
	CheckDBConnection() error
}

type UserRepository interface {
	GetAllUsers(limit, offset uint64) ([]models.ReadUserDTO, error)
	GetUserByID(id uint64) (*models.ReadUserDTO, error)
	GetUserByUserName(userName string) (*models.ReadAuthUserDataDTO, error)
	CreateUser(user models.CreateUserDTO) (*models.ReadAuthUserDataDTO, error)
	UpdateUser(id uint64, user models.UpdateUserDTO) (*models.ReadUserDTO, error)
	DeleteUser(id uint64) error
}

type PostRepository interface {
	GetAllPosts(dto models.FilterPostDTO) ([]models.ReadPostDTO, error)
	CreatePost(post models.CreatePostDTO) (*models.ReadPostDTO, error)
	DeletePost(id, ownerID uint64) error
	ViewPost(id, viewedByID uint64) error
	LikePost(id, likedByID uint64) error
	DislikePost(id, dislikedByID uint64) error
}

var ErrNotFound = fmt.Errorf("not found")
var ErrUserExists = fmt.Errorf("user with same name already exists")
var ErrNoFieldsToUpdate = fmt.Errorf("no fields to update")
var ErrAlreadyLiked = fmt.Errorf("already liked")
var ErrAlreadyViewed = fmt.Errorf("already viewed")
