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

func NewUserServiceImpl(repo repository.UserRepository, cfg *config.Config) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, cfg: cfg}
}

func (s *UserServiceImpl) GetAllUsers(limit, offset uint64) ([]models.ReadUserDTO, error) {
	return s.repo.GetAllUsers(limit, offset)
}

func (s *UserServiceImpl) GetUserByID(id uint64) (*models.ReadUserDTO, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserServiceImpl) UpdateUser(id uint64, user models.UpdateUserDTO) (*models.ReadUserDTO, error) {
	if user.Password != "" {
		user.PasswordHash = utils.HashPassword(user.Password)
	}
	return s.repo.UpdateUser(id, user)
}

func (s *UserServiceImpl) DeleteUser(id uint64) error {
	return s.repo.DeleteUser(id)
}
