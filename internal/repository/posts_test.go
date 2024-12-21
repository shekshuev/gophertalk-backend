package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPostRepositoryImpl_CreatePost(t *testing.T) {
	testCases := []struct {
		name      string
		createDTO models.CreatePostDTO
		readDTO   models.ReadPostDTO
		hasError  bool
	}{
		{
			name: "Success create",
			createDTO: models.CreatePostDTO{
				Text:   "Lorem ipsum dolor sit amet, consectetur adipiscing",
				UserID: 1,
			},
			readDTO: models.ReadPostDTO{
				ID:         1,
				Text:       "Lorem ipsum dolor sit amet, consectetur adipiscing",
				RepostOfID: 0,
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()
	r := &PostRepositoryImpl{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into posts (text, user_id) values ($1, $2)
					returning id, text, user_id, created_at;
					`)).
					WithArgs(
						tc.createDTO.Text,
						tc.createDTO.UserID).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "text", "user_id", "created_at"},
						).AddRow(
							tc.readDTO.ID,
							tc.readDTO.Text,
							tc.readDTO.UserID,
							tc.readDTO.CreatedAt,
						),
					)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into posts (text, user_id) values ($1, $2)
					returning id, text, user_id, created_at;
					`)).
					WithArgs(
						tc.createDTO.Text,
						tc.createDTO.UserID).
					WillReturnError(sql.ErrNoRows)
			}
			post, err := r.CreatePost(tc.createDTO)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, *post, "Post mismatch")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostRepositoryImpl_GetAllPosts(t *testing.T) {
	testCases := []struct {
		name     string
		readDTOs []models.ReadPostDTO
		hasError bool
	}{
		{
			name: "Success get all posts",
			readDTOs: []models.ReadPostDTO{
				{
					ID:         1,
					Text:       "Lorem ipsum dolor sit amet, consectetur adipiscing",
					UserID:     1,
					RepostOfID: 0,
					CreatedAt:  time.Now(),
				},
				{
					ID:         2,
					Text:       "Lorem ipsum dolor sit amet, consectetur adipiscing",
					UserID:     2,
					RepostOfID: 0,
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &PostRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "text", "user_id", "repost_of_id", "created_at",
				})
				for _, post := range tc.readDTOs {
					rows.AddRow(post.ID, post.Text, post.UserID, post.RepostOfID, post.CreatedAt)
				}

				mock.ExpectQuery(regexp.QuoteMeta(`select id, text, user_id, repost_of_id, created_at from posts where deleted_at is null offset $1 limit $2;`)).
					WithArgs(0, 100).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, text, user_id, repost_of_id, created_at from posts where deleted_at is null offset $1 limit $2;`)).
					WithArgs(0, 100).
					WillReturnError(sql.ErrNoRows)
			}

			posts, err := r.GetAllPosts(uint64(100), uint64(0))
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTOs, posts, "Posts mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostRepositoryImpl_GetPostByID(t *testing.T) {
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
				ID:         1,
				Text:       "Lorem ipsum dolor sit amet, consectetur adipiscing",
				UserID:     1,
				RepostOfID: 0,
				CreatedAt:  time.Now(),
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &PostRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "text", "user_id", "repost_of_id", "created_at",
				}).AddRow(
					tc.readDTO.ID, tc.readDTO.Text, tc.readDTO.UserID, tc.readDTO.RepostOfID, tc.readDTO.CreatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`select id, text, user_id, repost_of_id, created_at from posts where id = $1 and deleted_at is null;`)).
					WithArgs(tc.id).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, text, user_id, repost_of_id, created_at from posts where id = $1 and deleted_at is null;`)).
					WithArgs(tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			post, err := r.GetPostByID(tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, post, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, post, "Post mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostRepositoryImpl_DeletePost(t *testing.T) {
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
			name:     "Error on delete SQL",
			id:       2,
			hasError: true,
		},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &PostRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectExec(regexp.QuoteMeta(`update posts set deleted_at = now() where id = $1 and deleted_at is null;`)).
					WithArgs(tc.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`update posts set deleted_at = now() where id = $1 and deleted_at is null;`)).
					WithArgs(tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			err := r.DeletePost(tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}
