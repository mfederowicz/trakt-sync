// Package cmds used for commands modules
package cmds

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

var (
	_userName     = flag.String("u", cfg.DefaultConfig().UserName, consts.UserlistUsage)
	_strType      = flag.String("t", cfg.DefaultConfig().Type, consts.TypeUsage)
	_output       = flag.String("o", cfg.DefaultConfig().Output, consts.OutputUsage)
	_format       = flag.String("f", cfg.DefaultConfig().Format, consts.FormatUsage)
	_extendedInfo = flag.String("ex", "", consts.ExtendedInfoUsage)
	_query        = flag.String("query", "", "")
	_years        = flag.String("years", "", "")
	_genres       = flag.String("genres", "", "")
	_languages    = flag.String("languages", "", "")
	_countries    = flag.String("countries", "", "")
	_runtimes     = flag.String("runtimes", "", "")
	_studioIDs    = flag.String("studio_ids", "", "")
)

// Avflags contains all available flags
var Avflags = map[string]bool{
	"a":               true,
	"c":               true,
	"calendars":       true,
	"certifications":  true,
	"comments":        true,
	"comment_type":    true,
	"countries":       true,
	"country":         true,
	"checkin":         true,
	"collection":      true,
	"days":            true,
	"delete":          true,
	"ex":              true,
	"f":               true,
	"field":           true,
	"godoc":           true,
	"genres":          true,
	"help":            true,
	"history":         true,
	"i":               true,
	"include_replies": true,
	"trakt_id":        true,
	"comment_id":      true,
	"episode_code":    true,
	"episode_abs":     true,
	"id_type":         true,
	"language":        true,
	"languages":       true,
	"lists":           true,
	"msg":             true,
	"movies":          true,
	"networks":        true,
	"o":               true,
	"people":          true,
	"period":          true,
	"q":               true,
	"remove":          true,
	"releases":        true,
	"reply":           true,
	"s":               true,
	"comment":         true,
	"search":          true,
	"start_date":      true,
	"t":               true,
	"translations":    true,
	"u":               true,
	"users":           true,
	"v":               true,
	"version":         true,
	"watchlist":       true,
}

type fatal struct{}

// A Command represents a subcommand of trakt-sync.
type Command struct {
	Flag    flag.FlagSet
	Run     func(cmd *Command, args ...string) error
	Client  *internal.Client
	Config  *cfg.Config
	Options *str.Options
	Name    string
	Usage   string
	Summary string
	Help    string
	Abbrev  string
	exit    int
}

// UpdateMovieFlagsValues update movies flags values only in command
func (*Command) UpdateMovieFlagsValues() {
	if *_moviesSort == "" {
		switch *_moviesAction {
		case "comments":
			*_moviesSort = "newest"
		case "lists":
			*_moviesSort = "popular"
		}
	}
	if *_moviesType == "" {
		switch *_moviesAction {
		case "comments":
			*_moviesType = ""
		case "lists":
			*_moviesType = "personal"
		}
	}
}

// ValidPeriodForModule valid period options depends on action value
func (c *Command) ValidPeriodForModule(options *str.Options) error {
	switch options.Action {
	case "favorited":
		err := c.ValidPeriod(options)
		if err != nil {
			return errors.New(err.Error())
		}
	case "played":
		err := c.ValidPeriod(options)
		if err != nil {
			return errors.New(err.Error())
		}
	case "watched":
		err := c.ValidPeriod(options)
		if err != nil {
			return errors.New(err.Error())
		}
	case "collected":
		err := c.ValidPeriod(options)
		if err != nil {
			return errors.New(err.Error())
		}
	}
	return nil
}

// Helper function to handle the error
func handleHelpError(err error) {
	if err != nil {
		printer.Printf("error:%s", err)
	}
}

