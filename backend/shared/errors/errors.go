package errors

import "github.com/palantir/stacktrace"

var (
	// ErrCodeDatabaseError is thrown when there's an error with the database
	ErrCodeDatabaseError = stacktrace.ErrorCode(1)

	// ErrCodeEntityNotFound is thrown when there's is no entity in the database
	ErrCodeEntityNotFound = stacktrace.ErrorCode(2)

	// ErrCodeCannotDecodeIDFromString is thrown when the ID cannot be decoded from a string
	ErrCodeCannotDecodeIDFromString = stacktrace.ErrorCode(3)

	// ErrCodeCannotDecodeIDFromInterface is thrown when the ID cannot be decoded from a string
	ErrCodeCannotDecodeIDFromInterface = stacktrace.ErrorCode(4)

	// ErrCodeDatabaseHydrationError is thrown when there is an error when hydrating the database
	ErrCodeDatabaseHydrationError = stacktrace.ErrorCode(5)

	// ErrCodeInvalidRawRecordSource when the raw record source is invalid
	ErrCodeInvalidRawRecordSource = stacktrace.ErrorCode(6)
)

var (
	// ErrEntityNotFound is returned when an entity does not exist in the database
	ErrEntityNotFound = stacktrace.NewErrorWithCode(ErrCodeEntityNotFound, "entity not found")

	// ErrCannotDecodeIDFromString is thrown when the ID cannot be decoded from a string
	ErrCannotDecodeIDFromString = stacktrace.NewErrorWithCode(ErrCodeCannotDecodeIDFromString, "cannot decode uuid from string")

	// ErrCannotDecodeIDFromInterface is thrown when the ID cannot be decoded from a string
	ErrCannotDecodeIDFromInterface = stacktrace.NewErrorWithCode(ErrCodeCannotDecodeIDFromInterface, "cannot decode uuid from interface")

	// ErrInvalidRawRecordSource is thrown when the raw record source is invalid
	ErrInvalidRawRecordSource = stacktrace.NewErrorWithCode(ErrCodeInvalidRawRecordSource, "raw record source is invalid")

)
