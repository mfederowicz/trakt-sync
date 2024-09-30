// Package uri used for url operations
package uri

// Pagination represents pagination params
type Pagination struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`
	// For paginated result sets, the number of elements on one page.
	Limit int `url:"limit,omitempty"`
}

