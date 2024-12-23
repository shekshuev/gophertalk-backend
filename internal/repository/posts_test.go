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
				Text:      "Lorem ipsum dolor sit amet, consectetur adipiscing",
				UserID:    1,
				ReplyToID: nil,
			},
			readDTO: models.ReadPostDTO{
				ID:        1,
				Text:      "Lorem ipsum dolor sit amet, consectetur adipiscing",
				ReplyToID: nil,
				CreatedAt: time.Now(),
			},
			hasError: false,
		},
		{
			name: "Error on insert SQL",
			createDTO: models.CreatePostDTO{
				Text:      "Lorem ipsum dolor sit amet, consectetur adipiscing",
				UserID:    1,
				ReplyToID: nil,
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
					insert into posts (text, user_id, reply_to_id) values ($1, $2, $3)
					returning id, text, created_at, reply_to_id;
					`)).
					WithArgs(
						tc.createDTO.Text,
						tc.createDTO.UserID,
						tc.createDTO.ReplyToID).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "text", "created_at", "reply_to_id"},
						).AddRow(
							tc.readDTO.ID,
							tc.readDTO.Text,
							tc.readDTO.CreatedAt,
							tc.readDTO.ReplyToID,
						),
					)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into posts (text, user_id, reply_to_id) values ($1, $2, $3)
					returning id, text, created_at, reply_to_id;
					`)).
					WithArgs(
						tc.createDTO.Text,
						tc.createDTO.UserID,
						tc.createDTO.ReplyToID).
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
		name      string
		filterDTO models.FilterPostDTO
		readDTOs  []models.ReadPostDTO
		hasError  bool
	}{
		{
			name: "Success get all posts",
			filterDTO: models.FilterPostDTO{
				UserID:    1,
				Limit:     100,
				Offset:    0,
				ReplyToID: 1,
				Search:    "test",
			},
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
					ReplyToID:  nil,
					CreatedAt:  time.Now(),
					LikesCount: 10,
					ViewsCount: 100,
					UserLiked:  true,
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
					ReplyToID:  nil,
					CreatedAt:  time.Now(),
					LikesCount: 10,
					ViewsCount: 100,
					UserLiked:  true,
				},
			},
			hasError: false,
		},
		{
			name: "Error on SQL query",
			filterDTO: models.FilterPostDTO{
				UserID:    1,
				Limit:     100,
				Offset:    0,
				ReplyToID: 1,
				Search:    "test",
			},
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

	query := regexp.QuoteMeta(`
	with likes_count AS (
		select post_id, count(*) as likes_count
		from likes group by post_id
	),
	views_count as (
		select post_id, count(*) AS views_count
		from views group by post_id
	)
	select 
		p.id AS post_id,
		p.text,
		p.reply_to_id,
		p.created_at,
		u.id AS user_id,
		u.user_name,
		u.first_name,
		u.last_name,
		coalesce(lc.likes_count, 0) AS likes_count,
		coalesce(vc.views_count, 0) AS views_count,
		case 
			when l.user_id is not null then true
			else false
		end as user_liked
	from posts p
	join users u ON p.user_id = u.id
	left join likes_count lc ON p.id = lc.post_id
	left join views_count vc ON p.id = vc.post_id
	left join likes l on l.post_id = p.id and l.user_id = $1
	where p.deleted_at is null and p.text ilike $2 and p.reply_to_id = $3
	offset $4 limit $5
	`)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"p.id",
					"p.text",
					"p.reply_to_id",
					"p.created_at",
					"u.id",
					"u.user_name",
					"u.first_name",
					"u.last_name",
					"likes_count",
					"views_count",
					"user_liked",
				})
				for _, post := range tc.readDTOs {
					rows.AddRow(
						post.ID,
						post.Text,
						post.ReplyToID,
						post.CreatedAt,
						post.User.ID,
						post.User.UserName,
						post.User.FirstName,
						post.User.LastName,
						post.LikesCount,
						post.ViewsCount,
						post.UserLiked)
				}

				mock.ExpectQuery(query).
					WithArgs(tc.filterDTO.UserID, "%"+tc.filterDTO.Search+"%", tc.filterDTO.ReplyToID, tc.filterDTO.Offset, tc.filterDTO.Limit).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(query).
					WithArgs(tc.filterDTO.UserID, "%"+tc.filterDTO.Search+"%", tc.filterDTO.ReplyToID, tc.filterDTO.Offset, tc.filterDTO.Limit).
					WillReturnError(sql.ErrNoRows)
			}

			posts, err := r.GetAllPosts(tc.filterDTO)
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
		userID   uint64
		readDTO  *models.ReadPostDTO
		hasError bool
	}{
		{
			name:   "Success get post by ID",
			id:     1,
			userID: 1,
			readDTO: &models.ReadPostDTO{
				ID:   1,
				Text: "Lorem ipsum dolor sit amet, consectetur adipiscing",
				User: &models.ReadPostUserDTO{
					ID:        1,
					UserName:  "username",
					FirstName: "first_name",
					LastName:  "last_name",
				},
				ReplyToID:  nil,
				CreatedAt:  time.Now(),
				LikesCount: 10,
				ViewsCount: 100,
				UserLiked:  true,
			},
			hasError: false,
		},
		{
			name:     "Post not found",
			id:       1,
			userID:   1,
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

	query := regexp.QuoteMeta(`
	with likes_count AS (
		select post_id, count(*) as likes_count
		from likes group by post_id
	),
	views_count as (
		select post_id, count(*) AS views_count
		from views group by post_id
	)
	select 
		p.id AS post_id,
		p.text,
		p.reply_to_id,
		p.created_at,
		u.id AS user_id,
		u.user_name,
		u.first_name,
		u.last_name,
		coalesce(lc.likes_count, 0) AS likes_count,
		coalesce(vc.views_count, 0) AS views_count,
		case 
			when l.user_id is not null then true
			else false
		end as user_liked
	from posts p
	join users u ON p.user_id = u.id
	left join likes_count lc on p.id = lc.post_id
	left join views_count vc on p.id = vc.post_id
	left join likes l on l.post_id = p.id and l.user_id = $1
	where p.id = $2 and p.deleted_at is null;
	`)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"p.id",
					"p.text",
					"p.reply_to_id",
					"p.created_at",
					"u.id",
					"u.user_name",
					"u.first_name",
					"u.last_name",
					"likes_count",
					"views_count",
					"user_liked",
				}).AddRow(
					tc.readDTO.ID,
					tc.readDTO.Text,
					tc.readDTO.ReplyToID,
					tc.readDTO.CreatedAt,
					tc.readDTO.User.ID,
					tc.readDTO.User.UserName,
					tc.readDTO.User.FirstName,
					tc.readDTO.User.LastName,
					tc.readDTO.LikesCount,
					tc.readDTO.ViewsCount,
					tc.readDTO.UserLiked,
				)

				mock.ExpectQuery(query).
					WithArgs(tc.userID, tc.id).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(query).
					WithArgs(tc.userID, tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			post, err := r.GetPostByID(tc.userID, tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, post, "Post should be nil")
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
				mock.ExpectExec(regexp.QuoteMeta(`update posts set deleted_at = now() where id = $1 and user_id = $2 and deleted_at is null;`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`update posts set deleted_at = now() where id = $1 and user_id = $2 and deleted_at is null;`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnError(sql.ErrNoRows)
			}

			err := r.DeletePost(tc.id, uint64(1))
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

func TestPostRepositoryImpl_ViewPost(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint64
		hasError bool
	}{
		{
			name:     "Success view post",
			id:       1,
			hasError: false,
		},
		{
			name:     "Error on view SQL",
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
				mock.ExpectExec(regexp.QuoteMeta(`insert into views (post_id, user_id) values ($1, $2);`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`insert into views (post_id, user_id) values ($1, $2);`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnError(sql.ErrNoRows)
			}

			err := r.ViewPost(tc.id, uint64(1))
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

func TestPostRepositoryImpl_LikePost(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint64
		hasError bool
	}{
		{
			name:     "Success like post",
			id:       1,
			hasError: false,
		},
		{
			name:     "Error on like SQL",
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
				mock.ExpectExec(regexp.QuoteMeta(`insert into likes (post_id, user_id) values ($1, $2);`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`insert into likes (post_id, user_id) values ($1, $2);`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnError(sql.ErrNoRows)
			}

			err := r.LikePost(tc.id, uint64(1))
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

func TestPostRepositoryImpl_DislikePost(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint64
		hasError bool
	}{
		{
			name:     "Success dislike post",
			id:       1,
			hasError: false,
		},
		{
			name:     "Error on dislike SQL",
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
				mock.ExpectExec(regexp.QuoteMeta(`delete from likes where post_id = $1 and user_id = $2;`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`delete from likes where post_id = $1 and user_id = $2;`)).
					WithArgs(tc.id, uint64(1)).
					WillReturnError(sql.ErrNoRows)
			}

			err := r.DislikePost(tc.id, uint64(1))
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
