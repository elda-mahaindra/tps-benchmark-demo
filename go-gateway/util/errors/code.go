package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents semantic error categories in our system
// These are interface-agnostic and represent business/domain errors
type ErrorCode string

const (
	// Client Errors - 4xx category
	ErrorCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrorCodeConflict     ErrorCode = "CONFLICT"
	ErrorCodeValidation   ErrorCode = "VALIDATION_FAILED"
	ErrorCodeRateLimit    ErrorCode = "RATE_LIMIT_EXCEEDED"

	// Server Errors - 5xx category
	ErrorCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrorCodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
	ErrorCodeUnavailable    ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeTimeout        ErrorCode = "TIMEOUT"
	ErrorCodeDependency     ErrorCode = "DEPENDENCY_FAILED"

	// Auth-specific errors
	ErrorCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrorCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
	ErrorCodeInsufficientScope  ErrorCode = "INSUFFICIENT_SCOPE"

	// Data-specific errors
	ErrorCodeDataCorrupted       ErrorCode = "DATA_CORRUPTED"
	ErrorCodeConstraintViolation ErrorCode = "CONSTRAINT_VIOLATION"
	ErrorCodeConcurrencyConflict ErrorCode = "CONCURRENCY_CONFLICT"
)

// ErrorCodeToHTTPStatus maps our error codes to HTTP status codes
func ErrorCodeToHTTPStatus(code ErrorCode) int {
	switch code {
	// 4xx Client Errors
	case ErrorCodeBadRequest, ErrorCodeValidation:
		return http.StatusBadRequest
	case ErrorCodeUnauthorized, ErrorCodeInvalidCredentials, ErrorCodeTokenExpired, ErrorCodeTokenInvalid:
		return http.StatusUnauthorized
	case ErrorCodeForbidden, ErrorCodeInsufficientScope:
		return http.StatusForbidden
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeConflict, ErrorCodeConstraintViolation, ErrorCodeConcurrencyConflict:
		return http.StatusConflict
	case ErrorCodeRateLimit:
		return http.StatusTooManyRequests

	// 5xx Server Errors
	case ErrorCodeInternal, ErrorCodeDataCorrupted:
		return http.StatusInternalServerError
	case ErrorCodeNotImplemented:
		return http.StatusNotImplemented
	case ErrorCodeUnavailable, ErrorCodeDependency:
		return http.StatusServiceUnavailable
	case ErrorCodeTimeout:
		return http.StatusGatewayTimeout

	// Default to internal server error for unknown codes
	default:
		return http.StatusInternalServerError
	}
}

// ErrorCodeToGRPCStatus maps our error codes to gRPC status codes
func ErrorCodeToGRPCStatus(code ErrorCode, message string) *status.Status {
	var grpcCode codes.Code

	switch code {
	// Client Errors
	case ErrorCodeBadRequest, ErrorCodeValidation:
		grpcCode = codes.InvalidArgument
	case ErrorCodeUnauthorized, ErrorCodeInvalidCredentials, ErrorCodeTokenExpired, ErrorCodeTokenInvalid:
		grpcCode = codes.Unauthenticated
	case ErrorCodeForbidden, ErrorCodeInsufficientScope:
		grpcCode = codes.PermissionDenied
	case ErrorCodeNotFound:
		grpcCode = codes.NotFound
	case ErrorCodeConflict, ErrorCodeConstraintViolation, ErrorCodeConcurrencyConflict:
		grpcCode = codes.AlreadyExists
	case ErrorCodeRateLimit:
		grpcCode = codes.ResourceExhausted

	// Server Errors
	case ErrorCodeInternal, ErrorCodeDataCorrupted:
		grpcCode = codes.Internal
	case ErrorCodeNotImplemented:
		grpcCode = codes.Unimplemented
	case ErrorCodeUnavailable, ErrorCodeDependency:
		grpcCode = codes.Unavailable
	case ErrorCodeTimeout:
		grpcCode = codes.DeadlineExceeded

	// Default to internal error
	default:
		grpcCode = codes.Internal
	}

	return status.New(grpcCode, message)
}

