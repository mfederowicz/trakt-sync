// Package cfg used for process configuration
package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

// OptionsConfig represents the configuration options for each module
type OptionsConfig struct {
	SearchIDType []string
	CommentType  []string
	SearchType   []string
	SearchField  []string
	Type         []string
	Period       []string
	Sort         []string
	Format       []string
	Action       []string
	Privacy      []string
}

// SearchFieldConfig represents the configuration options for search_field depens on type
var SearchFieldConfig = map[string][]string{
	"movie":   {"title", "tagline", "overview", "people", "translations", "aliases"},
	"show":    {"title", "overview", "people", "translations", "aliases"},
	"episode": {"title", "overview"},
	"person":  {"name", "biography"},
	"list":    {"name", "description"},
}

// ModuleActionConfig represents the configuration options depens on module:action key
var ModuleActionConfig = map[string]OptionsConfig{

	"comments:trending": {
		Type: []string{"all", "movies", "shows", "seasons", "episodes", "lists"},
		Sort: []string{"newest", "oldest", "likes", "replies", "highest", "lowest", "plays"},
	},
	"comments:recent": {
		Type: []string{"all", "movies", "shows", "seasons", "episodes", "lists"},
		Sort: []string{"newest", "oldest", "likes", "replies", "highest", "lowest", "plays"},
	},
	"movies:comments": {
		Type: []string{},
		Sort: []string{"newest", "oldest", "likes", "replies", "highest", "lowest", "plays"},
	},
	"movies:lists": {
		Type: []string{"all", "personal", "official", "watchlists", "favorites"},
		Sort: []string{"popular", "likes", "comments", "items", "added", "updated"},
	},
	"users:watched": {
		Type: []string{"movies", "shows"},
		Sort: []string{},
	},
	"notes:notes": {
		Privacy: []string{"private", "friends", "public"},
	},
	"notes:note": {
		Privacy: []string{"private", "friends", "public"},
	},
	"notes:item": {
		Privacy: []string{"private", "friends", "public"},
	},
}

