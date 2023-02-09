package controllers

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrBadParamInput       = errors.New("bad param input")
)

func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
