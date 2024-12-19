package models

import "time"

const (
	StatusBlocked = iota
	StatusActive
)

type CreateUserDTO struct {
	UserName     string
	PasswordHash string
	FirstName    string
	LastName     string
}

type ReadUserDTO struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReadAuthUserDataDTO struct {
	ID           int    `json:"id"`
	UserName     string `json:"user_name"`
	PasswordHash string `json:"first_name"`
	Status       int    `json:"status"`
}

type UpdateUserDTO struct {
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	PasswordHash    string
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}