// ModuleConfig represents the configuration options for all modules
var ModuleConfig = map[string]OptionsConfig{
	"watchlist": {
		SearchIDType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"collection": {
		SearchIDType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"comments": {
		SearchIDType: []string{},
		SearchType:   []string{},
		CommentType:  []string{"all", "review", "shouts"},
		SearchField:  []string{},
		Type:         []string{"all", "movies", "shows", "seasons", "episodes", "lists"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},

	"history": {
		SearchIDType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"lists": {
		SearchIDType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"movies": {
		SearchIDType: []string{},
		SearchType:   []string{},
		CommentType:  []string{"all", "review", "shouts"},
		SearchField:  []string{},
		Type:         []string{"all", "movies", "shows", "seasons", "episodes", "lists"},
		Period:       []string{"all", "daily", "weekly", "monthly"},
		Sort:         []string{"newest", "oldest", "likes", "replies", "highest", "lowest", "plays"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},

	"people": {
		SearchIDType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons", "all", "personal", "official"},
		Sort:         []string{"rank", "added", "released", "title", "popular", "likes", "comments", "items", "added", "updated"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"search": {
		SearchIDType: []string{"trakt", "imdb", "tmdb", "tvdb"},
		SearchType:   []string{"movie", "show", "episode", "person", "list", "podcast", "podcast_episode"},
		SearchField:  []string{"title", "aliases", "biography", "description", "episode", "name", "overview", "people", "show", "tagline", "translations"},
		Type:         []string{"movies", "shows", "episodes", "persons", "all", "personal", "official"},
		Sort:         []string{"rank", "added", "released", "title", "popular", "likes", "comments", "items", "added", "updated"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},

	"shows": {
		SearchIDType: []string{},
		SearchType:   []string{},
		CommentType:  []string{"all", "review", "shouts"},
		SearchField:  []string{},
		Type:         []string{"all", "movies", "shows", "seasons", "episodes", "lists"},
		Period:       []string{"all", "daily", "weekly", "monthly"},
		Sort:         []string{"newest", "oldest", "likes", "replies", "highest", "lowest", "plays"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},

	"users": {
		SearchIDType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows"},
		Sort:         []string{},
		Format:       []string{},
		Action:       []string{},
	},
	"notes": {
		Privacy: []string{"private", "friends", "public"},
	},
}

// ValidateConfig validates if the provided configuration is allowed for the given module
func ValidateConfig(module string, config OptionsConfig) bool {
	allowedConfig := ModuleConfig[module]
	return isSubset(config.Type, allowedConfig.Type) &&
		isSubset(config.Sort, allowedConfig.Sort) &&
		isSubset(config.Format, allowedConfig.Format)
}

// isSubset checks if slice a is a subset of slice b
func isSubset(a, b []string) bool {
	set := map[string]bool{}
	for _, item := range b {
		set[item] = true
	}
	for _, item := range a {
		if !set[item] {
			return false
		}
	}
	return true
}

// SyncOptionsFromFlags reads options from user flags
func SyncOptionsFromFlags(fs afero.Fs, config *Config, flagMap map[string]string) (str.Options, error) {
	cfg, err := MergeConfigs(DefaultConfig(), config, flagMap)

	if err != nil {
		return str.Options{}, fmt.Errorf("error sync options from flags: %w", err)
	}

	return OptionsFromConfig(fs, cfg)
}

// OptionsFromConfig reads options from config file
func OptionsFromConfig(fs afero.Fs, config *Config) (str.Options, error) {
	options := &str.Options{}
	options.Verbose = config.Verbose
	options.Type = config.Type
	options.SearchType = config.SearchType
	options.SearchField = config.SearchField
	options.Module = config.Module
	options.Output = config.Output
	options.Format = config.Format
	options.List = config.List
	options.UserName = config.UserName
	options.ID = config.ID
	options.PerPage = config.PerPage
	options.Sort = config.Sort
	options.Action = config.Action
	options.PagesLimit = config.PagesLimit

	token, err := readTokenFromFile(fs, config.TokenPath)
	if err != nil {
		return str.Options{}, fmt.Errorf("error reading token:%w", err)
	}

	str.Headers["Authorization"] = "Bearer " + token.AccessToken
	str.Headers["trakt-api-key"] = config.ClientID

	if len(str.Headers["Authorization"].(string)) == consts.ZeroValue && len(str.Headers["trakt-api-key"].(string)) == consts.ZeroValue {
		return str.Options{}, errors.New("no valid Authorization header")
	}

	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := ModuleConfig[options.Module]
	if !ok {
		options.Module = "history"
		printer.Println("Forcing module to history")
	}

	// Check if the provided type is valid for the selected module
	if !IsValidConfigType(moduleConfig.Type, options.Type) {
		return str.Options{}, fmt.Errorf("type '%s' is not valid for module '%s'", options.Type, options.Module)
	}
	options = optionsFromModuleConfig(moduleConfig, options)
	options.Headers = str.Headers
	options.Token = *token
	options.Output = optionsFromConfigOutput(options)

	return *options, nil
}

func optionsFromModuleConfig(moduleConfig OptionsConfig, options *str.Options) *str.Options {
	if !IsValidConfigType(moduleConfig.Format, options.Format) {
		options.Format = "imdb"
		printer.Println("Forcing format to imdb")
	}

	if !IsValidConfigType(moduleConfig.Sort, options.Sort) {
		options.Sort = "rank"
		printer.Println("Forcing sort to rank")
	}

	if options.Type == "episodes" && options.Format == "imdb" {
		options.Format = "tmdb"
		printer.Println("Forcing format to tmdb for type episode")
	}
	return options
}

func optionsFromConfigOutput(options *str.Options) string {
	if len(options.Output) == consts.ZeroValue && options.Module == "lists" {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Type)
	}

	if len(options.Output) == consts.ZeroValue && options.Module == "people" {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Action)
	}

	if len(options.Output) == consts.ZeroValue {
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Type, options.Module)
	}

	return options.Output
}

// IsValidConfigType checks if the provided type is valid for the module
func IsValidConfigType(allowedTypes []string, userType string) bool {
	for _, t := range allowedTypes {
		if t == userType {
			return true
		}
	}
	return false
}

// IsValidConfigTypeSlice checks if all elements of userElements are in allowedElements,
// considering the counts of each element.
func IsValidConfigTypeSlice(allowedElements []string, userElements str.Slice) bool {
	if len(userElements) == consts.ZeroValue {
		return true
	}

	// Create a map to count occurrences of each element in allowedElements.
	mainCount := map[string]int{}
	for _, element := range allowedElements {
		mainCount[element]++
	}

	// Create a map to count occurrences of each element in userElements.
	subCount := map[string]int{}
	for _, element := range userElements {
		subCount[element]++
	}

	// Check if allowedElements contains all elements of userElements with the required counts.
	for element, count := range subCount {
		if mainCount[element] < count {
			return false
		}
	}

	return true
}

// ReadTokenFromFile reads the token from the specified file
func readTokenFromFile(fs afero.Fs, filePath string) (*str.Token, error) {
	data, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return nil, err
	}

	var token str.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// GetOptionTime config Time depends on Module name
func GetOptionTime(options *str.Options) string {
	switch options.Module {
	case "history":
		options.Time = "watched_at"
	case "watchlist", "collection":
		options.Time = "listed_at"
	default:
		if options.UserName != "" {
			options.Time = "listed_at"
		}
	}
	return options.Time
}

// GetOutputForModule generates output value depends on module name
func GetOutputForModule(options *str.Options) string {
	allOutputs := map[string]string{
		consts.Calendars:       getOutputForModuleCalendars(options),
		consts.Certifications:  getOutputForModuleCertifications(options),
		consts.Comments:        getOutputForModuleComments(options),
		consts.Countries:       getOutputForModuleCountries(options),
		consts.Genres:          getOutputForModuleGenres(options),
		consts.Languages:       getOutputForModuleLanguages(options),
		consts.Search:          getOutputForModuleSearch(options),
		consts.Shows:           getOutputForModuleShows(options),
		consts.Users:           getOutputForModuleUsers(options),
		consts.Lists:           getOutputForModuleLists(options),
		consts.Movies:          getOutputForModuleMovies(options),
		consts.Networks:        getOutputForModuleNetworks(options),
		consts.Notes:           getOutputForModuleNotes(options),
		consts.Recommendations: getOutputForModuleRecommendations(options),
	}

	if output, found := allOutputs[options.Module]; found {
		return output
	}
	return fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Type, options.Format)
}

func getOutputForModuleRecommendations(options *str.Options) string {
	switch options.Action {
	case consts.Movies, consts.Shows:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Action)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Action)
	}
	return options.Output
}

func getOutputForModuleNotes(options *str.Options) string {
	return fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Action, options.InternalID)
}

func getOutputForModuleNetworks(options *str.Options) string {
	switch options.Action {
	case consts.Lists:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Action)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Type)
	}

	return options.Output
}

func getOutputForModuleMovies(options *str.Options) string {
	switch options.Action {
	case consts.Trending, consts.Popular, consts.Anticipated, consts.Boxoffice, consts.Updates, consts.UpdatedIDs:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Action)
	case consts.Favorited, consts.Played, consts.Watched, consts.Collected:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Action, options.Period)
	case consts.Summary, consts.Aliases, consts.Releases, consts.Translations, consts.Comments, consts.Lists, consts.People, consts.Ratings, consts.Related, consts.Stats, consts.Studios, consts.Watching, consts.Videos:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Action, options.InternalID)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Type)
	}

	return options.Output
}

