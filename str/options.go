// Package str used for structs
package str

// Options represents a app opions.
type Options struct {
	Headers      map[string]any
	ExtendedInfo string
	List         string
	Type         string
	SearchIDType string
	StartDate    string
	Query        string
	Sort         string
	CommentsSort string
	Module       string
	Action       string
	Format       string
	UserName     string
	Time         string
	ID           string
	Output       string
	SearchType   Slice
	SearchField  Slice
	Token        Token
	PerPage      int
	Days         int
	TraktID      int
	Msg          string
	Verbose      bool
	Version      bool
	Remove       bool
}
