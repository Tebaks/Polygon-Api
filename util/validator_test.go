//go:build unit
// +build unit
package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	UserName          string `validate:"required"`
	Password          string `validate:"required,gt=15"`
	VeryRequiredField string `validate:"required"`
	Desc              string `validate:""`
}

func TestGetValidationErrorString(t *testing.T) {
	qt := assert.New(t)

	instance := testStruct{
		UserName: "username",
		Password: "secret",
	}

	tests := []struct {
		err error
		exp string
	}{
		{
			err: errors.New("Not validator error"),
			exp: "Not validator error",
		},
		{
			err: NewValidator().Validate(instance),
			exp: "Password gt 15; VeryRequiredField required",
		},
	}

	for _, test := range tests {
		res := GetValidationErrorString(test.err)
		qt.Equal(test.exp, res)
	}

}
