package errors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

// AppError represents a structured error in our application
// It carries semantic meaning and can be mapped to different interfaces
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"` // Original error, not exposed to clients
}

// Error implements the error interface
func (appError *AppError) Error() string {
	if appError.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", appError.Code, appError.Message, appError.Details)
	}

	return fmt.Sprintf("%s: %s", appError.Code, appError.Message)
}

// ToHTTPStatus converts an AppError to HTTP status code
func (appError *AppError) ToHTTPStatus() int {
	return ErrorCodeToHTTPStatus(appError.Code)
}

// ToGRPCStatus converts an AppError to gRPC status
func (appError *AppError) ToGRPCStatus() *status.Status {
	return ErrorCodeToGRPCStatus(appError.Code, appError.Message)
}

// Unwrap returns the underlying error for error wrapping
func (appError *AppError) Unwrap() error {
	return appError.Cause
}

// Is allows error comparison using errors.Is()
func (appError *AppError) Is(target error) bool {
	if appErr, ok := target.(*AppError); ok {
		return appError.Code == appErr.Code
	}

	return false
}

// WithDetails adds additional details to an AppError
func (appError *AppError) WithDetails(details string) *AppError {
	appError.Details = details
	return appError
}

// WithDetailsf adds formatted details to an AppError
func (appError *AppError) WithDetailsf(format string, args ...any) *AppError {
	appError.Details = fmt.Sprintf(format, args...)
	return appError
}

// New creates a new AppError with explicit error code
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new AppError with formatted message
func Newf(code ErrorCode, format string, args ...any) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(code ErrorCode, cause error, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(code ErrorCode, cause error, format string, args ...any) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// FromError converts any error to an AppError
// If it's already an AppError, returns as-is
// Otherwise wraps it as an internal error
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return Wrap(ErrorCodeInternal, err, "internal error occurred")
}

// IsCode checks if an error is an AppError with a specific code
func IsCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}

	return false
}

// IsClientError checks if an error represents a client error (4xx)
func IsClientError(err error) bool {
	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}

	switch appErr.Code {
	case ErrorCodeBadRequest, ErrorCodeUnauthorized, ErrorCodeForbidden,
		ErrorCodeNotFound, ErrorCodeConflict, ErrorCodeValidation,
		ErrorCodeRateLimit, ErrorCodeInvalidCredentials, ErrorCodeTokenExpired,
		ErrorCodeTokenInvalid, ErrorCodeInsufficientScope, ErrorCodeConstraintViolation,
		ErrorCodeConcurrencyConflict:
		return true

	default:
		return false
	}
}

// IsServerError checks if an error represents a server error (5xx)
func IsServerError(err error) bool {
	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}

	return !IsClientError(appErr)
}
