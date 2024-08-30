package cmds

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

var (
	_userName     = flag.String("u", cfg.DefaultConfig().UserName, UserlistUsage)
	_strType      = flag.String("t", cfg.DefaultConfig().Type, TypeUsage)
	_output       = flag.String("o", cfg.DefaultConfig().Output, OutputUsage)
	_format       = flag.String("f", cfg.DefaultConfig().Format, FormatUsage)
	_extendedInfo = flag.String("ex", "", ExtendedInfoUsage)
	_query        = flag.String("query", "", "")
	_years        = flag.String("years", "", "")
	_genres       = flag.String("genres", "", "")
	_languages    = flag.String("languages", "", "")
	_countries    = flag.String("countries", "", "")
	_runtimes     = flag.String("runtimes", "", "")
	_studio_ids   = flag.String("studio_ids", "", "")
)

const (
	ConfigUsage       = "allow to overwrite default config filename"
	VersionUsage      = "get trakt-sync version"
	VerboseUsage      = "print additional verbose information"
	OutputUsage       = "allow to overwrite default output filename"
	FormatUsage       = "allow to overwrite default ID type format"
	TypeUsage         = "allow to overwrite type"
	ListUsage         = "allow to overwrite list"
	UserlistUsage     = "allow to export a user custom list"
	StartDateUsage    = "allow to overwrite start_date"
	DaysUsage         = "allow to overwrite days"
	ListIdUsage       = "allow to export a specific custom list"
	ModuleUsage       = "allow use selected module"
	ActionUsage       = "allow use selected action"
	QueryUsage        = "allow use selected query"
	FieldUsage        = "allow use selected field"
	SortUsage         = "allow to overwrite sort"
	ExtendedInfoUsage = "allow to overwrite extended flag"
)

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
	Run     func(cmd *Command, args ...string)
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

func (c *Command) Exec(fs afero.Fs, client *internal.Client, config *cfg.Config, args []string) {

	c.Client = client
	c.Config = config
	c.Flag.Usage = func() {
		HelpFunc(c, c.Name)
	}
	c.registerGlobalFlagsInSet(&c.Flag)
	c.Flag.Parse(args)
	m := c.fetchFlagsMap()
	options := cfg.SyncOptionsFromFlags(fs, c.Config, m)

	options.Type = *_strType
	options.Module = c.Name
	switch c.Name {
	case "people":
		options.Action = *_action
		options.Id = *_personId
		options.Type = *_action
	case "calendars":
		options.Action = *_cal_action
		options.StartDate = *_cal_startDate
		options.Days = *_cal_days
	case "search":
		options.Action = *_search_action
		options.SearchType = _search_type
		options.SearchField = _search_field
		options.Id = *_search_id
		options.SearchIdType = *_search_id_type
	case "watchlist":
	case "collection":
	case "history":
		options.Format = *_format
	}

	c.Options = &options

	if !c.ValidFlags() {
		fmt.Println("invalid flags")
		os.Exit(1)
	}

	if options.Verbose {
		fmt.Println("Authorization header:" + options.Headers["Authorization"].(string))
		fmt.Println("trakt-api-key header:" + options.Headers["trakt-api-key"].(string))
		fmt.Println("token expiration in seconds:" + strconv.Itoa(options.Token.ExpiritySeconds()))
		fmt.Println("Extended info:" + *_extendedInfo)
		if len(options.Module) > 0 {
			fmt.Println("selected module:" + options.Module)
		}
		fmt.Println(
			str.Format("selected user: {0}, module: {1}, type: {2}, per_page: {3}, format: {4}, action: {5}, sort: {6}",
				options.UserName, options.Module, options.Type, options.PerPage, options.Format, options.Action, options.Sort),
		)

	}
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(fatal); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	c.Run(c, c.Flag.Args()...)

	if c.exit != 0 {
		os.Exit(c.exit)
	}
}

func (c *Command) BadArgs(errFormat string, args ...interface{}) {
	fmt.Fprintf(stdout, "error: "+errFormat+"\n\n", args...)
	HelpFunc(c, c.Name)
	panic(fatal{})
}

