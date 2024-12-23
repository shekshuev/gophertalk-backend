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
	ID        uint64    `json:"id"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReadAuthUserDataDTO struct {
	ID           uint64 `json:"id"`
	UserName     string `json:"user_name"`
	PasswordHash string `json:"first_name"`
	Status       uint8  `json:"status"`
}

type UpdateUserDTO struct {
	UserName        string `json:"user_name" validate:"omitempty,min=5,max=30,alphanumunderscore,startswithalpha"`
	Password        string `json:"password" validate:"omitempty,password"`
	PasswordConfirm string `json:"password_confirm" validate:"omitempty,password,eqfield=Password"`
	PasswordHash    string
	FirstName       string `json:"first_name" validate:"omitempty,min=1,max=30,alphaunicode"`
	LastName        string `json:"last_name" validate:"omitempty,min=1,max=30,alphaunicode"`
}
