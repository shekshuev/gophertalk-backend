package service

import (
	"fmt"

	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type UserService interface {
	GetAllUsers(limit, offset int) ([]models.ReadUserDTO, error)
	GetUserByID(id int) (*models.ReadUserDTO, error)
	UpdateUser(id int, user models.UpdateUserDTO) (*models.ReadUserDTO, error)
	DeleteUser(id int) error
}

type AuthService interface {
	Login(dto models.LoginUserDTO) (*models.ReadTokenDTO, error)
	Register(dto models.RegisterUserDTO) (*models.ReadTokenDTO, error)
}

var ErrUserNotFound = fmt.Errorf("user not found")
var ErrWrongPassword = fmt.Errorf("wrong password")
