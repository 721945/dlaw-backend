package libs

import (
	"fmt"
	"net/http"
)

var (
	ErrInternalServerError = fmt.Errorf("internal server error")
	ErrBadRequest          = fmt.Errorf("bad request")
	ErrUnauthorized        = fmt.Errorf("unauthorized")
	ErrForbidden           = fmt.Errorf("forbidden")
	ErrNotFound            = fmt.Errorf("not found")
	ErrConflict            = fmt.Errorf("conflict")
	ErrBadParamInput       = fmt.Errorf("bad param input")
	ErrBadParamInputFormat = fmt.Errorf("bad param input format")
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

//var (
//	NotFoundError = fmt.Errorf("resource could not be found")
//)
