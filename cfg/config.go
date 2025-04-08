// Package cfg used for process configuration
package cfg

import (
	"errors"
	"flag"
	"fmt"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

// Config struct for app.
type Config struct {
	Action            string    `toml:"action"`
	ClientID          string    `toml:"client_id"`
	ClientSecret      string    `toml:"client_secret"`
	Comment           string    `toml:"comment"`
	CommentID         int       `toml:"comment_id"`
	CommentType       string    `toml:"comment_type"`
	CommentsSort      string    `toml:"sort"`
	ConfigPath        string    `toml:"config_path"`
	Days              int       `toml:"days"`
	Delete            bool      `toml:"delete"`
	Episode           int       `toml:"episode"`
	EpisodeAbs        int       `toml:"episode_abs"`
	EpisodeCode       string    `toml:"episode_code"`
	ErrorCode         int       `toml:"errorCode"`
	Field             string    `toml:"field"`
	Format            string    `toml:"format"`
	Hide              bool      `toml:"hide"`
	ID                string    `toml:"id"`
	IgnoreCollected   string    `toml:"ignore_collected"`
	IgnoreWatchlisted string    `toml:"ignore_watchlisted"`
	IncludeReplies    string    `toml:"include_replies"`
	InternalID        string    `toml:"trakt_id"`
	Item              string    `toml:"item"`
	List              string    `toml:"list"`
	Module            string    `toml:"module"`
	MoviesCountry     string    `toml:"country"`
	MoviesLanguage    string    `toml:"language"`
	MoviesPeriod      string    `toml:"period"`
	MoviesSort        string    `toml:"sort"`
	MoviesType        string    `toml:"type"`
	Msg               string    `toml:"msg"`
	Notes             string    `toml:"notes"`
	NotesID           int       `toml:"notes_id"`
	Output            string    `toml:"output"`
	PagesLimit        int       `toml:"pages_limit"`
	PerPage           int       `toml:"per_page"`
	Privacy           string    `toml:"privacy"`
	Progress          float64   `toml:"progress"`
	Query             string    `toml:"query"`
	RedirectURI       string    `toml:"redirect_uri"`
	Remove            bool      `toml:"remove"`
	Reply             string    `toml:"reply"`
	SearchField       str.Slice `toml:"search_field"`
	SearchIDType      string    `toml:"search_id_type"`
	SearchType        str.Slice `toml:"search_type"`
	Season            int       `toml:"season"`
	Sort              string    `toml:"sort"`
	Spoiler           bool      `toml:"spoiler"`
	TokenPath         string    `toml:"token_path"`
	TraktID           int       `toml:"trakt_id"`
	Type              string    `toml:"type"`
	UserName          string    `toml:"username"`
	Verbose           bool      `toml:"verbose"`
	WarningCode       int       `toml:"warningCode"`
}

var (
	versionFlag bool
)

// InitConfig of app
func InitConfig(fs afero.Fs) (*Config, error) {
	flagMap := make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		flagMap[f.Name] = f.Value.String()
	})

	if len(flagMap["c"]) == consts.ZeroValue {
		return nil, errors.New("config file not exists")
	}

	configFromFile, err := ReadConfigFromFile(fs, flagMap["c"])
	if err != nil {
		return nil, fmt.Errorf("init config error : %w", err)
	}

	return MergeConfigs(DefaultConfig(), configFromFile, flagMap)
}

// GenUsedFlagMap map of used flags
func GenUsedFlagMap() map[string]bool {
	flagset := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		key := string(f.Name[0])
		flagset[key] = true
	})

	return flagset
}

