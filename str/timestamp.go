// Package str used for structs
package str

import (
	"time"
)

// Timestamp object
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

