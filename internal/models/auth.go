package models

type LoginUserDTO struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type RegisterUserDTO struct {
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

type ReadTokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