// MergeConfigs from two sources file and flags
func MergeConfigs(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string) (*Config, error) {
	flagset := GenUsedFlagMap()

	tokenPath, err := processOptionTokenPath(defaultConfig, fileConfig, flagConfig, flagset)
	if err != nil {
		return nil, fmt.Errorf("config error : %w", err)
	}
	defaultConfig.TokenPath = tokenPath
	defaultConfig.ClientID = processOptionClientID(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.ClientSecret = processOptionClientSecret(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.RedirectURI = processOptionRedirectURI(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.ErrorCode = processOptionErrorCode(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.WarningCode = processOptionWarningCode(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.PerPage = processOptionPerPage(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.PagesLimit = processOptionPagesLimit(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Verbose = processOptionVerbose(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.ConfigPath = processOptionConfigPath(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Output = processOptionOutput(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Type = processOptionType(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Format = processOptionFormat(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.UserName = processOptionUsername(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.List = processOptionList(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.ID = processOptionID(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Module = processOptionModule(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Action = processOptionAction(defaultConfig, fileConfig, flagConfig, flagset)
	defaultConfig.Sort = processOptionSort(defaultConfig, fileConfig, flagConfig, flagset)

	err = normalizeConfig(defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("config error : %w", err)
	}

	return defaultConfig, nil
}

func processOptionPagesLimit(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) int {
	// process if field is set in config file
	if fileConfig.PagesLimit > consts.ZeroValue && fileConfig.PagesLimit != defaultConfig.PagesLimit {
		defaultConfig.PagesLimit = fileConfig.PagesLimit
	}
	return defaultConfig.PagesLimit
}

func processOptionClientID(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) string {
	// Use values from fileConfig if present
	if len(fileConfig.ClientID) > consts.ZeroValue && fileConfig.ClientID != defaultConfig.ClientID {
		defaultConfig.ClientID = fileConfig.ClientID
	}
	return defaultConfig.ClientID
}

func processOptionClientSecret(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.ClientSecret) > consts.ZeroValue && fileConfig.ClientSecret != defaultConfig.ClientSecret {
		defaultConfig.ClientSecret = fileConfig.ClientSecret
	}
	return defaultConfig.ClientSecret
}

func processOptionRedirectURI(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.RedirectURI) > consts.ZeroValue && fileConfig.RedirectURI != defaultConfig.RedirectURI {
		defaultConfig.RedirectURI = fileConfig.RedirectURI
	}
	return defaultConfig.RedirectURI
}

func processOptionErrorCode(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) int {
	// process if field is set in config file
	if fileConfig.ErrorCode != consts.ZeroValue {
		defaultConfig.ErrorCode = fileConfig.ErrorCode
	}
	return defaultConfig.ErrorCode
}

func processOptionWarningCode(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) int {
	// process if field is set in config file
	if fileConfig.WarningCode != consts.ZeroValue {
		defaultConfig.WarningCode = fileConfig.WarningCode
	}
	return defaultConfig.WarningCode
}

func processOptionPerPage(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) int {
	// process if field is set in config file
	if fileConfig.PerPage > consts.ZeroValue && fileConfig.PerPage != defaultConfig.PerPage {
		defaultConfig.PerPage = fileConfig.PerPage
	}
	return defaultConfig.PerPage
}

func processOptionTokenPath(defaultConfig *Config, fileConfig *Config, _ map[string]string, _ map[string]bool) (string, error) {
	// process if field is set in config file
	if len(fileConfig.TokenPath) > consts.ZeroValue && fileConfig.TokenPath != defaultConfig.TokenPath {
		defaultConfig.TokenPath = fileConfig.TokenPath
	}

	tokenPath, err := expandTilde(defaultConfig.TokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand tilde from tokenPath: %w", err)
	}
	return tokenPath, nil
}

func processOptionVerbose(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, _ map[string]bool) bool {
	// process if field is set in config file
	if fileConfig.Verbose {
		defaultConfig.Verbose = fileConfig.Verbose
	}
	// process if flag is set
	f := "v"
	boolValue, err := strconv.ParseBool(flagConfig[f])
	if err == nil {
		defaultConfig.Verbose = boolValue
	}
	return defaultConfig.Verbose
}

func processOptionSort(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.Sort) > consts.ZeroValue && fileConfig.Sort != defaultConfig.Sort {
		defaultConfig.Sort = fileConfig.Sort
	}
	// process if flag is set
	f := "s"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Sort = flagConfig[f]
	} else if fileConfig.Sort != consts.EmptyString {
		defaultConfig.Sort = fileConfig.Sort
	}
	return defaultConfig.Sort
}

func processOptionAction(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.Action) > consts.ZeroValue && fileConfig.Action != defaultConfig.Action {
		defaultConfig.Action = fileConfig.Action
	}
	// process if flag is set
	f := "a"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Action = flagConfig[f]
	} else if fileConfig.Action != consts.EmptyString {
		defaultConfig.Action = fileConfig.Action
	}
	return defaultConfig.Action
}

func processOptionModule(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.Module) > consts.ZeroValue && fileConfig.Module != defaultConfig.Module {
		defaultConfig.Module = fileConfig.Module
	}
	// process if flag is set
	f := "m"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Module = flagConfig[f]
	} else if fileConfig.Module != consts.EmptyString {
		defaultConfig.Module = fileConfig.Module
	}
	return defaultConfig.Module
}

func processOptionID(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if fileConfig.ID == consts.EmptyString && fileConfig.ID != defaultConfig.ID {
		defaultConfig.ID = fileConfig.ID
	}
	// process if flag is set
	f := "i"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.ID = flagConfig[f]
	} else if fileConfig.ID != consts.EmptyString {
		defaultConfig.ID = fileConfig.ID
	}
	return defaultConfig.ID
}

func processOptionList(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.List) > consts.ZeroValue && fileConfig.List != defaultConfig.List {
		defaultConfig.List = fileConfig.List
	}
	// process if flag is set
	f := "l"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.List = flagConfig[f]
	} else if fileConfig.List != consts.EmptyString {
		defaultConfig.List = fileConfig.List
	}
	return defaultConfig.List
}

func processOptionUsername(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.UserName) > consts.ZeroValue && fileConfig.UserName != defaultConfig.UserName {
		defaultConfig.UserName = fileConfig.UserName
	}
	// process if flag is set
	f := "u"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.UserName = flagConfig[f]
	} else if fileConfig.UserName != consts.EmptyString {
		defaultConfig.UserName = fileConfig.UserName
	}

	return defaultConfig.UserName
}

func processOptionFormat(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.Format) > consts.ZeroValue && fileConfig.Format != defaultConfig.Format {
		defaultConfig.Format = fileConfig.Format
	}
	// process if flag is set
	f := "f"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Format = flagConfig[f]
	} else if fileConfig.Format != consts.EmptyString {
		defaultConfig.Format = fileConfig.Format
	}
	return defaultConfig.Format
}

func processOptionType(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.Type) > consts.ZeroValue && fileConfig.Type != defaultConfig.Type {
		defaultConfig.Type = fileConfig.Type
	}
	// process if flag is set
	f := "t"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Type = flagConfig[f]
	} else if fileConfig.Type != consts.EmptyString {
		defaultConfig.Type = fileConfig.Type
	}
	return defaultConfig.Type
}

