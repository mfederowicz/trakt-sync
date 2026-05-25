// Package str used for structs
package str

// UserReportResult represents JSON user report result object
type UserReportResult struct {
	Message *string `json:"message,omitempty"`
}

func (u UserReportResult) String() string {
	return Stringify(u)
}
