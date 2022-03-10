package handler

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

func SuccessWithCodeResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, Response{
		Success: true,
		Error:   "",
		Data:    data,
	})
}

func ErrorWithCodeResponse(c echo.Context, code int, err error) error {
	return c.JSON(code, Response{
		Success: false,
		Error:   err.Error(),
		Data:    nil,
	})
}
