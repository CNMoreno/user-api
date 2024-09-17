package utils

import "github.com/go-playground/validator/v10"

// NewValidator validate files.
func NewValidator() *validator.Validate {
	validate := validator.New()
	RegisterCustomValidators(validate)

	return validate
}
