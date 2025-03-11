// Package consts used to store const for application
package consts

// usage strings
const (
	ConfigUsage            = "allow to overwrite default config filename"
	VersionUsage           = "get trakt-sync version"
	VerboseUsage           = "print additional verbose information"
	OutputUsage            = "allow to overwrite default output filename"
	FormatUsage            = "allow to overwrite default ID type format"
	TypeUsage              = "allow to overwrite type"
	ListUsage              = "allow to overwrite list"
	UserlistUsage          = "allow to export a user custom list"
	StartDateUsage         = "allow to overwrite start_date"
	DaysUsage              = "allow to overwrite days"
	ListIDUsage            = "allow to export a specific custom list"
	ListLikeRemoveUsage    = "allow remove like for list"
	ListCommentSortUsage   = "allow to overwrite comments sort"
	ModuleUsage            = "allow use selected module"
	ActionUsage            = "allow use selected action"
	TraktIDUsage           = "allow to overwrite trakt_id"
	CommentIDUsage         = "allow to overwrite comment_id"
	CheckInMsgUsage        = "allow to overwrite msg"
	DeleteUsage            = "allow delete item"
	RemoveUsage            = "allow remove item"
	SpoilerUsage           = "allow to overwrite spoiler"
	CommentUsage           = "allow to overwrite comment"
	CommentTypeUsage       = "allow to overwrite comment_type"
	IncludeRepliesUsage    = "allow to overwrite include_replies"
	ReplyUsage             = "allow to overwrite reply"
	EpisodeCodeUsage       = "episode_code format 01x24"
	EpisodeAbsUsage        = "episode_abs 1234"
	QueryUsage             = "allow use selected query"
	FieldUsage             = "allow use selected field"
	SortUsage              = "allow to overwrite sort"
	ExtendedInfoUsage      = "allow to overwrite extended flag"
	ActionTypeAll          = "all"
	NoShowTitle            = "no show title"
	NoEpisodeTitle         = "no episode title"
	EmptyPersonIDMsg       = "set personId ie: -i john-wayne"
	EmptyListIDMsg         = "set traktId ie: -trakt_id 55"
	EmptyTraktIDMsg        = "set traktId ie: -trakt_id 55"
	EmptyCommentIDMsg      = "set commentId ie: -comment_id 123"
	EmptyIncludeReplies    = "set includeReplies ie: -include_replies true or false"
	ErrorRender            = "error render: %w"
	CMD                    = "cmd"
	TestURL                = "test-url"
	TestURLNext            = "test-url-next"
	RatingRageMin          = 0
	RatingRangeMax         = 100
	RatingRageMinFloat     = 0.00
	RatingRangeMaxFloat    = 100.00
	VotesRangeMin          = 0
	VotesRangeMax          = 100000
	ImdbVotesRangeMin      = 0
	ImdbVotesRangeMax      = 3000000
	TmdbRatingRangeMin     = 0.00
	TmdbRatingRangeMax     = 10.00
	X644                   = 0644
	X755                   = 0755
	NextPageStep           = 1
	SleepNumberOfSeconds   = 2
	OneValue               = 1
	BaseInt                = 10
	BitSize                = 64
	MaxAcceptedStatus      = 299
	DefaultPerPage         = 100
	PagesLimit             = 2
	PerPage                = 50
	ZeroValue              = 0
	TwoValue               = 2
	MinSeasonNumberLength  = 3
	DefaultExitCode        = 0
	FirstArgElement        = 0
	DefaultPage            = 1
	EmptyResult            = "empty result"
	RangeFormatDigits      = "%d-%d"
	StringDigit            = "%s%d"
	ErrorsPlaceholders     = "%v %v: %d %v"
	RangeFormatFloats      = "%.1f-%.1f"
	EmptyString            = ""
	CommaString            = ","
	SeparatorString        = CommaString
	JSONDataFormat         = "  "
	DefaultStartDateFormat = "2006-01-02T15:00Z"
	EpisodesType           = "episodes"
	TmdbFormat             = "tmdb"
	TvdbFormat             = "tvdb"
	ImdbFormat             = "imdb"
	TmdbIDFormat           = "Tmdb"
	TvdbIDFormat           = "Tvdb"
	ImdbIDFormat           = "Imdb"
	DefaultOutputFormat1   = "export_%s.json"
	DefaultOutputFormat2   = "export_%s_%s.json"
	DefaultOutputFormat3   = "export_%s_%s_%s.json"
	NewLine                = "\n"
	Trending               = "trending"
	Recent                 = "recent"
	EmptyBuildInfoLen      = 0
)
