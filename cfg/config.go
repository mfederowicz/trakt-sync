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
	UserName     string    `toml:"username"`
	ConfigPath   string    `toml:"config_path"`
	Action       string    `toml:"action"`
	TokenPath    string    `toml:"token_path"`
	Type         string    `toml:"type"`
	Output       string    `toml:"output"`
	ID           string    `toml:"id"`
	SearchIDType string    `toml:"search_id_type"`
	RedirectURI  string    `toml:"redirect_uri"`
	ClientSecret string    `toml:"client_secret"`
	List         string    `toml:"list"`
	Format       string    `toml:"format"`
	ClientID     string    `toml:"client_id"`
	Query        string    `toml:"query"`
	Field        string    `toml:"field"`
	Sort         string    `toml:"sort"`
	Module       string    `toml:"module"`
	SearchField  str.Slice `toml:"search_field"`
	SearchType   str.Slice `toml:"search_type"`
	WarningCode  int       `toml:"warningCode"`
	ErrorCode    int       `toml:"errorCode"`
	Days         int       `toml:"days"`
	PerPage      int       `toml:"per_page"`
	Verbose      bool      `toml:"verbose"`
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
		return nil, fmt.Errorf("config file not exists")
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
	//fmt.Println(flagConfig)

	// Use values from fileConfig if present
	if len(fileConfig.ClientID) > consts.ZeroValue && fileConfig.ClientID != defaultConfig.ClientID {
		defaultConfig.ClientID = fileConfig.ClientID
	}

	if len(fileConfig.ClientSecret) > consts.ZeroValue && fileConfig.ClientSecret != defaultConfig.ClientSecret {
		defaultConfig.ClientSecret = fileConfig.ClientSecret
	}

	if len(fileConfig.RedirectURI) > consts.ZeroValue && fileConfig.RedirectURI != defaultConfig.RedirectURI {
		defaultConfig.RedirectURI = fileConfig.RedirectURI
	}

	if len(fileConfig.TokenPath) > consts.ZeroValue && fileConfig.TokenPath != defaultConfig.TokenPath {
		defaultConfig.TokenPath = fileConfig.TokenPath
	}

	tokenPath, err := expandTilde(defaultConfig.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand tilde from tokenPath: %w", err)
	}
	defaultConfig.TokenPath = tokenPath

	if len(fileConfig.ConfigPath) > consts.ZeroValue && fileConfig.ConfigPath != defaultConfig.ConfigPath {
		defaultConfig.ConfigPath = fileConfig.ConfigPath
	}

	if len(fileConfig.Output) > consts.ZeroValue && fileConfig.Output != defaultConfig.Output {
		defaultConfig.Output = fileConfig.Output
	}

	if len(fileConfig.Action) > consts.ZeroValue && fileConfig.Action != defaultConfig.Action {
		defaultConfig.Action = fileConfig.Action
	}

	if fileConfig.ErrorCode != consts.ZeroValue {
		defaultConfig.ErrorCode = fileConfig.ErrorCode
	}

	if fileConfig.WarningCode != consts.ZeroValue {
		defaultConfig.WarningCode = fileConfig.WarningCode
	}

	//fmt.Printf("%v",fileConfig.Verbose)

	if fileConfig.Verbose {
		defaultConfig.Verbose = fileConfig.Verbose
	}

	if len(fileConfig.Type) > consts.ZeroValue && fileConfig.Type != defaultConfig.Type {
		defaultConfig.Type = fileConfig.Type
	}

	if len(fileConfig.Sort) > consts.ZeroValue && fileConfig.Sort != defaultConfig.Sort {
		defaultConfig.Sort = fileConfig.Sort
	}

	if len(fileConfig.Module) > consts.ZeroValue && fileConfig.Module != defaultConfig.Module {
		defaultConfig.Module = fileConfig.Module
	}

	if len(fileConfig.Format) > consts.ZeroValue && fileConfig.Format != defaultConfig.Format {
		defaultConfig.Format = fileConfig.Format
	}

	if len(fileConfig.UserName) > consts.ZeroValue && fileConfig.UserName != defaultConfig.UserName {
		defaultConfig.UserName = fileConfig.UserName
	}

	if len(fileConfig.List) > consts.ZeroValue && fileConfig.List != defaultConfig.List {
		defaultConfig.List = fileConfig.List
	}

	if fileConfig.ID == consts.EmptyString && fileConfig.ID != defaultConfig.ID {
		defaultConfig.ID = fileConfig.ID
	}

	if fileConfig.PerPage > consts.ZeroValue && fileConfig.PerPage != defaultConfig.PerPage {
		defaultConfig.PerPage = fileConfig.PerPage
	}

	// Override with values from flagConfig, if present
	boolValue, err := strconv.ParseBool(flagConfig["v"])
	if err == nil {
		defaultConfig.Verbose = boolValue
	}

	f := "c"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.ConfigPath = flagConfig[f]
	} else if fileConfig.ConfigPath != consts.EmptyString {
		defaultConfig.ConfigPath = fileConfig.ConfigPath
	}

	f = "o"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Output = flagConfig[f]
	} else if fileConfig.Output != consts.EmptyString {
		defaultConfig.Output = fileConfig.Output
	}
	f = "t"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Type = flagConfig[f]
	} else if fileConfig.Type != consts.EmptyString {
		defaultConfig.Type = fileConfig.Type
	}

	f = "f"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Format = flagConfig[f]
	} else if fileConfig.Format != consts.EmptyString {
		defaultConfig.Format = fileConfig.Format
	}

	f = "u"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.UserName = flagConfig[f]
	} else if fileConfig.UserName != consts.EmptyString {
		defaultConfig.UserName = fileConfig.UserName
	}

	f = "l"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.List = flagConfig[f]
	} else if fileConfig.List != consts.EmptyString {
		defaultConfig.List = fileConfig.List
	}
	f = "i"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.ID = flagConfig[f]
	} else if fileConfig.ID != consts.EmptyString {
		defaultConfig.ID = fileConfig.ID
	}
	f = "m"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Module = flagConfig[f]
	} else if fileConfig.Module != consts.EmptyString {
		defaultConfig.Module = fileConfig.Module
	}
	f = "a"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Action = flagConfig[f]
	} else if fileConfig.Action != consts.EmptyString {
		defaultConfig.Action = fileConfig.Action
	}
	f = "s"
	if flagset[f] && flagConfig[f] != consts.EmptyString {
		defaultConfig.Sort = flagConfig[f]
	} else if fileConfig.Sort != consts.EmptyString {
		defaultConfig.Sort = fileConfig.Sort
	}

	err = normalizeConfig(defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("config error : %w", err)
	}

	return defaultConfig, nil
}