// Exec core command function
func (c *Command) Exec(fs afero.Fs, client *internal.Client, config *cfg.Config, args []string) error {
	c.Client = client
	c.Config = config
	c.Flag.Usage = func() {
		handleHelpError(HelpFunc(c, c.Name))
	}
	c.registerGlobalFlagsInSet(&c.Flag)
	_ = c.Flag.Parse(args)
	m := c.fetchFlagsMap()
	options, err := cfg.SyncOptionsFromFlags(fs, c.Config, m)

	if err != nil {
		return fmt.Errorf("%s", err)
	}

	options.Type = *_strType

	options.Module = c.Name
	options = setOptionsDependsOnModule(c.Name, options)
	c.Options = &options

	if !c.ValidFlags() {
		return errors.New("invalid flags")
	}

	processVerbose(&options)

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(fatal); ok {
				err = errors.New("fatal error")
			} else {
				err = fmt.Errorf("panic error:%s", r)
			}
		}
	}()

	err = c.Run(c, c.Flag.Args()...)

	return err
}

func processVerbose(options *str.Options) {
	if options.Verbose {
		printer.Println("Authorization header:" + options.Headers["Authorization"].(string))
		printer.Println("trakt-api-key header:" + options.Headers["trakt-api-key"].(string))
		printer.Println("token expiration in seconds:" + strconv.Itoa(options.Token.ExpiritySeconds()))
		printer.Println("Extended info:" + *_extendedInfo)
		if len(options.Module) > consts.ZeroValue {
			printer.Println("selected module:" + options.Module)
		}
		printer.Println(
			str.Format("selected user: {0}, module: {1}, type: {2}, per_page: {3}, format: {4}, action: {5}, sort: {6}",
				options.UserName, options.Module, options.Type, options.PerPage, options.Format, options.Action, options.Sort),
		)
	}
}

func selectFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if len(v) > consts.ZeroValue {
			return v
		}
	}
	return ""
}

func setOptionsDependsOnModule(module string, options str.Options) str.Options {
	switch module {
	case "comments":
		options.Action = *_commentsAction
		options.InternalID = selectFirstNonEmpty(*_commentsTraktID, *_commentsInternalID)
		options.CommentID = *_commentsCommentID
		options.CommentType = *_commentsCommentType
	case "checkin":
		options.Action = *_checkinAction
		options.TraktID = *_checkinTraktID
	case "lists":
		options.Action = *_listsAction
		options.InternalID = selectFirstNonEmpty(*_listTraktID, *_listInternalID)
		options.Sort = *_listSort
	case "movies":
		options.Action = *_moviesAction
		options.Period = *_moviesPeriod
		options.StartDate = *_moviesStartDate
		options.InternalID = *_moviesInternalID
		options.Country = *_moviesCountry
		options.Language = *_moviesLanguage
		options.Sort = *_moviesSort
		options.Type = *_moviesType
	case "networks":
		options.Action = *_networksAction
	case "users":
		options.Action = *_usersAction
	case "people":
		options.Action = *_action
		options.ID = *_personID
		options.Type = *_action
	case "calendars":
		options.Action = *_calAction
		options.StartDate = *_calStartDate
		options.Days = *_calDays
	case "search":
		options.Action = *_searchAction
		options.SearchType = _searchType
		options.SearchField = _searchField
		options.ID = *_searchID
		options.SearchIDType = *_searchIDType
	case "watchlist":
	case "collection":
	case "history":
		options.Format = *_format
	}

	return options
}

// BadArgs shows error if command have invalid arguments
func (c *Command) BadArgs(errFormat string, args ...any) {
	printer.Fprintf(stdout, "error: "+errFormat+"\n\n", args...)
	err := HelpFunc(c, c.Name)
	if err != nil {
		printer.Printf("error:%s", err)
	}
}

// Errorf prints out a formatted error with the right prefixes.
func (c *Command) Errorf(errFormat string, args ...any) {
	printer.Fprintf(stdout, c.Name+": error: "+errFormat+"\n", args...)
	if c.exit == consts.ZeroValue {
		c.exit = consts.DefaultExitCode
	}
}

