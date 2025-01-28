package service

import (
	"strconv"

	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/repository"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
)

type AuthServiceImpl struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewAuthServiceImpl(repo repository.UserRepository, cfg *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repo, cfg: cfg}
}

func (s *AuthServiceImpl) Login(dto models.LoginUserDTO) (*models.ReadTokenDTO, error) {
	user, err := s.repo.GetUserByUserName(dto.UserName)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if !utils.VerifyPassword(dto.Password, user.PasswordHash) {
		return nil, ErrWrongPassword
	}
	return s.generateTokenPair(*user)
}

func (s *AuthServiceImpl) Register(dto models.RegisterUserDTO) (*models.ReadTokenDTO, error) {
	createDTO := models.CreateUserDTO{
		UserName:     dto.UserName,
		PasswordHash: utils.HashPassword(dto.Password),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
	}
	user, err := s.repo.CreateUser(createDTO)
	if err != nil {
		return nil, err
	}
	return s.generateTokenPair(*user)
}

func (s *AuthServiceImpl) generateTokenPair(user models.ReadAuthUserDataDTO) (*models.ReadTokenDTO, error) {
	accessToken, err := utils.CreateToken(
		s.cfg.AccessTokenSecret,
		strconv.FormatUint(user.ID, 10),
		s.cfg.AccessTokenExpires,
	)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.CreateToken(
		s.cfg.RefreshTokenSecret,
		strconv.FormatUint(user.ID, 10),
		s.cfg.RefreshTokenExpires,
	)
	if err != nil {
		return nil, err
	}
	return &models.ReadTokenDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
