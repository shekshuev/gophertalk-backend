package service

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/mocks"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthServiceImpl_Login(t *testing.T) {
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	authService := &AuthServiceImpl{repo: repo, cfg: &cfg}

	testCases := []struct {
		name     string
		dto      models.LoginUserDTO
		hasError bool
		mockSet  func()
	}{
		{
			name: "Success",
			dto: models.LoginUserDTO{
				UserName: "testuser",
				Password: "password123",
			},
			hasError: false,
			mockSet: func() {
				repo.EXPECT().GetUserByUserName("testuser").Return(&models.ReadAuthUserDataDTO{
					ID:           1,
					UserName:     "testuser",
					Status:       models.StatusActive,
					PasswordHash: utils.HashPassword("password123"),
				}, nil)
			},
		},
		{
			name: "User not found",
			dto: models.LoginUserDTO{
				UserName: "nonexistentuser",
				Password: "password123",
			},
			hasError: true,
			mockSet: func() {
				repo.EXPECT().GetUserByUserName("nonexistentuser").Return(nil, sql.ErrNoRows)
			},
		},
		{
			name: "Wrong password",
			dto: models.LoginUserDTO{
				UserName: "testuser",
				Password: "wrongpassword",
			},
			hasError: true,
			mockSet: func() {
				repo.EXPECT().GetUserByUserName("testuser").Return(&models.ReadAuthUserDataDTO{
					ID:           1,
					UserName:     "testuser",
					PasswordHash: utils.HashPassword("password123"),
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSet()
			_, err := authService.Login(tc.dto)
			if tc.hasError {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestAuthServiceImpl_Register(t *testing.T) {
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	authService := &AuthServiceImpl{repo: repo, cfg: &cfg}

	fixedPasswordHash := "$2a$10$CmIxNqxCFrgFoji4qyka0.UvTV4wG54LN5UJjV7mfH6q0caiNGUvK"

	testCases := []struct {
		name        string
		dto         models.RegisterUserDTO
		expectedErr bool
		mockSet     func()
	}{
		{
			name: "Success",
			dto: models.RegisterUserDTO{
				UserName:        "testuser",
				Password:        "password123",
				PasswordConfirm: "password123",
				FirstName:       "Test",
				LastName:        "User",
			},
			expectedErr: false,
			mockSet: func() {
				repo.EXPECT().CreateUser(gomock.Any()).Return(&models.ReadAuthUserDataDTO{
					ID:           1,
					UserName:     "testuser",
					Status:       models.StatusActive,
					PasswordHash: fixedPasswordHash,
				}, nil)
			},
		},
		{
			name: "Password mismatch",
			dto: models.RegisterUserDTO{
				UserName:        "testuser",
				Password:        "password123",
				PasswordConfirm: "password124",
				FirstName:       "Test",
				LastName:        "User",
			},
			expectedErr: true,
			mockSet:     func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSet()
			_, err := authService.Register(tc.dto)
			if tc.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Expected no error but got one")
			}
		})
	}
}