// ReadConfigFromFile reads config from file stored on disc
func ReadConfigFromFile(fs afero.Fs, filename string) (*Config, error) {
	var config Config

	file, err := afero.ReadFile(fs, filename)

	if err != nil {
		return nil, fmt.Errorf("cannot read the config file : %w", err)
	}

	if len(string(file)) == consts.ZeroValue {
		return nil, fmt.Errorf("empty file content")
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

	normalizeConfig(config)
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
		ClientID:     consts.EmptyString,
		ClientSecret: consts.EmptyString,
		RedirectURI:  consts.EmptyString,
		WarningCode:  consts.ZeroValue,
		ErrorCode:    consts.ZeroValue,
		Verbose:      false,
		TokenPath:    consts.EmptyString,
		ConfigPath:   buildDefaultConfigPath(),
		Output:       consts.EmptyString,
		Format:       "imdb",
		Module:       "history",
		Action:       consts.EmptyString,
		Type:         "movies",
		SearchIDType: "trakt",
		SearchType:   []string{},
		SearchField:  []string{},
		Sort:         "rank",
		List:         "history",
		UserName:     "me",
		ID:           consts.EmptyString,
		PerPage:      consts.DefaultPerPage,
	}
}

func normalizeConfig(config *Config) error {
	if len(config.ClientID) == consts.ZeroValue || len(config.ClientSecret) == consts.ZeroValue {
		return fmt.Errorf("client_id and client_secret are required fields, update your config file")
	}

	if len(config.TokenPath) == consts.ZeroValue || (config.TokenPath != consts.EmptyString && !strings.HasSuffix(config.TokenPath, "json")) {
		return fmt.Errorf("token_path should be json file, update your config file")
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