// Fatalf is like Errorf except the stack unwinds up to the Exec call before
// exiting the application with status code 1.
func (c *Command) Fatalf(errFormat string, args ...any) {
	c.Errorf(errFormat, args...)
	panic(fatal{})
}

func (c *Command) fetchFlagsMap() map[string]string {
	flagMap := make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		flagMap[f.Name] = f.Value.String()
	})
	c.Flag.VisitAll(func(f *flag.Flag) {
		flagMap[f.Name] = f.Value.String()
	})

	return flagMap
}

func cleanKey(arg string) string {
	return strings.TrimLeft(arg, "-")
}
func processArgsItem(arg string, key string, argMap map[string]bool) (string, map[string]bool) {
	// If the argument starts with "-", consider it a key
	if arg[consts.FirstArgElement] == '-' {
		// If we already have a key, it means it's a single argument without a value
		if key != "" {
			argMap[cleanKey(key)] = true // Set the key to true for bool map
		}
		key = arg
	} else {
		// If we have a key, assign the value to it
		if key != "" {
			argMap[cleanKey(key)] = true // Set the key to true for bool map
			key = consts.EmptyString
		} else {
			// If we don't have a key, consider it a standalone argument
			argMap[arg] = true // Set the key to true for bool map
		}
	}

	return key, argMap
}
func argsToMap(args []string) map[string]bool {
	argMap := make(map[string]bool)
	var key string

	for _, arg := range args {
		key, argMap = processArgsItem(arg, key, argMap)
	}

	// If we still have a key at the end, it means it's a single argument without a value
	if key != consts.EmptyString {
		argMap[cleanKey(key)] = true // Set the key to true for bool map
	}

	return argMap
}

// ValidFlags validate if flag is in our list
func (*Command) ValidFlags() bool {
	for f := range argsToMap(flag.Args()) {
		if _, ok := Avflags[f]; !ok {
			return false
		}
	}
	return true
}

func (*Command) registerGlobalFlagsInSet(fset *flag.FlagSet) {
	flag.VisitAll(func(f *flag.Flag) {
		if fset.Lookup(f.Name) == nil {
			fset.Var(f.Value, f.Name, f.Usage)
		}
	})
}

// IsImdbMovie check movie imdb format
func (*Command) IsImdbMovie(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != consts.EpisodesType && data.Movie != nil && options.Format == consts.ImdbFormat
}

// IsImdbShow check show imdb format
func (*Command) IsImdbShow(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != consts.EpisodesType && data.Show != nil && data.Show.IDs.HaveID(consts.ImdbIDFormat) &&
		options.Format == consts.ImdbFormat
}

// IsImdbEpisode check episode imdb format
func (*Command) IsImdbEpisode(options *str.Options, data *str.ExportlistItem) bool {
	return data.Episode != nil && data.Episode.IDs.HaveID(consts.ImdbIDFormat) && options.Format == consts.ImdbFormat
}

// IsTmdbMovie check movie tmdb format
func (*Command) IsTmdbMovie(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != consts.EpisodesType && data.Movie != nil &&
		data.Movie.IDs.HaveID(consts.TmdbIDFormat) && options.Format == consts.TmdbFormat
}

// IsTmdbShow check show tmdb format
func (*Command) IsTmdbShow(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != consts.EpisodesType && data.Show != nil &&
		data.Show.IDs.HaveID(consts.TmdbIDFormat) && options.Format == consts.TmdbFormat
}

// IsTmdbEpisode check episode tmdb format
func (*Command) IsTmdbEpisode(options *str.Options, data *str.ExportlistItem) bool {
	return data.Episode.IDs.HaveID(consts.TmdbIDFormat) && options.Format == consts.TmdbFormat
}

// IsTvdbEpisode check episode tvdb format
func (*Command) IsTvdbEpisode(options *str.Options, data *str.ExportlistItem) bool {
	return data.Episode.IDs.HaveID("Tvdb") && options.Format == "tvdb"
}

