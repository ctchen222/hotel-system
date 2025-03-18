package response

import "net/http"

type Error struct {
	Code   int    `json:"code"`
	Extras string `json:"extras"`
}

func (e Error) Error() string {
	return e.Extras
}

func NewError(code int, message string) Error {
	return Error{
		Code:   code,
		Extras: message,
	}
}

func ErrInvalidId() Error {
	return NewError(http.StatusBadRequest, "Invalid ID")
}

func ErrUnAuthenticated() Error {
	return NewError(http.StatusUnauthorized, "Forbidden")
}

func ErrUnAuthorized() Error {
	return NewError(http.StatusForbidden, "Forbidden")
}

func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "Invalid JSON Request")
}

func ErrResourceNotFound() Error {
	return NewError(http.StatusBadRequest, "Resource not found")
}
