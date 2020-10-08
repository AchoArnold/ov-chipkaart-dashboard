package types

import (
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errors"
	"github.com/palantir/stacktrace"
)

const (
	// RawRecordSourceAPI is an API source
	RawRecordSourceAPI = RawRecordSource("API")

	// RawRecordSourceCSV is a CSV source
	RawRecordSourceCSV = RawRecordSource("CSV")
)

// RawRecordSource is the source for a raw record entry
type RawRecordSource string

// String converts a raw record source to a string
func (source RawRecordSource) String() string {
	return string(source)
}

func RawRecordSourceFromString(source string) (rawRecordSource RawRecordSource,err error) {
	if source != RawRecordSourceAPI.String() && source != RawRecordSourceCSV.String() {
		return rawRecordSource, stacktrace.Propagate(errors.ErrInvalidRawRecordSource, source)
	}
	return RawRecordSource(source), nil
}
