// Package cmds used for commands modules
package cmds

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
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
	"a":          true,
	"c":          true,
	"calendars":  true,
	"collection": true,
	"days":       true,
	"ex":         true,
	"f":          true,
	"field":      true,
	"godoc":      true,
	"help":       true,
	"history":    true,
	"i":          true,
	"id_type":    true,
	"lists":      true,
	"o":          true,
	"people":     true,
	"q":          true,
	"search":     true,
	"start_date": true,
	"t":          true,
	"u":          true,
	"v":          true,
	"version":    true,
	"watchlist":  true,
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

// Exec core command function
func (c *Command) Exec(fs afero.Fs, client *internal.Client, config *cfg.Config, args []string) error {

	c.Client = client
	c.Config = config
	c.Flag.Usage = func() {
		HelpFunc(c, c.Name)
	}
	c.registerGlobalFlagsInSet(&c.Flag)
	c.Flag.Parse(args)
	m := c.fetchFlagsMap()
	options, err := cfg.SyncOptionsFromFlags(fs, c.Config, m)

	if err != nil {
		return fmt.Errorf("error command exec: %w", err)
	}

	options.Type = *_strType
	options.Module = c.Name
	switch c.Name {
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

	c.Options = &options

	if !c.ValidFlags() {
		return fmt.Errorf("invalid flags")
	}

	if options.Verbose {
		fmt.Println("Authorization header:" + options.Headers["Authorization"].(string))
		fmt.Println("trakt-api-key header:" + options.Headers["trakt-api-key"].(string))
		fmt.Println("token expiration in seconds:" + strconv.Itoa(options.Token.ExpiritySeconds()))
		fmt.Println("Extended info:" + *_extendedInfo)
		if len(options.Module) > consts.ZeroValue {
			fmt.Println("selected module:" + options.Module)
		}
		fmt.Println(
			str.Format("selected user: {0}, module: {1}, type: {2}, per_page: {3}, format: {4}, action: {5}, sort: {6}",
				options.UserName, options.Module, options.Type, options.PerPage, options.Format, options.Action, options.Sort),
		)

	}
	defer func() error {
		if r := recover(); r != nil {
			if _, ok := r.(fatal); ok {
				return fmt.Errorf("fatal error")
			}
			return fmt.Errorf("panic error:%s", r)
		}
		return nil
	}()

	err = c.Run(c, c.Flag.Args()...)

	if err != nil {
		return fmt.Errorf("error command run: %w", err)
	}

	return nil
}

// BadArgs shows error if command have invalid arguments
func (c *Command) BadArgs(errFormat string, args ...interface{}) {
	fmt.Fprintf(stdout, "error: "+errFormat+"\n\n", args...)
	HelpFunc(c, c.Name)
	panic(fatal{})
}

// Errorf prints out a formatted error with the right prefixes.
func (c *Command) Errorf(errFormat string, args ...interface{}) {
	fmt.Fprintf(stdout, c.Name+": error: "+errFormat+"\n", args...)
	if c.exit == consts.ZeroValue {
		c.exit = consts.DefaultExitCode
	}
}

// Fatalf is like Errorf except the stack unwinds up to the Exec call before
// exiting the application with status code 1.
func (c *Command) Fatalf(errFormat string, args ...interface{}) {
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

func argsToMap(args []string) map[string]bool {
	argMap := make(map[string]bool)
	var key string

	for _, arg := range args {
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
	}

	// If we still have a key at the end, it means it's a single argument without a value
	if key != consts.EmptyString {
		argMap[cleanKey(key)] = true // Set the key to true for bool map
	}

	return argMap
}

// ValidFlags validate if flag is in our list
func (c *Command) ValidFlags() bool {

	for flag := range argsToMap(flag.Args()) {
		if _, ok := Avflags[flag]; !ok {
			return false
		}
	}
	return true
}

func (c *Command) registerGlobalFlagsInSet(fset *flag.FlagSet) {
	flag.VisitAll(func(f *flag.Flag) {
		if fset.Lookup(f.Name) == nil {
			fset.Var(f.Value, f.Name, f.Usage)
		}
	})
}

// Uptime update item time fields
func (c *Command) Uptime(item *str.ExportlistItemJSON, options *str.Options, data *str.ExportlistItem) {

	switch options.Time {

	case "watched_at":
		item.WatchedAt = data.WatchedAt
	case "listed_at":
		item.ListedAt = data.ListedAt
	case "collected_at":
		item.CollectedAt = data.CollectedAt
	case "last_collected_at":
		item.LastCollectedAt = data.LastCollectedAt
	case "updated_at":
		item.UpdatedAt = data.UpdatedAt
	case "last_updated_at":
		item.LastUpdatedAt = data.LastUpdatedAt

	}
}

func (c *Command) IsImdbMovie(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != "episodes" && data.Movie != nil && options.Format == "imdb"
}

func (c *Command) IsImdbShow(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != "episodes" && data.Show != nil && data.Show.IDs.HaveID("Imdb") && 
	options.Format == "imdb"
}

func (c *Command) IsTmdbMovie(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != consts.EpisodesType && data.Movie != nil && 
	data.Movie.IDs.HaveID("Tmdb") && options.Format == "tmdb"
}

func (c *Command) IsTmdbShow(options *str.Options, data *str.ExportlistItem) bool {
	return options.Type != consts.EpisodesType && data.Show != nil &&
		data.Show.IDs.HaveID("Tmdb") && options.Format == "tmdb"
}

func (c *Command) IsTmdbEpisode(options *str.Options, data *str.ExportlistItem) bool {
	return data.Episode.IDs.HaveID(consts.TmdbIDFormat) && options.Format == consts.TmdbFormat
}

func (c *Command) IsTvdbEpisode(options *str.Options, data *str.ExportlistItem) bool {
	return data.Episode.IDs.HaveID("Tvdb") && options.Format == "tvdb"
}

func (c *Command) IsImdbEpisode(options *str.Options, data *str.ExportlistItem) bool {
	return data.Episode.IDs.HaveID(consts.ImdbIDFormat) && options.Format == consts.ImdbFormat
}

// ExportListProcess process list items
func (c *Command) ExportListProcess(
	data *str.ExportlistItem, options *str.Options,
	findDuplicates []any, exportJSON []str.ExportlistItemJSON,
) ([]any, []str.ExportlistItemJSON) {

	//If movie or show export by format imdb
	if c.IsImdbMovie(options, data) {

		//fmt.Println("movie or show by format imdb")
		if !data.Movie.IDs.HaveID("Imdb") {
			noImdb := "no-imdb"
			data.Movie.IDs.Imdb = &noImdb
		}

		findDuplicates = append(findDuplicates, *data.Movie.IDs.Imdb)
		emap := str.ExportlistItemJSON{
			Imdb:  data.Movie.IDs.Imdb,
			Trakt: data.Movie.IDs.Trakt,
			Title: data.Movie.Title}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt
		emap.Year = data.Movie.Year
		emap.Metadata = data.Metadata
		exportJSON = append(exportJSON, emap)

	} else if c.IsImdbShow(options, data) {

		findDuplicates = append(findDuplicates, *data.Show.IDs.Imdb)
		emap := str.ExportlistItemJSON{
			Imdb:  data.Show.IDs.Imdb,
			Trakt: data.Show.IDs.Trakt,
			Title: data.Show.Title}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt

		exportJSON = append(exportJSON, emap)

	} else if c.IsTmdbMovie(options, data) {

		findDuplicates = append(findDuplicates, *data.Movie.IDs.Tmdb)
		emap := str.ExportlistItemJSON{
			Tmdb:  data.Movie.IDs.Tmdb,
			Trakt: data.Movie.IDs.Trakt,
			Title: data.Movie.Title}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt
		exportJSON = append(exportJSON, emap)

	} else if c.IsTmdbShow(options, data) {

		findDuplicates = append(findDuplicates, *data.Show.IDs.Tmdb)
		emap := str.ExportlistItemJSON{
			Tmdb:  data.Show.IDs.Tmdb,
			Trakt: data.Show.IDs.Trakt,
			Title: data.Show.Title}
		c.Uptime(&emap, options, data)
		exportJSON = append(exportJSON, emap)

	} else if c.IsTvdbEpisode(options, data) {

		//fmt.Println("episode export by format tvdb")
		findDuplicates = append(findDuplicates, *data.Episode.IDs.Tvdb)

		if len(*data.Episode.Title) == consts.ZeroValue {
			notitle := consts.NoEpisodeTitle
			data.Episode.Title = &notitle
		}

		if len(*data.Show.Title) == consts.ZeroValue {
			notitle := consts.NoShowTitle
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJSON{
			Tvdb:  data.Episode.IDs.Tvdb,
			Trakt: data.Episode.IDs.Trakt}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
		emap.Show = &str.Show{Title: data.Show.Title}

		exportJSON = append(exportJSON, emap)

	} else if c.IsImdbEpisode(options, data) {

		//fmt.Println("episode export by format imdb")
		findDuplicates = append(findDuplicates, *data.Episode.IDs.Imdb)

		if len(*data.Episode.Title) == consts.ZeroValue {
			notitle := consts.NoEpisodeTitle
			data.Episode.Title = &notitle
		}

		if len(*data.Show.Title) == consts.ZeroValue {
			notitle := consts.NoShowTitle
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJSON{
			Imdb:  data.Episode.IDs.Imdb,
			Trakt: data.Episode.IDs.Trakt}
		c.Uptime(&emap, options, data)

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
		emap.Show = &str.Show{Title: data.Show.Title}

		exportJSON = append(exportJSON, emap)

	} else if c.IsTmdbEpisode(options, data) {

		//fmt.Println("episode export by format tmdb")
		findDuplicates = append(findDuplicates, *data.Episode.IDs.Tmdb)

		if len(*data.Episode.Title) == consts.ZeroValue {
			notitle := consts.NoEpisodeTitle
			data.Episode.Title = &notitle
		}

		if len(*data.Show.Title) == consts.ZeroValue {
			notitle := consts.NoShowTitle
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJSON{
			Tmdb:  data.Episode.IDs.Tmdb,
			Trakt: data.Episode.IDs.Trakt}
		c.Uptime(&emap, options, data)

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
		emap.Show = &str.Show{Title: data.Show.Title}

		exportJSON = append(exportJSON, emap)
	}

	return findDuplicates, exportJSON

}

// PrepareQueryString for remove or replace unwanted signs from query string
func (c *Command) PrepareQueryString(q string) *string {
	return &q
}

// UpdateOptionsWithCommandFlags update options depends on command flags
func (c *Command) UpdateOptionsWithCommandFlags(options *str.Options) *str.Options {

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
		switch options.Module {
		case "calendars":
			switch options.Action {
			case "my-shows", "all-shows":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"shows",
					strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-new-shows", "all-new-shows":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"new_shows",
					strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-season-premieres", "all-season-premieres":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"season_premieres",
					strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-finales", "all-finales":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"finales",
					strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-movies", "all-movies":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"movies",
					strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-dvd", "all-dvd":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"dvd",
					strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))

			default:
				options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
			}
		case "search":
			switch options.Action {
			case "text-query":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"query",
					strings.ReplaceAll(options.Type, ",", consts.EmptyString))
			case "id-lookup":
				options.Output = fmt.Sprintf(
					consts.DefaultOutputFormat3,
					options.Module,
					"lookup",
					strings.ReplaceAll(options.SearchIDType, ",", consts.EmptyString))
			default:
				options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
			}

		default:
			options.Output = fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Type, options.Format)
		}

	}

	return options

}