func getOutputForModuleShows(options *str.Options) string {
	switch options.Action {
	case consts.Trending, consts.Popular, consts.Anticipated, consts.Boxoffice, consts.Updates, consts.UpdatedIDs:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Action)
	case consts.Favorited, consts.Played, consts.Watched, consts.Collected:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Action, options.Period)
	case consts.Summary, consts.Aliases, consts.Releases, consts.Translations, consts.Comments, consts.Lists, consts.People, consts.Ratings, consts.Related, consts.Stats, consts.Studios, consts.Watching, consts.Videos:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat3, options.Module, options.Action, options.InternalID)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Type)
	}

	return options.Output
}

func getOutputForModuleLanguages(options *str.Options) string {
	switch options.Type {
	case consts.Movies, consts.Shows:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Type)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
	}

	return options.Output
}

func getOutputForModuleGenres(options *str.Options) string {
	switch options.Type {
	case consts.Movies, consts.Shows:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Type)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
	}

	return options.Output
}

func getOutputForModuleCountries(options *str.Options) string {
	switch options.Type {
	case consts.Movies, consts.Shows:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Type)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
	}

	return options.Output
}

func getOutputForModuleLists(options *str.Options) string {
	switch options.Action {
	case consts.Trending:
	case consts.Popular:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Action)
	case consts.List:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringString, "trakt_", options.InternalID),
		)

	case consts.Likes:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringString, "likes_trakt_", options.InternalID),
		)
	case consts.Items:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringString, "items_trakt_", options.InternalID),
		)
	case consts.Comments:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringString, "comments_trakt_", options.InternalID),
		)

	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Type)
	}

	return options.Output
}