func processOptionOutput(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.Output) > consts.ZeroValue && fileConfig.Output != defaultConfig.Output {
		defaultConfig.Output = fileConfig.Output
	}
	// process if flag is set
	f := "o"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Output = flagConfig[f]
	} else if fileConfig.Output != consts.EmptyString {
		defaultConfig.Output = fileConfig.Output
	}
	return defaultConfig.Output
}

func processOptionConfigPath(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string, flagset map[string]bool) string {
	// process if field is set in config file
	if len(fileConfig.ConfigPath) > consts.ZeroValue && fileConfig.ConfigPath != defaultConfig.ConfigPath {
		defaultConfig.ConfigPath = fileConfig.ConfigPath
	}
	// process if flag is set
	f := "c"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.ConfigPath = flagConfig[f]
	} else if fileConfig.ConfigPath != consts.EmptyString {
		defaultConfig.ConfigPath = fileConfig.ConfigPath
	}
	return defaultConfig.ConfigPath
}

// ReadConfigFromFile reads config from file stored on disc
func ReadConfigFromFile(fs afero.Fs, filename string) (*Config, error) {
	var config Config

	file, err := afero.ReadFile(fs, filename)

	if err != nil {
		return nil, fmt.Errorf("cannot read the config file : %w", err)
	}

	if len(string(file)) == consts.ZeroValue {
		return nil, errors.New("empty file content")
	}

	_, err = toml.Decode(string(file), &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse the config file : %w", err)
	}

	return &config, nil
}

