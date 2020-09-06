package errors

import (
	"github.com/palantir/stacktrace"
	"github.com/pkg/errors"
)

// Customer facing errors
var (
	// ErrInternalServerError is thrown when there's a server error
	ErrInternalServerError = errors.New("internal server error")

	// ErrValidationError the input is invalid
	ErrValidationError = errors.New("input validation errors")

	// ErrInvalidJWTToken is thrown when the JWT token is invalid
	ErrInvalidJWTToken = errors.New("invalid JWT token")
)

var (
	// ErrCodeMissingJWT code when he jwt token is not present
	ErrCodeMissingJWT = stacktrace.ErrorCode(71)
)
