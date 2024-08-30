package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

// OptionsConfig represents the configuration options for each module
type OptionsConfig struct {
	SearchIdType []string
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
		SearchIdType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"collection": {
		SearchIdType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"history": {
		SearchIdType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"lists": {
		SearchIdType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons"},
		Sort:         []string{"rank", "added", "released", "title"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"people": {
		SearchIdType: []string{},
		SearchType:   []string{},
		SearchField:  []string{},
		Type:         []string{"movies", "shows", "episodes", "persons", "all", "personal", "official"},
		Sort:         []string{"rank", "added", "released", "title", "popular", "likes", "comments", "items", "added", "updated"},
		Format:       []string{"imdb", "tmdb", "tvdb", "tvrage", "trakt"},
		Action:       []string{},
	},
	"search": {
		SearchIdType: []string{"trakt", "imdb", "tmdb", "tvdb"},
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
func SyncOptionsFromFlags(fs afero.Fs, config *Config, flagMap map[string]string) str.Options {
	cfg := MergeConfigs(DefaultConfig(), config, flagMap)
	return OptionsFromConfig(fs, cfg)
}

func OptionsFromConfig(fs afero.Fs, config *Config) str.Options {

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
	options.Id = config.Id
	options.PerPage = config.PerPage
	options.Sort = config.Sort
	options.Action = config.Action

	token, err := readTokenFromFile(fs, config.TokenPath)
	if err != nil {
		fmt.Println("Error reading token:", err)
		os.Exit(1)
	}

	str.Headers["Authorization"] = "Bearer " + token.AccessToken
	str.Headers["trakt-api-key"] = config.ClientId

	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := ModuleConfig[options.Module]
	if !ok {
		options.Module = "history"
		fmt.Println("Forcing module to history")
	}

	// Check if the provided type is valid for the selected module
	if !IsValidConfigType(moduleConfig.Type, options.Type) {
		fmt.Printf("Type '%s' is not valid for module '%s'\n", options.Type, options.Module)
		os.Exit(1)
	}

	if !IsValidConfigType(moduleConfig.Format, options.Format) {
		options.Format = "imdb"
		fmt.Println("Forcing format to imdb")
	}

	if !IsValidConfigType(moduleConfig.Sort, options.Sort) {
		options.Sort = "rank"
		fmt.Println("Forcing sort to rank")
	}

	if len(options.Output) == 0 && options.Module == "lists" {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Type)
	}
	if len(options.Output) == 0 && options.Module == "people" {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Action)
	}

	if len(options.Output) == 0 {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Type, options.Module)
	}

	if options.Type == "episodes" && options.Format == "imdb" {
		options.Format = "tmdb"
		fmt.Println("Forcing format to tmdb for type episode")
	}

	if len(str.Headers["Authorization"].(string)) == 0 && len(str.Headers["trakt-api-key"].(string)) == 0 {
		fmt.Println("No valid Authorization header")
		os.Exit(1)

	}

	options.Headers = str.Headers
	options.Token = *token

	return *options
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
func IsValidConfigTypeSlice(allowedElements []string, userElements str.StrSlice) bool {
	if len(userElements) == 0 {
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
