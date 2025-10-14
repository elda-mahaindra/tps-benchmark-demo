package errors

import (
	"net/http"
)

// HTTPErrorResponse converts any error to an HTTP response map
// Returns status code and response body
func HTTPErrorResponse(err error) (int, map[string]any) {
	if err == nil {
		return http.StatusOK, nil
	}

	appErr := FromError(err)

	status := appErr.ToHTTPStatus()
	response := map[string]any{
		"error": map[string]any{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	}

	return status, response
}

// ToGRPCError converts any error to a gRPC error
// Uses FromError() internally to handle both AppError and regular errors
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	appErr := FromError(err)

	return appErr.ToGRPCStatus().Err()
}