// Errorf prints out a formatted error with the right prefixes.
func (c *Command) Errorf(errFormat string, args ...interface{}) {
	fmt.Fprintf(stdout, c.Name+": error: "+errFormat+"\n", args...)
	if c.exit == 0 {
		c.exit = 1
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
		if arg[0] == '-' {
			// If we already have a key, it means it's a single argument without a value
			if key != "" {
				argMap[cleanKey(key)] = true // Set the key to true for bool map
			}
			key = arg
		} else {
			// If we have a key, assign the value to it
			if key != "" {
				argMap[cleanKey(key)] = true // Set the key to true for bool map
				key = ""
			} else {
				// If we don't have a key, consider it a standalone argument
				argMap[arg] = true // Set the key to true for bool map
			}
		}
	}

	// If we still have a key at the end, it means it's a single argument without a value
	if key != "" {
		argMap[cleanKey(key)] = true // Set the key to true for bool map
	}

	return argMap
}

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

func (c *Command) Uptime(item *str.ExportlistItemJson, options *str.Options, data *str.ExportlistItem) {

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

func (c *Command) ExportListProcess(
	data *str.ExportlistItem, options *str.Options,
	find_duplicates []any, export_json []str.ExportlistItemJson,
) ([]any, []str.ExportlistItemJson) {

	//If movie or show export by format imdb
	if options.Type != "episodes" && data.Movie != nil && options.Format == "imdb" {

		//fmt.Println("movie or show by format imdb")
		if !data.Movie.Ids.HaveId("Imdb") {
			noImdb := "no-imdb"
			data.Movie.Ids.Imdb = &noImdb
		}

		find_duplicates = append(find_duplicates, *data.Movie.Ids.Imdb)
		emap := str.ExportlistItemJson{
			Imdb:  data.Movie.Ids.Imdb,
			Trakt: data.Movie.Ids.Trakt,
			Title: data.Movie.Title}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt
		emap.Year = data.Movie.Year
		emap.Metadata = data.Metadata
		export_json = append(export_json, emap)

	} else if options.Type != "episodes" && data.Show != nil && data.Show.Ids.HaveId("Imdb") && options.Format == "imdb" {

		find_duplicates = append(find_duplicates, *data.Show.Ids.Imdb)
		emap := str.ExportlistItemJson{
			Imdb:  data.Show.Ids.Imdb,
			Trakt: data.Show.Ids.Trakt,
			Title: data.Show.Title}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt

		export_json = append(export_json, emap)

	} else if options.Type != "episodes" && data.Movie != nil && data.Movie.Ids.HaveId("Tmdb") && options.Format == "tmdb" {

		find_duplicates = append(find_duplicates, *data.Movie.Ids.Tmdb)
		emap := str.ExportlistItemJson{
			Tmdb:  data.Movie.Ids.Tmdb,
			Trakt: data.Movie.Ids.Trakt,
			Title: data.Movie.Title}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt
		export_json = append(export_json, emap)

	} else if options.Type != "episodes" && data.Show != nil && data.Show.Ids.HaveId("Tmdb") && options.Format == "tmdb" {

		find_duplicates = append(find_duplicates, *data.Show.Ids.Tmdb)
		emap := str.ExportlistItemJson{
			Tmdb:  data.Show.Ids.Tmdb,
			Trakt: data.Show.Ids.Trakt,
			Title: data.Show.Title}
		c.Uptime(&emap, options, data)
		export_json = append(export_json, emap)

	} else if data.Episode.Ids.HaveId("Tmdb") && options.Format == "tmdb" {
		//fmt.Println("episode export by format tmdb")
		find_duplicates = append(find_duplicates, *data.Episode.Ids.Tmdb)

		if len(*data.Episode.Title) == 0 {
			notitle := "no episode title"
			data.Episode.Title = &notitle
		}

		if data.Show != nil && len(*data.Show.Title) == 0 {
			notitle := "no show title"
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJson{
			Tmdb:  data.Episode.Ids.Tmdb,
			Trakt: data.Episode.Ids.Trakt}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}

		if data.Show != nil {
			emap.Show = data.Show
		} else {
			notitle := "no show title"
			emap.Show = &str.Show{Title: &notitle}
		}

		export_json = append(export_json, emap)

	} else if data.Episode.Ids.HaveId("Tvdb") && options.Format == "tvdb" {

		//fmt.Println("episode export by format tvdb")
		find_duplicates = append(find_duplicates, *data.Episode.Ids.Tvdb)

		if len(*data.Episode.Title) == 0 {
			notitle := "no episode title"
			data.Episode.Title = &notitle
		}

		if len(*data.Show.Title) == 0 {
			notitle := "no show title"
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJson{
			Tvdb:  data.Episode.Ids.Tvdb,
			Trakt: data.Episode.Ids.Trakt}
		c.Uptime(&emap, options, data)

		emap.UpdatedAt = data.UpdatedAt

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
		emap.Show = &str.Show{Title: data.Show.Title}

		export_json = append(export_json, emap)

	} else if data.Episode.Ids.HaveId("Imdb") && options.Format == "imdb" {

		//fmt.Println("episode export by format imdb")
		find_duplicates = append(find_duplicates, *data.Episode.Ids.Imdb)

		if len(*data.Episode.Title) == 0 {
			notitle := "no episode title"
			data.Episode.Title = &notitle
		}

		if len(*data.Show.Title) == 0 {
			notitle := "no show title"
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJson{
			Imdb:  data.Episode.Ids.Imdb,
			Trakt: data.Episode.Ids.Trakt}
		c.Uptime(&emap, options, data)

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
		emap.Show = &str.Show{Title: data.Show.Title}

		export_json = append(export_json, emap)

	} else if data.Episode.Ids.HaveId("Tmdb") && options.Format == "tmdb" {

		//fmt.Println("episode export by format tmdb")
		find_duplicates = append(find_duplicates, *data.Episode.Ids.Tmdb)

		if len(*data.Episode.Title) == 0 {
			notitle := "no episode title"
			data.Episode.Title = &notitle
		}

		if len(*data.Show.Title) == 0 {
			notitle := "no show title"
			data.Show.Title = &notitle
		}

		emap := str.ExportlistItemJson{
			Tmdb:  data.Episode.Ids.Tmdb,
			Trakt: data.Episode.Ids.Trakt}
		c.Uptime(&emap, options, data)

		emap.Season = &str.Season{Number: data.Episode.Season}
		emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
		emap.Show = &str.Show{Title: data.Show.Title}

		export_json = append(export_json, emap)
	}

	return find_duplicates, export_json

}
func (c *Command) PrepareQueryString(q string) *string {
	return &q
}

func (c *Command) UpdateOptionsWithCommandFlags(options *str.Options) *str.Options {

	// Check if the provided module exists in ModuleConfig
	//moduleConfig, ok := ModuleConfig[options.Module]
	// if !ok {
	// 	options.Module = "history"
	// 	fmt.Println("Forcing module to history")
	// }

	if len(*_search_query) > 0 {
		options.Query = *c.PrepareQueryString(*_search_query)
	}

	if len(*_extendedInfo) > 0 {
		options.ExtendedInfo = *_extendedInfo
	}

	if len(*_format) > 0 {
		options.Format = *_format
	}

	if len(*_output) > 0 {
		options.Output = *_output
	} else {
		switch options.Module {
		case "calendars":
			switch options.Action {
			case "my-shows", "all-shows":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "shows", strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-new-shows", "all-new-shows":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "new_shows", strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-season-premieres", "all-season-premieres":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "season_premieres", strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-finales", "all-finales":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "finales", strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-movies", "all-movies":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "movies", strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))
			case "my-dvd", "all-dvd":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "dvd", strings.ReplaceAll(options.StartDate, "-", "")+"_"+strconv.Itoa(options.Days))

			default:
				options.Output = fmt.Sprintf("export_%s.json", options.Module)
			}
		case "search":
			switch options.Action {
			case "text-query":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "query", strings.ReplaceAll(options.Type, ",", ""))
			case "id-lookup":
				options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, "lookup", strings.ReplaceAll(options.SearchIdType, ",", ""))
			default:
				options.Output = fmt.Sprintf("export_%s.json", options.Module)
			}

		default:
			options.Output = fmt.Sprintf("export_%s_%s_%s.json", options.Module, options.Type, options.Format)
		}

	}

	return options

}
