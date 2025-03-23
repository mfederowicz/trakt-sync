// Package str used for structs
package str

import (
	"fmt"
	"time"
)

// Timestamp object
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// Define the possible formats
const (
	dateFormat     = "2006-01-02"
	dateTimeFormat = time.RFC3339 // "2006-01-02T15:04:05Z07:00"
	minStrLen      = 2
	start          = 1
	zero           = 0
)

// UnmarshalJSON supports both date and datetime formats
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	// Remove quotes from JSON string
	s := string(b)
	if len(s) >= minStrLen {
		s = s[start : len(s)-start] // Trim surrounding quotes
	}

	// Try parsing as full timestamp (RFC3339)
	parsedTime, err := time.Parse(dateTimeFormat, s)
	if err == nil {
		t.Time = parsedTime
		return nil
	}

	// If that fails, try parsing as a simple date (YYYY-MM-DD)
	parsedTime, err = time.Parse(dateFormat, s)
	if err == nil {
		t.Time = parsedTime
		return nil
	}

	return fmt.Errorf("invalid timestamp format: %s", s)
}

// MarshalJSON marshal json object to string
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// If time is exactly midnight (00:00:00 UTC), assume it was a date-only input
	if t.Hour() == zero && t.Minute() == zero && t.Second() == zero {
		return fmt.Appendf(nil, `"%s"`, t.Format("2006-01-02")), nil
	}

	// Otherwise, return full RFC3339 format
	return fmt.Appendf(nil, `"%s"`, t.Format(time.RFC3339)), nil
}
