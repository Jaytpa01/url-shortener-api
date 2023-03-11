package api

import (
	"errors"
	"net/http"
	"strings"
)

// ErrorType holds a string and represents the integer HTTP status code for the error
type ErrorType string

const (
	Unknown              ErrorType = "UNKNOWN"
	NotFound             ErrorType = "NOT_FOUND"
	Internal             ErrorType = "INTERNAL"
	UnsupportedMediaType ErrorType = "UNSUPPORTED"
)

// ApiError is a custom error for the application.
// It is helpful in returning a consistent
// error response from API endpoints.
type ApiError struct {
	Type    ErrorType `json:"type,omitempty"`    // the type of error
	Code    string    `json:"code,omitempty"`    // a unique code for every insance an error is created
	Message string    `json:"message,omitempty"` // human readable msg about error
	Action  string    `json:"action,omitempty"`  // actions you can take to resolve the error
	Debug   string    `json:"-"`                 // the actual internal error thrown
}

// ApiError.Error() satisfies the standard error interface
func (a *ApiError) Error() string {
	return a.Message
}

func EnsureApiError(err error) *ApiError {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return NewUnknown(err.Error())
}

// Status maps an ErrorType to a HTTP Status Code
func (a *ApiError) Status() int {
	switch a.Type {
	case NotFound:
		return http.StatusNotFound
	case Internal:
		return http.StatusInternalServerError
	case Unknown:
		return http.StatusInternalServerError
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}

/*
	Error Factories
*/

// NewInternal for HTTP Status 500 Server Errors
func NewInternal(code, debug string) *ApiError {
	return &ApiError{
		Type:  Internal,
		Code:  code,
		Debug: debug,
	}
}

// NewUnknown for unknown errors. Uses HTTP Status 500
func NewUnknown(debug string) *ApiError {
	return &ApiError{
		Type:  Unknown,
		Code:  "unknown",
		Debug: debug,
	}
}

// func NewBadRequest()

// NewUnsupportedMediaType to creat HTTP Status 415 errors.
// Often used when the Content-Type header of a request is not what its expected to be.
func NewUnsupportedMediaType(code, msg, action string, debug ...string) *ApiError {
	return &ApiError{
		Type:    UnsupportedMediaType,
		Code:    code,
		Message: msg,
		Action:  action,
		Debug:   strings.Join(debug, " "),
	}
}
