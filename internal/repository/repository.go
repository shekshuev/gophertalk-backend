package repository

import (
	"fmt"

	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type DatabaseChecker interface {
	CheckDBConnection() error
}

type UserRepository interface {
	GetAllUsers() []models.ReadUserDTO
	GetUserByID(id int) (*models.ReadUserDTO, error)
	RegisterUser(user models.RegisterUserDTO) (*models.ReadUserDTO, error)
	UpdateUser(id int, user models.UpdateUserDTO) (*models.ReadUserDTO, error)
	DeleteUser(id int) error
}

var ErrNotFound = fmt.Errorf("not found")
var ErrUserExists = fmt.Errorf("user with same name already exists")
var ErrNoFieldsToUpdate = fmt.Errorf("no fields to update")
var ErrPasswordMismatch = fmt.Errorf("password mismatch")