func getOutputForModuleUsers(options *str.Options) string {
	switch options.Action {
	case consts.Watched:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat3,
			options.Module,
			options.Action,
			strings.ReplaceAll(options.Type, consts.CommaString, consts.EmptyString))
	case consts.Stats:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Action)
	case consts.Lists:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat3,
			options.Module,
			options.Action,
			strings.ReplaceAll(options.Type, consts.CommaString, consts.EmptyString))
	case consts.SavedFilters:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat3,
			options.Module,
			options.Action,
			strings.ReplaceAll(options.Type, consts.CommaString, consts.EmptyString))
	case consts.Settings:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Action)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Type)
	}

	return options.Output
}

func getOutputForModuleSearch(options *str.Options) string {
	switch options.Action {
	case consts.TextQuery:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat3,
			options.Module,
			consts.Query,
			strings.ReplaceAll(options.Type, consts.CommaString, consts.EmptyString))
	case consts.IDLookup:
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat3,
			options.Module,
			consts.Lookup,
			strings.ReplaceAll(options.SearchIDType, consts.CommaString, consts.EmptyString))
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
	}

	return options.Output
}

func getOutputForModuleCalendars(options *str.Options) string {
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
	return options.Output
}

func getOutputForModuleCertifications(options *str.Options) string {
	switch options.Type {
	case "movies", "shows":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			options.Type)
	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat1, options.Module)
	}

	return options.Output
}

func getOutputForModuleComments(options *str.Options) string {
	switch options.Action {
	case "comment":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringDigit, "comment_", options.CommentID),
		)
	case "replies":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringDigit, "replies_", options.CommentID),
		)
	case "item":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringDigit, "item_", options.CommentID),
		)
	case "likes":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			fmt.Sprintf(consts.StringDigit, "likes_", options.CommentID),
		)
	case "trending":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			consts.Trending,
		)
	case "recent":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			consts.Recent,
		)
	case "updates":
		options.Output = fmt.Sprintf(
			consts.DefaultOutputFormat2,
			options.Module,
			consts.Updates,
		)

	default:
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, options.Type)
	}

	return options.Output
}
