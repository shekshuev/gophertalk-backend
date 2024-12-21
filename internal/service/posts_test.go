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

func TestPostServiceImpl_GetAllPosts(t *testing.T) {
	testCases := []struct {
		name     string
		readDTOs []models.ReadPostDTO
		hasError bool
	}{
		{
			name: "Success get all posts",
			readDTOs: []models.ReadPostDTO{
				{
					ID:   1,
					Text: "Lorem ipsum dolor sit amet, consectetur adipiscing",
					User: &models.ReadPostUserDTO{
						ID:        1,
						UserName:  "username",
						FirstName: "first_name",
						LastName:  "last_name",
					},
					RepostOfID: nil,
					CreatedAt:  time.Now(),
				},
				{
					ID:   2,
					Text: "Lorem ipsum dolor sit amet, consectetur adipiscing",
					User: &models.ReadPostUserDTO{
						ID:        1,
						UserName:  "username",
						FirstName: "first_name",
						LastName:  "last_name",
					},
					RepostOfID: nil,
					CreatedAt:  time.Now(),
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
	m := mocks.NewMockPostRepository(ctrl)
	s := &PostServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().GetAllPosts(uint64(100), uint64(0)).Return(tc.readDTOs, nil)
			} else {
				m.EXPECT().GetAllPosts(uint64(100), uint64(0)).Return(nil, sql.ErrNoRows)
			}
			users, err := s.GetAllPosts(uint64(100), uint64(0))
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTOs, users, "Posts mismatch")
			}
		})
	}
}

func TestPostServiceImpl_GetPostByID(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint64
		readDTO  *models.ReadPostDTO
		hasError bool
	}{
		{
			name: "Success get post by ID",
			id:   1,
			readDTO: &models.ReadPostDTO{
				ID:   1,
				Text: "Lorem ipsum dolor sit amet, consectetur adipiscing",
				User: &models.ReadPostUserDTO{
					ID:        1,
					UserName:  "username",
					FirstName: "first_name",
					LastName:  "last_name",
				},
				RepostOfID: nil,
				CreatedAt:  time.Now(),
			},
			hasError: false,
		},
		{
			name:     "Post not found",
			id:       2,
			readDTO:  nil,
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockPostRepository(ctrl)
	s := &PostServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().GetPostByID(tc.id).Return(tc.readDTO, nil)
			} else {
				m.EXPECT().GetPostByID(tc.id).Return(nil, sql.ErrNoRows)
			}

			user, err := s.GetPostByID(tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "Post should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "Post mismatch")
			}

		})
	}
}

func TestPostServiceImpl_CreatePost(t *testing.T) {
	testCases := []struct {
		name      string
		createDTO models.CreatePostDTO
		readDTO   *models.ReadPostDTO
		hasError  bool
	}{
		{
			name: "Success create",
			createDTO: models.CreatePostDTO{
				Text:   "Lorem ipsum dolor sit amet, consectetur adipiscing",
				UserID: 1,
			},
			readDTO: &models.ReadPostDTO{
				ID:         1,
				Text:       "Lorem ipsum dolor sit amet, consectetur adipiscing",
				RepostOfID: nil,
				CreatedAt:  time.Now(),
			},
			hasError: false,
		},
		{
			name: "Error on insert SQL",
			createDTO: models.CreatePostDTO{
				Text:   "Lorem ipsum dolor sit amet, consectetur adipiscing",
				UserID: 1,
			},
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockPostRepository(ctrl)
	s := &PostServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().CreatePost(tc.createDTO).Return(tc.readDTO, nil)
			} else {
				m.EXPECT().CreatePost(tc.createDTO).Return(nil, sql.ErrNoRows)
			}

			user, err := s.CreatePost(tc.createDTO)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "Post should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "Post mismatch")
			}
		})
	}
}

func TestPostServiceImpl_DeletePost(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint64
		hasError bool
	}{
		{
			name:     "Success delete post",
			id:       1,
			hasError: false,
		},
		{
			name:     "Error on delete",
			id:       2,
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockPostRepository(ctrl)
	s := &PostServiceImpl{cfg: &cfg, repo: m}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				m.EXPECT().DeletePost(tc.id, uint64(0)).Return(nil)
			} else {
				m.EXPECT().DeletePost(tc.id, uint64(0)).Return(sql.ErrNoRows)
			}
			err := s.DeletePost(tc.id, uint64(0))
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}
		})
	}
}
