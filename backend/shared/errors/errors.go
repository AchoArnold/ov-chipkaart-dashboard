package errors

import "github.com/palantir/stacktrace"

var (
	// ErrCodeDatabaseError is thrown when there's an error with the database
	ErrCodeDatabaseError = stacktrace.ErrorCode(1)

	// ErrCodeEntityNotFound is thrown when there's is no entity in the database
	ErrCodeEntityNotFound = stacktrace.ErrorCode(2)
)

var (
	// ErrEntityNotFound is returned when an entity does not exist in the database
	ErrEntityNotFound = stacktrace.NewErrorWithCode(ErrCodeEntityNotFound, "entity not found")
)
