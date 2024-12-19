package service

import (
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/repository"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
)

type UserServiceImpl struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserServiceImpl(repo *repository.UserRepository, cfg *config.Config) *UserServiceImpl {
	return &UserServiceImpl{repo: *repo, cfg: cfg}
}

func (s *UserServiceImpl) GetAllUsers() ([]models.ReadUserDTO, error) {
	return s.repo.GetAllUsers()
}

func (s *UserServiceImpl) GetUserByID(id int) (*models.ReadUserDTO, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserServiceImpl) UpdateUser(id int, user models.UpdateUserDTO) (*models.ReadUserDTO, error) {
	if user.Password != "" {
		if user.Password != user.PasswordConfirm {
			return nil, ErrPasswordMismatch
		}
		user.PasswordHash = utils.HashPassword(user.Password)
	}
	return s.repo.UpdateUser(id, user)
}

func (s *UserServiceImpl) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
