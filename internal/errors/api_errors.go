package errors

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrInternalServer    = errors.New("internal server error")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrNotFound          = errors.New("not found")
)

func ErrStatusCode(err error) int {
	switch {
	case strings.Contains(err.Error(), ErrInternalServer.Error()):
		return http.StatusInternalServerError
	case strings.Contains(err.Error(), ErrInvalidCredential.Error()):
		return http.StatusUnauthorized
	case strings.Contains(err.Error(), ErrNotFound.Error()):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
