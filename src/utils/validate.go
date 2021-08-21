package utils

import (
	"github.com/go-playground/validator/v10"
)

func Validate(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}
