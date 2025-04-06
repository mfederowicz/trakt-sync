// Package str used for structs
package str

// Options represents a app opions.
type Options struct {
	Action            string
	Comment           string
	CommentID         int
	CommentType       string
	CommentsSort      string
	Country           string
	Days              int
	Delete            bool
	Episode           int
	EpisodeAbs        int
	EpisodeCode       string
	ExtendedInfo      string
	Format            string
	Headers           map[string]any
	Hide              bool
	ID                string
	IncludeReplies    string
	IgnoreCollected   string
	IgnoreWatchlisted string
	InternalID        string
	Item              string
	Language          string
	List              string
	Module            string
	Msg               string
	Notes             string
	Output            string
	PagesLimit        int
	PerPage           int
	Period            string
	Privacy           string
	Query             string
	Remove            bool
	Reply             string
	SearchField       Slice
	SearchIDType      string
	SearchType        Slice
	Season            int
	Sort              string
	Spoiler           bool
	StartDate         string
	Time              string
	Token             Token
	TraktID           int
	Type              string
	UserName          string
	Verbose           bool
	Version           bool
}
