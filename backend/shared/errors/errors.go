package errors

import "github.com/palantir/stacktrace"

var (
	// ErrCodeDatabaseError is thrown when there's an error with the database
	ErrCodeDatabaseError = stacktrace.ErrorCode(1)
)
