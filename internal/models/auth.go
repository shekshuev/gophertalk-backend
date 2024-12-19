package models

type LoginUserDTO struct {
	UserName string `json:"user_name" validate:"required,min=5,max=30,alphanumunderscore,startswithalpha"`
	Password string `json:"password" validate:"required,password"`
}

type RegisterUserDTO struct {
	UserName        string `json:"user_name" validate:"required,min=5,max=30,alphanumunderscore,startswithalpha"`
	Password        string `json:"password" validate:"required,password"`
	PasswordConfirm string `json:"password_confirm" validate:"required,password"`
	FirstName       string `json:"first_name" validate:"required,min=1,max=30,alphaunicode"`
	LastName        string `json:"last_name" validate:"required,min=1,max=30,alphaunicode"`
}

type ReadTokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
