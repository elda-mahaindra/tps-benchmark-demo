package errs

const (
	CODE_SUCCESS             = "00" // Success.
	CODE_ERR_UNANTICIPATED   = "99" // Unanticipated error.
	CODE_ERR_DATABASE        = "98" // Error from database.
	CODE_ERR_IO              = "97" // External I/O error such as network failure.
	CODE_ERR_VALIDATION      = "96" // Input validation error.
	CODE_ERR_NOT_EXIST       = "95" // Item does not exist.
	CODE_ERR_UNAUTHENTICATED = "94" // Unauthenticated request
	CODE_ERR_UNAUTHORIZED    = "93" // Authenticated but unauthorized request.
)
