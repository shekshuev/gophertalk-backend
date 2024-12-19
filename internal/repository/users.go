package repository

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type UserRepositoryImpl struct {
	db  *sql.DB
	cfg *config.Config
}

func NewUserRepositoryImpl(cfg *config.Config) *UserRepositoryImpl {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Error connecting to database", err)
		return nil
	}
	repository := &UserRepositoryImpl{cfg: cfg, db: db}
	return repository
}

func (r *UserRepositoryImpl) CreateUser(dto models.CreateUserDTO) (*models.ReadAuthUserDataDTO, error) {
	query := `
		insert into users (user_name, first_name, last_name, password_hash) values ($1, $2, $3, $4)
		returning id, user_name, password_hash, status;
	`
	var user models.ReadAuthUserDataDTO
	err := r.db.QueryRow(
		query, dto.UserName, dto.FirstName, dto.LastName, dto.PasswordHash).Scan(&user.ID, &user.UserName, &user.PasswordHash, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetAllUsers() ([]models.ReadUserDTO, error) {
	query := `
		select id, user_name, first_name, last_name, status, created_at, updated_at from users;
	`
	var readDTO []models.ReadUserDTO
	rows, err := r.db.Query(query)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.ReadUserDTO
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		readDTO = append(readDTO, user)
	}
	return readDTO, nil
}

func (r *UserRepositoryImpl) GetUserByID(id int) (*models.ReadUserDTO, error) {
	query := `
		select 
			id, user_name, first_name, last_name, status, created_at, updated_at 
		from users where id = $1 and deleted_at is null;
	`
	var user models.ReadUserDTO
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByUserName(userName string) (*models.ReadAuthUserDataDTO, error) {
	query := `
		select 
			id, user_name, password_hash, status 
		from users where user_name = $1 and deleted_at is null;
	`
	var user models.ReadAuthUserDataDTO
	err := r.db.QueryRow(query, userName).Scan(&user.ID, &user.UserName, &user.PasswordHash, &user.Status)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(id int, dto models.UpdateUserDTO) (*models.ReadUserDTO, error) {
	fields := make([]string, 0)
	args := make([]interface{}, 0)
	if dto.PasswordHash != "" {
		fields = append(fields, fmt.Sprintf("password_hash = $%d", len(args)+1))
		args = append(args, dto.PasswordHash)
	}
	if dto.UserName != "" {
		fields = append(fields, fmt.Sprintf("user_name = $%d", len(args)+1))
		args = append(args, dto.UserName)
	}
	if dto.FirstName != "" {
		fields = append(fields, fmt.Sprintf("first_name = $%d", len(args)+1))
		args = append(args, dto.FirstName)
	}
	if dto.LastName != "" {
		fields = append(fields, fmt.Sprintf("last_name = $%d", len(args)+1))
		args = append(args, dto.LastName)
	}
	if len(fields) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	fields = append(fields, "updated_at = now()")

	query := fmt.Sprintf("update users set %s where id = $%d and deleted_at is null returning *",
		strings.Join(fields, ", "), len(args)+1)

	args = append(args, id)

	var user models.ReadUserDTO
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) DeleteUser(id int) error {
	query := `
        update users set deleted_at = now() where id = $1;
    `
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
