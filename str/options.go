// Package str used for structs
package str

// Options represents a app opions.
type Options struct {
	Headers        map[string]any
	ExtendedInfo   string
	List           string
	Type           string
	SearchIDType   string
	StartDate      string
	Query          string
	Sort           string
	Period         string
	CommentsSort   string
	CommentType    string
	Module         string
	Action         string
	Format         string
	UserName       string
	Time           string
	ID             string
	InternalID     string
	Language       string
	Output         string
	SearchType     Slice
	SearchField    Slice
	Token          Token
	PerPage        int
	PagesLimit     int
	Days           int
	TraktID        int
	CommentID      int
	Msg            string
	Comment        string
	Reply          string
	Verbose        bool
	Version        bool
	Remove         bool
	Delete         bool
	Spoiler        bool
	IncludeReplies string
	EpisodeAbs     int
	EpisodeCode    string
	Country        string
	Episode        int
	Season         int
	Notes          string
}
