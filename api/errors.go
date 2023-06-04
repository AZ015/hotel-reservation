package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}

	apiErr := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiErr.Code).JSON(apiErr.Err)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidID() Error {
	return NewError(http.StatusBadRequest, "invalid id given")
}

func ErrNotAuthorized() Error {
	return NewError(http.StatusUnauthorized, "unauthorized request")
}

func ErrTokenExpired() Error {
	return NewError(http.StatusUnauthorized, "token expired")
}

func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "invalid JSON request")
}

func ErrResourceNotFound(resource string) Error {
	return NewError(http.StatusNotFound, fmt.Sprintf("%s resource not found", resource))
}
