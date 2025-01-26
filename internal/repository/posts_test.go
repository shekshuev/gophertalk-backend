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
				OwnerID:   0,
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
					UserViewed: true,
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
					UserViewed: true,
				},
			},
			hasError: false,
		},
		{
			name: "Error on SQL query",
			filterDTO: models.FilterPostDTO{
				UserID:    1,
				OwnerID:   0,
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
	select 
		p.id AS post_id,
		p.text,
		p.reply_to_id,
		p.created_at,
		u.id AS user_id,
		u.user_name,
		u.first_name,
		u.last_name,
		p.likes_count,
		p.views_count,
		p.replies_count,
		case 
			when l.user_id is not null then true
			else false
		end as user_liked,
		case 
			when v.user_id is not null then true
			else false
		end as user_viewed
	from posts p
	join users u ON p.user_id = u.id
	left join likes l on l.post_id = p.id and l.user_id = $1
	left join views v on v.post_id = p.id and v.user_id = $1
	where p.deleted_at is null and p.text ilike $2 and p.reply_to_id = $3
	order by p.created_at asc
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
					"replies_count",
					"user_liked",
					"user_viewed",
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
						post.RepliesCount,
						post.UserLiked,
						post.UserViewed)
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
				UserViewed: true,
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
	select 
		p.id AS post_id,
		p.text,
		p.reply_to_id,
		p.created_at,
		u.id AS user_id,
		u.user_name,
		u.first_name,
		u.last_name,
		p.likes_count,
		p.views_count,
		p.replies_count,
		case 
			when l.user_id is not null then true
			else false
		end as user_liked,
		case 
			when v.user_id is not null then true
			else false
		end as user_viewed
	from posts p
	join users u ON p.user_id = u.id
	left join likes l on l.post_id = p.id and l.user_id = $1
	left join views v on v.post_id = p.id and v.user_id = $1
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
					"replies_count",
					"user_liked",
					"user_viewed",
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
					tc.readDTO.RepliesCount,
					tc.readDTO.UserLiked,
					tc.readDTO.UserViewed,
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
		name       string
		id         uint64
		maxRecords int
	}{
		{
			name:       "Success view post",
			id:         1,
			maxRecords: 1,
		},
		{
			name:       "Success cache post",
			id:         1,
			maxRecords: 2,
		},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vb := &ViewBuffer{
				buffer:     make([]View, 0, tc.maxRecords),
				maxRecords: tc.maxRecords,
				timer:      time.Second,
			}
			r := &PostRepositoryImpl{cfg: &cfg, db: db, vb: vb}
			if tc.maxRecords == 1 {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("create temp table tmp_views (post_id bigint, user_id bigint, created_at timestamp) on commit drop;")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`insert into tmp_views (post_id, user_id, created_at) values ($1, $2, $3)`)).
					WithArgs(tc.id, uint64(1), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`
					insert into views (post_id, user_id, created_at)
					select post_id, user_id, min(created_at) as created_at
					from tmp_views group by post_id, user_id on conflict (user_id, post_id) do nothing;
					`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`
					update posts set views_count = v.count from (
						select post_id, count(post_id) as count from views where post_id in (
							select distinct post_id from tmp_views
						)
						group by post_id
					) v where posts.id = v.post_id;
					`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			err := r.ViewPost(tc.id, uint64(1))

			assert.Nil(t, err, "Error is not nil")

			if tc.maxRecords > 1 {
				assert.Len(t, r.vb.buffer, 1, "Max records exceeded")
			}

			if err := mock.ExpectationsWereMet(); err != nil && tc.maxRecords == 1 {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostRepositoryImpl_LikePost(t *testing.T) {
	testCases := []struct {
		name       string
		id         uint64
		maxRecords int
	}{
		{
			name:       "Success like post",
			id:         1,
			maxRecords: 1,
		},
		{
			name:       "Success cache post",
			id:         1,
			maxRecords: 2,
		},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lb := &LikeBuffer{
				likeBuffer: make([]Like, 0, tc.maxRecords),
				maxRecords: tc.maxRecords,
				timer:      time.Second,
			}
			r := &PostRepositoryImpl{cfg: &cfg, db: db, lb: lb}
			if tc.maxRecords == 1 {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("create temp table tmp_likes (post_id bigint, user_id bigint, created_at timestamp) on commit drop;")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`insert into tmp_likes (post_id, user_id, created_at) values ($1, $2, $3)`)).
					WithArgs(tc.id, uint64(1), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`
					insert into likes (post_id, user_id, created_at)
					select post_id, user_id, min(created_at) as created_at
					from tmp_likes group by post_id, user_id on conflict (user_id, post_id) do nothing;
					`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`
					update posts set likes_count = l.count from (
						select post_id, count(post_id) as count from likes where post_id in (
							select distinct post_id from tmp_likes
						)
						group by post_id
					) l where posts.id = l.post_id;
					`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			err := r.LikePost(tc.id, uint64(1))

			assert.Nil(t, err, "Error is not nil")

			if tc.maxRecords > 1 {
				assert.Len(t, r.lb.likeBuffer, 1, "Max records exceeded")
			}

			if err := mock.ExpectationsWereMet(); err != nil && tc.maxRecords == 1 {
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