// ExportListProcess process list items
func (c *Command) ExportListProcess(
	data *str.ExportlistItem, options *str.Options,
	findDuplicates []any, exportJSON []str.ExportlistItemJSON,
) ([]any, []str.ExportlistItemJSON, error) {
	var handler handlers.ItemsHandler
	switch {
	case c.IsImdbMovie(options, data):
		handler = handlers.ImdbMovieHandler{}
	case c.IsImdbShow(options, data):
		handler = handlers.ImdbShowHandler{}
	case c.IsImdbEpisode(options, data):
		handler = handlers.ImdbEpisodeHandler{}
	case c.IsTmdbMovie(options, data):
		handler = handlers.TmdbMovieHandler{}
	case c.IsTmdbShow(options, data):
		handler = handlers.TmdbShowHandler{}
	case c.IsTmdbEpisode(options, data):
		handler = handlers.TmdbEpisodeHandler{}
	case c.IsTvdbEpisode(options, data):
		handler = handlers.TvdbEpisodeHandler{}
	default:
		handler = handlers.DefaultHandler{}
	}

	return handler.Handle(options, data, findDuplicates, exportJSON)
}

// PrepareQueryString for remove or replace unwanted signs from query string
func (*Command) PrepareQueryString(q string) *string {
	return &q
}

// ValidModuleType check if type is valid
func (*Command) ValidModuleType(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	// Check if the provided type is valid for the selected module
	if !cfg.IsValidConfigType(moduleConfig.Type, options.Type) {
		return fmt.Errorf("type '%s' is not valid for module '%s' and action '%s', avaliable types:%s", options.Type, options.Module, options.Action, moduleConfig.Type)
	}

	return nil
}

// ValidModuleActionType check if type is valid
func (*Command) ValidModuleActionType(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	_, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	
	// Check if the provided type is valid for the selected module
	prefix := options.Module + ":" + options.Action
	if len(cfg.ModuleActionConfig[prefix].Type) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.ModuleActionConfig[prefix].Type, options.Type) {
		return fmt.Errorf("type '%s' is not valid for module '%s' and action '%s', avaliable type:%s", options.Type, options.Module, options.Action, cfg.ModuleActionConfig[prefix].Type)
	}

	return nil
}

// ValidSort check if sort is valid
func (*Command) ValidSort(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	_, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	// Check if the provided sort is valid for the selected module
	prefix := options.Module + ":" + options.Action
	if len(cfg.ModuleActionConfig[prefix].Sort) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.ModuleActionConfig[prefix].Sort, options.Sort) {
		return fmt.Errorf("sort '%s' is not valid for module '%s' and action '%s', avaliable sort:%s", options.Sort, options.Module, options.Action, cfg.ModuleActionConfig[prefix].Sort)
	}

	return nil
}

// ValidPeriod check if period is valid
func (*Command) ValidPeriod(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	// Check if the provided period is valid for the selected module
	if !cfg.IsValidConfigType(moduleConfig.Period, options.Period) {
		return fmt.Errorf("period '%s' is not valid for module '%s' and action '%s', avaliable periods:%s", options.Period, options.Module, options.Action, moduleConfig.Period)
	}

	return nil
}

