package time

import "time"

const (
	// DefaultFormat is the default format for timestamps
	DefaultFormat = "2006-01-02 15:04:05"

	// DateFormat is the default format for dates
	DateFormat = "2006-01-02"

	// StringFormat is the format when time is converted to string
	StringFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

// FromDate returns a time from a date string
func FromDate(date string) (time.Time, error) {
	return time.Parse(DateFormat, date)
}
