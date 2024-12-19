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

func TestPostgresUserRepository_CreateUser(t *testing.T) {
	testCases := []struct {
		name      string
		createDTO models.CreateUserDTO
		readDTO   models.ReadAuthUserDataDTO
		hasError  bool
	}{
		{
			name: "Success create",
			createDTO: models.CreateUserDTO{
				UserName:     "john",
				FirstName:    "John",
				LastName:     "Doe",
				PasswordHash: "password",
			},
			readDTO: models.ReadAuthUserDataDTO{
				ID:           1,
				UserName:     "john",
				Status:       1,
				PasswordHash: "password",
			},
			hasError: false,
		},
		{
			name: "Error on insert SQL",
			createDTO: models.CreateUserDTO{
				UserName:     "john",
				FirstName:    "John",
				LastName:     "Doe",
				PasswordHash: "password",
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
	r := &UserRepositoryImpl{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into users (user_name, first_name, last_name, password_hash) values ($1, $2, $3, $4) 
					returning id, user_name, password_hash, status;
					`)).
					WithArgs(
						tc.createDTO.UserName,
						tc.createDTO.FirstName,
						tc.createDTO.LastName,
						tc.createDTO.PasswordHash).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "user_name", "password_hash", "status"},
						).AddRow(
							tc.readDTO.ID,
							tc.readDTO.UserName,
							tc.readDTO.PasswordHash,
							tc.readDTO.Status,
						),
					)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into users (user_name, first_name, last_name, password_hash) values ($1, $2, $3, $4) 
					returning id, user_name, password_hash, status;
					`)).
					WithArgs(
						tc.createDTO.UserName,
						tc.createDTO.FirstName,
						tc.createDTO.LastName,
						tc.createDTO.PasswordHash).
					WillReturnError(sql.ErrNoRows)
			}
			user, err := r.CreateUser(tc.createDTO)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, *user, "User mismatch")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresUserRepository_GetAllUsers(t *testing.T) {
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "user_name", "first_name", "last_name", "status", "created_at", "updated_at",
				})
				for _, user := range tc.readDTOs {
					rows.AddRow(user.ID, user.UserName, user.FirstName, user.LastName, user.Status, user.CreatedAt, user.UpdatedAt)
				}

				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, first_name, last_name, status, created_at, updated_at from users where deleted_at is null offset $1 limit $2;`)).
					WithArgs(0, 100).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, first_name, last_name, status, created_at, updated_at from users where deleted_at is null offset $1 limit $2;`)).
					WithArgs(0, 100).
					WillReturnError(sql.ErrNoRows)
			}

			users, err := r.GetAllUsers(100, 0)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTOs, users, "Users mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresUserRepository_GetUserByID(t *testing.T) {
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "user_name", "first_name", "last_name", "status", "created_at", "updated_at",
				}).AddRow(
					tc.readDTO.ID, tc.readDTO.UserName, tc.readDTO.FirstName, tc.readDTO.LastName, tc.readDTO.Status, tc.readDTO.CreatedAt, tc.readDTO.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, first_name, last_name, status, created_at, updated_at from users where id = $1 and deleted_at is null;`)).
					WithArgs(tc.id).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, first_name, last_name, status, created_at, updated_at from users where id = $1 and deleted_at is null;`)).
					WithArgs(tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			user, err := r.GetUserByID(tc.id)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "User mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresUserRepository_GetUserByUserName(t *testing.T) {
	testCases := []struct {
		name     string
		userName string
		readDTO  *models.ReadAuthUserDataDTO
		hasError bool
	}{
		{
			name:     "Success get user by ID",
			userName: "john",
			readDTO: &models.ReadAuthUserDataDTO{
				ID:           1,
				UserName:     "john",
				PasswordHash: "password",
				Status:       models.StatusActive,
			},
			hasError: false,
		},
		{
			name:     "User not found",
			userName: "notfound",
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

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "user_name", "password_hash", "status",
				}).AddRow(
					tc.readDTO.ID, tc.readDTO.UserName, tc.readDTO.PasswordHash, tc.readDTO.Status,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, password_hash, status from users where user_name = $1 and deleted_at is null;`)).
					WithArgs(tc.userName).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, password_hash, status from users where user_name = $1 and deleted_at is null;`)).
					WithArgs(tc.userName).
					WillReturnError(sql.ErrNoRows)
			}

			user, err := r.GetUserByUserName(tc.userName)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "User mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresUserRepository_UpdateUser(t *testing.T) {
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "user_name", "first_name", "last_name", "status", "created_at", "updated_at",
				}).AddRow(
					tc.readDTO.ID, tc.readDTO.UserName, tc.readDTO.FirstName, tc.readDTO.LastName, tc.readDTO.Status, tc.readDTO.CreatedAt, tc.readDTO.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`update users set password_hash = $1, user_name = $2, first_name = $3, last_name = $4, updated_at = now() where id = $5 and deleted_at is null returning *`)).
					WithArgs(
						tc.updateDTO.PasswordHash,
						tc.updateDTO.UserName,
						tc.updateDTO.FirstName,
						tc.updateDTO.LastName,
						tc.id).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`update users set user_name = $1, updated_at = now() where id = $2 and deleted_at is null returning *`)).
					WithArgs(tc.updateDTO.UserName, tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			user, err := r.UpdateUser(tc.id, tc.updateDTO)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "User mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresUserRepository_DeleteUser(t *testing.T) {
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectExec(regexp.QuoteMeta(`update users set deleted_at = now() where id = $1;`)).
					WithArgs(tc.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`update users set deleted_at = now() where id = $1;`)).
					WithArgs(tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			err := r.DeleteUser(tc.id)
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
