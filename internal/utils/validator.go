package utils

import (
	"log"

	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	err := validate.RegisterValidation("alphanumunderscore", alphanumUnderscore)
	if err != nil {
		log.Fatalf("Error registering validator: %v", err)
	}
	err = validate.RegisterValidation("startswithalpha", startsWithAlpha)
	if err != nil {
		log.Fatalf("Error registering validator: %v", err)
	}
	err = validate.RegisterValidation("password", passwordValidation)
	if err != nil {
		log.Fatalf("Error registering validator: %v", err)
	}
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
