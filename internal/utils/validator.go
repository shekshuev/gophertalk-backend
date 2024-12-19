package utils

import (
	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("alphanumunderscore", alphanumUnderscore)
	validate.RegisterValidation("startswithalpha", startsWithAlpha)
	validate.RegisterValidation("password", passwordValidation)
	return validate
}

func alphanumUnderscore(fl validator.FieldLevel) bool {
	return regexValidation(fl.Field().String(), `^[a-zA-Z0-9_]+$`)
}

func startsWithAlpha(fl validator.FieldLevel) bool {
	return regexValidation(fl.Field().String(), `^[^0-9]`)
}

func passwordValidation(fl validator.FieldLevel) bool {
	return regexValidation(fl.Field().String(), `^(?=.*[a-zA-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{5,30}$`)
}

func regexValidation(field, regex string) bool {
	re := regexp2.MustCompile(regex, regexp2.None)
	matched, _ := re.MatchString(field)
	return matched
}
