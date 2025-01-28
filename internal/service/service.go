package service

import (
	"fmt"

	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type UserService interface {
	GetAllUsers(limit, offset uint64) ([]models.ReadUserDTO, error)
	GetUserByID(id uint64) (*models.ReadUserDTO, error)
	UpdateUser(id uint64, user models.UpdateUserDTO) (*models.ReadUserDTO, error)
	DeleteUser(id uint64) error
}

type AuthService interface {
	Login(dto models.LoginUserDTO) (*models.ReadTokenDTO, error)
	Register(dto models.RegisterUserDTO) (*models.ReadTokenDTO, error)
}

type PostService interface {
	GetAllPosts(dto models.FilterPostDTO) ([]models.ReadPostDTO, error)
	CreatePost(post models.CreatePostDTO) (*models.ReadPostDTO, error)
	DeletePost(id, ownerID uint64) error
	ViewPost(id, viewedByID uint64) error
	LikePost(id, likedByID uint64) error
	DislikePost(id, dislikedByID uint64) error
}

var ErrUserNotFound = fmt.Errorf("user not found")
var ErrWrongPassword = fmt.Errorf("wrong password")
