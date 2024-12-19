package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/mocks"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryImpl_GetAllUsers(t *testing.T) {
	testCases := []struct {
		name     string
		readDTOs []models.ReadUserDTO
		hasError bool
	}{
		{
			name: "Success get all users",
			readDTOs: []models.ReadUserDTO{
				{
					ID:        1,
					UserName:  "john",
					FirstName: "John",
					LastName:  "Doe",
					Status:    1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        2,
					UserName:  "jane",
					FirstName: "Jane",
					LastName:  "Smith",
					Status:    1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			hasError: false,
		},
		{
			name:     "Error on SQL query",
			readDTOs: nil,
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockUserRepository(ctrl)
	s := &UserServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().GetAllUsers().Return(tc.readDTOs, nil)
			} else {
				m.EXPECT().GetAllUsers().Return(nil, sql.ErrNoRows)
			}
			users, err := s.GetAllUsers()
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTOs, users, "Users mismatch")
			}
		})
	}
}

func TestUserRepositoryImpl_GetUserByID(t *testing.T) {
	testCases := []struct {
		name     string
		id       int
		readDTO  *models.ReadUserDTO
		hasError bool
	}{
		{
			name: "Success get user by ID",
			id:   1,
			readDTO: &models.ReadUserDTO{
				ID:        1,
				UserName:  "john",
				FirstName: "John",
				LastName:  "Doe",
				Status:    1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			hasError: false,
		},
		{
			name:     "User not found",
			id:       2,
			readDTO:  nil,
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockUserRepository(ctrl)
	s := &UserServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().GetUserByID(tc.id).Return(tc.readDTO, nil)
			} else {
				m.EXPECT().GetUserByID(tc.id).Return(nil, sql.ErrNoRows)
			}

			user, err := s.GetUserByID(tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "User mismatch")
			}

		})
	}
}

func TestUserRepositoryImpl_UpdateUser(t *testing.T) {
	testCases := []struct {
		name      string
		id        int
		updateDTO models.UpdateUserDTO
		readDTO   *models.ReadUserDTO
		hasError  bool
	}{
		{
			name: "Success update user",
			id:   1,
			updateDTO: models.UpdateUserDTO{
				UserName:     "john_updated",
				FirstName:    "John",
				LastName:     "Doe",
				PasswordHash: "password",
			},
			readDTO: &models.ReadUserDTO{
				ID:        1,
				UserName:  "john_updated",
				FirstName: "John",
				LastName:  "Doe",
				Status:    1,
				CreatedAt: time.Now().Add(-time.Hour),
				UpdatedAt: time.Now(),
			},
			hasError: false,
		},
		{
			name: "Error on update SQL",
			id:   2,
			updateDTO: models.UpdateUserDTO{
				UserName: "john_updated",
			},
			readDTO:  nil,
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockUserRepository(ctrl)
	s := &UserServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().UpdateUser(tc.id, tc.updateDTO).Return(tc.readDTO, nil)
			} else {
				m.EXPECT().UpdateUser(tc.id, tc.updateDTO).Return(nil, sql.ErrNoRows)
			}

			user, err := s.UpdateUser(tc.id, tc.updateDTO)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "User mismatch")
			}
		})
	}
}

func TestUserRepositoryImpl_DeleteUser(t *testing.T) {
	testCases := []struct {
		name     string
		id       int
		hasError bool
	}{
		{
			name:     "Success delete user",
			id:       1,
			hasError: false,
		},
		{
			name:     "Error on delete SQL",
			id:       2,
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockUserRepository(ctrl)
	s := &UserServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().DeleteUser(tc.id).Return(nil)
			} else {
				m.EXPECT().DeleteUser(tc.id).Return(sql.ErrNoRows)
			}
			err := s.DeleteUser(tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}
		})
	}
}
