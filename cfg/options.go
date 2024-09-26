// Package cfg used for process configuration
package cfg

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

// OptionsConfig represents the configuration options for each module
type OptionsConfig struct {
	SearchIDType []string
	SearchType   []string
	SearchField  []string
	Type         []string
	Sort         []string
	Format       []string
	Action       []string
}

// SearchFieldConfig represents the configuration options for search_field depens on type
var SearchFieldConfig = map[string][]string{
	"movie":   {"title", "tagline", "overview", "people", "translations", "aliases"},
	"show":    {"title", "overview", "people", "translations", "aliases"},
	"episode": {"title", "overview"},
	"person":  {"name", "biography"},
	"list":    {"name", "description"},
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
	set := make(map[string]bool)
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

// OptionsFromConfig reads optionf from config file
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

	token, err := readTokenFromFile(fs, config.TokenPath)
	if err != nil {
		return str.Options{}, fmt.Errorf("error reading token:%w", err)
	}

	str.Headers["Authorization"] = "Bearer " + token.AccessToken
	str.Headers["trakt-api-key"] = config.ClientID

	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := ModuleConfig[options.Module]
	if !ok {
		options.Module = "history"
		fmt.Println("Forcing module to history")
	}

	// Check if the provided type is valid for the selected module
	if !IsValidConfigType(moduleConfig.Type, options.Type) {
		return str.Options{}, fmt.Errorf("type '%s' is not valid for module '%s'", options.Type, options.Module)
	}

	if !IsValidConfigType(moduleConfig.Format, options.Format) {
		options.Format = "imdb"
		fmt.Println("Forcing format to imdb")
	}

	if !IsValidConfigType(moduleConfig.Sort, options.Sort) {
		options.Sort = "rank"
		fmt.Println("Forcing sort to rank")
	}

	if len(options.Output) == consts.ZeroValue && options.Module == "lists" {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Type)
	}
	if len(options.Output) == consts.ZeroValue && options.Module == "people" {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Action)
	}

	if len(options.Output) == consts.ZeroValue {
		options.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Type, options.Module)
	}

	if options.Type == "episodes" && options.Format == "imdb" {
		options.Format = "tmdb"
		fmt.Println("Forcing format to tmdb for type episode")
	}

	if len(str.Headers["Authorization"].(string)) == consts.ZeroValue && len(str.Headers["trakt-api-key"].(string)) == consts.ZeroValue {
		return str.Options{}, fmt.Errorf("no valid Authorization header")
	}

	options.Headers = str.Headers
	options.Token = *token

	return *options, nil
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
	mainCount := make(map[string]int)
	for _, element := range allowedElements {
		mainCount[element]++
	}

	// Create a map to count occurrences of each element in userElements.
	subCount := make(map[string]int)
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

	if options.Module == "history" {
		options.Time = "watched_at"
	} else if options.Module == "watchlist" {
		options.Time = "listed_at"
	} else if options.Module == "collection" {
		options.Time = "collected_at"
	} else if options.UserName != "" {
		options.Time = "listed_at"
	}

	return options.Time

}
