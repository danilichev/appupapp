package services

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	validator := validator.New()
	validator.RegisterValidation("notblank", notBlank)

	return &Validator{validator}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}
