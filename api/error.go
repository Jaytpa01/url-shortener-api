package api

import (
	"errors"
	"net/http"
)

// ErrorType holds a string and represents the integer HTTP status code for the error
type ErrorType string

const (
	Unknown              ErrorType = "UNKNOWN"
	NotFound             ErrorType = "NOT_FOUND"
	Internal             ErrorType = "INTERNAL"
	UnsupportedMediaType ErrorType = "UNSUPPORTED"
	BadRequest           ErrorType = "BAD_REQUEST"
	PayloadTooLarge      ErrorType = "PAYLOAD_TOO_LARGE"
)

// ApiError is a custom error for the application.
// It is helpful in returning a consistent
// error response from API endpoints.
type ApiError struct {
	Type    ErrorType `json:"type,omitempty"`    // the type of error
	Code    string    `json:"code,omitempty"`    // a unique code for every insance an error is created
	Message string    `json:"message,omitempty"` // human readable message about the error

	// Optional fields
	Action string `json:"action,omitempty"` // actions you can take to resolve the error
	Debug  string `json:"-"`                // the actual internal error thrown
}

// ApiError.Error() satisfies the standard error interface
func (ae *ApiError) Error() string {
	return ae.Message
}

// Option type allows us to implement optional functional parameters
// to fill optional fields of an ApiError.
// https://levelup.gitconnected.com/optional-function-parameter-pattern-in-golang-c1acc829307b
//
// Currently implemented options:
// WithAction(action string),
// WithDebug(debug string)
type ErrorOption func(*ApiError)

// WithAction lets us optionally append an action to
// resolve the it to an ApiError
func WithAction(action string) ErrorOption {
	return func(ae *ApiError) {
		ae.Action = action
	}
}

// WithDebug allows us to optionally attach a debug message
// to an error. This is usually the error itself.
// ie. error.Error()
func WithDebug(debug string) ErrorOption {
	return func(ae *ApiError) {
		ae.Debug = debug
	}
}

func EnsureApiError(err error) *ApiError {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return NewUnknown(WithDebug(err.Error()))
}

// Status maps an ErrorType to a HTTP Status Code
func (ae *ApiError) Status() int {
	switch ae.Type {
	case NotFound:
		return http.StatusNotFound
	case Internal:
		return http.StatusInternalServerError
	case Unknown:
		return http.StatusInternalServerError
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case BadRequest:
		return http.StatusBadRequest
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

/*
	Error Factories
*/

// NewInternal for HTTP Status 500 Server Errors
func NewInternal(code string, opts ...ErrorOption) *ApiError {
	ae := &ApiError{
		Type:    Internal,
		Code:    code,
		Message: "An internal server error occured.",
	}

	applyErrorOptions(ae, opts...)
	return ae
}

// NewUnknown for unknown errors. Uses HTTP Status 500
func NewUnknown(opts ...ErrorOption) *ApiError {
	ae := &ApiError{
		Type: Unknown,
		Code: "unknown",
	}

	applyErrorOptions(ae, opts...)

	return ae
}

func NewBadRequest(code, msg string, opts ...ErrorOption) *ApiError {
	ae := &ApiError{
		Type:    BadRequest,
		Code:    code,
		Message: msg,
	}
	applyErrorOptions(ae, opts...)
	return ae
}

// NewUnsupportedMediaType to creat HTTP Status 415 errors.
// Often used when the Content-Type header of a request is not what its expected to be.
func NewUnsupportedMediaType(code, msg string, opts ...ErrorOption) *ApiError {
	ae := &ApiError{
		Type:    UnsupportedMediaType,
		Code:    code,
		Message: msg,
	}

	applyErrorOptions(ae, opts...)
	return ae
}

// NewRequestPayloadTooLarge is used when returning a HTTP Status 413 error to the client.
func NewRequestPayloadTooLarge(code, msg string, opts ...ErrorOption) *ApiError {
	ae := &ApiError{
		Type:    PayloadTooLarge,
		Code:    code,
		Message: msg,
	}
	applyErrorOptions(ae, opts...)
	return ae
}

// applyErrorOptions is a helper function to apply any options
// in our wrror factories
func applyErrorOptions(ae *ApiError, opts ...ErrorOption) {
	for _, opt := range opts {
		opt(ae)
	}
}
