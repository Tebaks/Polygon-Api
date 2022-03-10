package util

import (
	"strings"

	v "github.com/go-playground/validator"
)

type Validator interface {
	Validate(i interface{}) error
}

type validator struct {
	instance *v.Validate
}

func NewValidator() Validator {
	return &validator{
		instance: v.New(),
	}
}

func (v *validator) Validate(i interface{}) error {
	return v.instance.Struct(i)
}

func GetValidationErrorString(err error) string {
	ve, ok := err.(v.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var errs []string
	for _, e := range ve {
		item := strings.Join([]string{e.Field(), e.Tag(), e.Param()}, " ")
		item = strings.TrimSpace(item)
		errs = append(errs, item)
	}
	return strings.Join(errs, "; ")
}