// HTTPStatusToErrorCode converts HTTP status codes back to our ErrorCode
// This is useful for adapter layers and when calling external HTTP services
func HTTPStatusToErrorCode(httpStatus int) ErrorCode {
	switch httpStatus {
	// 4xx Client Errors
	case http.StatusBadRequest:
		return ErrorCodeBadRequest
	case http.StatusUnauthorized:
		return ErrorCodeUnauthorized
	case http.StatusForbidden:
		return ErrorCodeForbidden
	case http.StatusNotFound:
		return ErrorCodeNotFound
	case http.StatusConflict:
		return ErrorCodeConflict
	case http.StatusTooManyRequests:
		return ErrorCodeRateLimit

	// 5xx Server Errors
	case http.StatusInternalServerError:
		return ErrorCodeInternal
	case http.StatusNotImplemented:
		return ErrorCodeNotImplemented
	case http.StatusServiceUnavailable:
		return ErrorCodeUnavailable
	case http.StatusGatewayTimeout:
		return ErrorCodeTimeout

	// Default mappings for less common status codes
	case http.StatusUnprocessableEntity:
		return ErrorCodeValidation
	case http.StatusRequestTimeout:
		return ErrorCodeTimeout
	case http.StatusBadGateway:
		return ErrorCodeDependency

	// Default to internal error for unknown status codes
	default:
		if httpStatus >= 400 && httpStatus < 500 {
			return ErrorCodeBadRequest
		}
		return ErrorCodeInternal
	}
}

// GRPCStatusToErrorCode converts gRPC status codes back to our ErrorCode
// This is useful for gRPC client calls and inter-service communication
func GRPCStatusToErrorCode(grpcCode codes.Code) ErrorCode {
	switch grpcCode {
	// Client Errors
	case codes.InvalidArgument:
		return ErrorCodeBadRequest
	case codes.Unauthenticated:
		return ErrorCodeUnauthorized
	case codes.PermissionDenied:
		return ErrorCodeForbidden
	case codes.NotFound:
		return ErrorCodeNotFound
	case codes.AlreadyExists:
		return ErrorCodeConflict
	case codes.ResourceExhausted:
		return ErrorCodeRateLimit

	// Server Errors
	case codes.Internal:
		return ErrorCodeInternal
	case codes.Unimplemented:
		return ErrorCodeNotImplemented
	case codes.Unavailable:
		return ErrorCodeUnavailable
	case codes.DeadlineExceeded:
		return ErrorCodeTimeout

	// Additional mappings
	case codes.FailedPrecondition:
		return ErrorCodeValidation
	case codes.Aborted:
		return ErrorCodeConflict
	case codes.OutOfRange:
		return ErrorCodeBadRequest
	case codes.DataLoss:
		return ErrorCodeDataCorrupted

	// Default to internal error for unknown codes
	default:
		return ErrorCodeInternal
	}
}

// IsRetryable determines if an error with this code should be retried
// This provides a centralized way to classify errors for retry logic
func (code ErrorCode) IsRetryable() bool {
	switch code {
	// Client errors - never retry (4xx equivalent)
	// These represent permanent failures that won't be fixed by retrying
	case ErrorCodeBadRequest, ErrorCodeUnauthorized, ErrorCodeForbidden,
		ErrorCodeNotFound, ErrorCodeConflict, ErrorCodeValidation,
		ErrorCodeInvalidCredentials, ErrorCodeTokenExpired,
		ErrorCodeTokenInvalid, ErrorCodeInsufficientScope,
		ErrorCodeConstraintViolation, ErrorCodeConcurrencyConflict:
		return false

	// Server/infrastructure errors - should retry (5xx equivalent)
	// These represent temporary failures that might be resolved by retrying
	case ErrorCodeInternal, ErrorCodeNotImplemented, ErrorCodeUnavailable,
		ErrorCodeTimeout, ErrorCodeDependency, ErrorCodeDataCorrupted:
		return true

	// Rate limiting - retryable with backoff
	case ErrorCodeRateLimit:
		return true

	// Conservative default - don't retry unknown error codes
	default:
		return false
	}
}
