// Package str used for structs
package str

// ListReportResult represents JSON list report result object
type ListReportResult struct {
	Message *string `json:"message,omitempty"`
}

func (l ListReportResult) String() string {
	return Stringify(l)
}