// UpdateOptionsWithCommandFlags update options depends on command flags
func (c *Command) UpdateOptionsWithCommandFlags(options *str.Options) *str.Options {
	if len(*_userName) > consts.ZeroValue {
		options.UserName = *_userName
	}

	if len(*_searchQuery) > consts.ZeroValue {
		options.Query = *c.PrepareQueryString(*_searchQuery)
	}

	if len(*_extendedInfo) > consts.ZeroValue {
		options.ExtendedInfo = *_extendedInfo
	}

	if len(*_format) > consts.ZeroValue {
		options.Format = *_format
	}

	if len(*_output) > consts.ZeroValue {
		options.Output = *_output
	} else {
		options.Output = cfg.GetOutputForModule(options)
	}

	if len(*_startDate) > consts.ZeroValue {
		options.StartDate = convertDateString(*_startDate, consts.DefaultStartDateFormat)
	} else {
		options.StartDate = time.Now().Format(consts.DefaultStartDateFormat)
	}

	if len(*_usersListID) > consts.ZeroValue {
		options.ID = *_usersListID
	}

	if len(*_listTraktID) > consts.ZeroValue || len(*_listInternalID) > consts.ZeroValue {
		options.InternalID = selectFirstNonEmpty(*_listTraktID, *_listInternalID)
	}

	if *_listLikeRemove {
		options.Remove = *_listLikeRemove
	}

	if *_commentsDelete {
		options.Delete = *_commentsDelete
	}

	if *_commentsRemove {
		options.Remove = *_commentsRemove
	}

	if len(*_listSort) > consts.ZeroValue {
		options.CommentsSort = *_listSort
	}

	if len(*_checkinMsg) > consts.ZeroValue {
		options.Msg = *_checkinMsg
	}

	if *_checkinEpisodeAbs > consts.ZeroValue {
		options.EpisodeAbs = *_checkinEpisodeAbs
	}

	if len(*_checkinEpisodeCode) > consts.ZeroValue {
		options.EpisodeCode = *_checkinEpisodeCode
	}

	if len(*_commentsComment) > consts.ZeroValue {
		options.Comment = *_commentsComment
	}

	if *_commentsCommentID > consts.ZeroValue {
		options.CommentID = *_commentsCommentID
	}

	if len(*_commentsReply) > consts.ZeroValue {
		options.Reply = *_commentsReply
	}

	if len(*_commentsIncludeReplies) > consts.ZeroValue {
		options.IncludeReplies = *_commentsIncludeReplies
	}

	if len(*_moviesAction) > consts.ZeroValue {
		options.Action = *_moviesAction
	}

	if len(*_moviesPeriod) > consts.ZeroValue {
		options.Period = *_moviesPeriod
	}

	if len(*_moviesStartDate) > consts.ZeroValue {
		options.StartDate = convertDateString(*_moviesStartDate, consts.DefaultStartDateFormat)
	} else {
		options.StartDate = time.Now().Format(consts.DefaultStartDateFormat)
	}

	if len(*_moviesInternalID) > consts.ZeroValue {
		options.InternalID = *_moviesInternalID
	}

	if len(*_moviesSort) > consts.ZeroValue {
		options.Sort = *_moviesSort
	}

	if len(*_moviesType) > consts.ZeroValue {
		options.Type = *_moviesType
	}

	return options
}

// convertDateString takes a date string and converts it to date time format,
// if empty return current date
func convertDateString(dateStr string, outputFormat string) string {
	// Parse the input date string using YYYY-MM-DD
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now().Format(consts.DefaultStartDateFormat)
	}

	// Get the current time
	currentTime := time.Now()

	// Combine the parsed date with the current time's hour, minute, second
	finalDateTime := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
		currentTime.Nanosecond(),
		currentTime.Location(),
	)

	// Format the parsed time into the output format
	formattedDate := finalDateTime.Format(outputFormat)
	return formattedDate
}

// GenActionsUsage prints a usage message when an invalid action is provided.
func (c *Command) GenActionsUsage(actions []string) {
	printer.Println("Usage: ./trakt-sync " + c.Name + " -a [action]")
	printer.Println("Available actions:")
	for _, action := range actions {
		printer.Printf("  - %s\n", action)
	}
}

// GenTypeUsage prints a usage message when an invalid type is provided.
func (c *Command) GenTypeUsage(types []string) {
	printer.Println("Usage: ./trakt-sync " + c.Name + " -t [type]")
	printer.Println("Available types:")
	for _, t := range types {
		printer.Printf("  - %s\n", t)
	}
}

// GetHandlerForMap choose handler from map
func (*Command) GetHandlerForMap(action string, allHandlers map[string]handlers.Handler) (handlers.Handler, error) {
	// Lookup and execute handler
	if handler, found := allHandlers[action]; found {
		return handler, nil
	}

	return nil, errors.New("unknown handler")
}