// GetConfig yields the configuration
func GetConfig(fs afero.Fs, configPath string) (*Config, error) {
	config := &Config{}
	switch {
	case configPath != consts.EmptyString:
		err := parseConfig(fs, configPath, config)
		if err != nil {
			return nil, err
		}

	default: // no configuration provided
		config = DefaultConfig()
	}

	err := normalizeConfig(config)
	if err != nil {
		return nil, fmt.Errorf("config normalize error : %s", err)
	}

	return config, nil
}

func parseConfig(fs afero.Fs, path string, config *Config) error {
	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return errors.New("cannot read the config file")
	}
	_, err = toml.Decode(string(file), config)
	if err != nil {
		return fmt.Errorf("cannot parse the config file: %v", err)
	}

	return nil
}

// DefaultConfig config with default values
func DefaultConfig() *Config {
	return &Config{
		ClientID:       consts.EmptyString,
		ClientSecret:   consts.EmptyString,
		RedirectURI:    consts.EmptyString,
		WarningCode:    consts.ZeroValue,
		ErrorCode:      consts.ZeroValue,
		Verbose:        false,
		TokenPath:      consts.EmptyString,
		ConfigPath:     buildDefaultConfigPath(),
		Output:         consts.EmptyString,
		Format:         "imdb",
		Module:         "history",
		Action:         consts.EmptyString,
		Type:           "movies",
		CommentType:    "all",
		SearchIDType:   "trakt",
		SearchType:     []string{},
		SearchField:    []string{},
		Sort:           "rank",
		MoviesSort:     consts.EmptyString,
		MoviesType:     consts.EmptyString,
		MoviesPeriod:   "weekly",
		MoviesCountry:  consts.EmptyString,
		MoviesLanguage: consts.EmptyString,
		List:           "history",
		UserName:       "me",
		Hide:           false,
		ID:             consts.EmptyString,
		PerPage:        consts.DefaultPerPage,
		PagesLimit:     consts.PagesLimit,
		Progress:       consts.DefaultProgress,
		Remove:         false,
		Delete:         false,
		Spoiler:        false,
		Privacy:        "private",
		IncludeReplies: consts.EmptyString,
		Msg:            consts.EmptyString,
		InternalID:     consts.EmptyString,
		NotesID:        consts.ZeroValue,
		Item:           consts.EmptyString,
	}
}

func normalizeConfig(config *Config) error {
	if len(config.ClientID) == consts.ZeroValue || len(config.ClientSecret) == consts.ZeroValue {
		return errors.New("client_id and client_secret are required fields, update your config file")
	}

	if len(config.TokenPath) == consts.ZeroValue || (config.TokenPath != consts.EmptyString && !strings.HasSuffix(config.TokenPath, "json")) {
		return errors.New("token_path should be json file, update your config file")
	}

	return nil
}

func expandTilde(path string) (string, error) {
	if len(path) > consts.ZeroValue && path[consts.ZeroValue] == '~' {
		usr, err := user.Current()
		if err != nil {
			return consts.EmptyString, err
		}
		return filepath.Join(usr.HomeDir, path[1:]), nil
	}
	return path, nil
}

func buildDefaultConfigPath() string {
	absPath, err := expandTilde("~/trakt-sync.toml")
	if err != nil {
		panic(err)
	}
	return absPath
}
