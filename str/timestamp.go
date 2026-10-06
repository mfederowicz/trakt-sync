package str

import (
	"fmt"
	"time"
)

// Timestamp is a date and time from the Trakt API. A value decoded from a
// date without a time ("2006-01-02") is encoded as a date again.
type Timestamp struct {
	time.Time
	// DateOnly is set when the value was decoded from a date without a time
	DateOnly bool
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// UTC returns a copy of Timestamp with time converted to UTC.
func (t Timestamp) UTC() *Timestamp {
	return &Timestamp{Time: t.Time.UTC(), DateOnly: t.DateOnly}
}

// Define the possible formats
const (
	dateFormat     = "2006-01-02"
	dateTimeFormat = time.RFC3339 // "2006-01-02T15:04:05Z07:00"
	jsonNull       = "null"
	minStrLen      = 2
	start          = 1
)

// UnmarshalJSON supports both date and datetime formats
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	// A JSON null leaves the value as it is, like the standard library types do
	if string(b) == jsonNull {
		return nil
	}

	// Remove quotes from JSON string
	s := string(b)
	if len(s) >= minStrLen {
		s = s[start : len(s)-start] // Trim surrounding quotes
	}

	// Try parsing as full timestamp (RFC3339)
	parsedTime, err := time.Parse(dateTimeFormat, s)
	if err == nil {
		t.Time = parsedTime
		t.DateOnly = false
		return nil
	}

	// If that fails, try parsing as a simple date (YYYY-MM-DD)
	parsedTime, err = time.Parse(dateFormat, s)
	if err == nil {
		t.Time = parsedTime
		t.DateOnly = true
		return nil
	}

	return fmt.Errorf("invalid timestamp format: %s", s)
}

// MarshalJSON marshal json object to string
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// A date has no timezone: write the day it was decoded from, even when
	// the client moved the value to the user's timezone
	if t.DateOnly {
		return fmt.Appendf(nil, `"%s"`, t.Time.UTC().Format(dateFormat)), nil
	}

	// Otherwise, return full RFC3339 format
	return fmt.Appendf(nil, `"%s"`, t.Format(time.RFC3339)), nil
}
