package service

import (
	"fmt"

	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type UserService interface {
	GetAllUsers() []models.ReadUserDTO
	GetUserByID(id int) (*models.ReadUserDTO, error)
	UpdateUser(id int, user models.UpdateUserDTO) (*models.ReadUserDTO, error)
	DeleteUser(id int) error
}

var ErrPasswordMismatch = fmt.Errorf("password mismatch")
