package str
// Trakt config.
type Options struct {
	Headers      map[string]any
	ExtendedInfo string
	List         string
	Type         string
	SearchIdType string
	StartDate    string
	Query        string
	Sort         string
	Module       string
	Action       string
	Format       string
	UserName     string
	Time         string
	Id           string
	Output       string
	SearchType   StrSlice
	SearchField  StrSlice
	Token        Token
	PerPage      int
	Days         int
	Verbose      bool
	Version      bool
}

