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
func InitConfig() (*Config, error) {
	flagMap := make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		flagMap[f.Name] = f.Value.String()
	})
	fs := afero.NewOsFs()
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
		if len(f.Name) > 1 && f.Name[1] == '-' {
			key = f.Name[1:]
		}
		flagset[key] = true
	})

	return flagset

}

// MergeConfigs from two sources file and flags
func MergeConfigs(defaultConfig *Config, fileConfig *Config, flagConfig map[string]string) (*Config, error) {

	flagset := GenUsedFlagMap()
	//fmt.Println(flagConfig)

	// Use values from fileConfig if present
	if len(fileConfig.ClientID) > 0 && fileConfig.ClientID != defaultConfig.ClientID {
		defaultConfig.ClientID = fileConfig.ClientID
	}

	if len(fileConfig.ClientSecret) > 0 && fileConfig.ClientSecret != defaultConfig.ClientSecret {
		defaultConfig.ClientSecret = fileConfig.ClientSecret
	}

	if len(fileConfig.RedirectURI) > 0 && fileConfig.RedirectURI != defaultConfig.RedirectURI {
		defaultConfig.RedirectURI = fileConfig.RedirectURI
	}

	if len(fileConfig.TokenPath) > 0 && fileConfig.TokenPath != defaultConfig.TokenPath {
		defaultConfig.TokenPath = fileConfig.TokenPath
	}

	tokenPath, err := expandTilde(defaultConfig.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand tilde from tokenPath: %w", err)
	}
	defaultConfig.TokenPath = tokenPath

	if len(fileConfig.ConfigPath) > 0 && fileConfig.ConfigPath != defaultConfig.ConfigPath {
		defaultConfig.ConfigPath = fileConfig.ConfigPath
	}

	if len(fileConfig.Output) > 0 && fileConfig.Output != defaultConfig.Output {
		defaultConfig.Output = fileConfig.Output
	}

	if len(fileConfig.Action) > 0 && fileConfig.Action != defaultConfig.Action {
		defaultConfig.Action = fileConfig.Action
	}

	if fileConfig.ErrorCode != 0 {
		defaultConfig.ErrorCode = fileConfig.ErrorCode
	}

	if fileConfig.WarningCode != 0 {
		defaultConfig.WarningCode = fileConfig.WarningCode
	}

	//fmt.Printf("%v",fileConfig.Verbose)

	if fileConfig.Verbose {
		defaultConfig.Verbose = fileConfig.Verbose
	}

	if len(fileConfig.Type) > 0 && fileConfig.Type != defaultConfig.Type {
		defaultConfig.Type = fileConfig.Type
	}

	if len(fileConfig.Sort) > 0 && fileConfig.Sort != defaultConfig.Sort {
		defaultConfig.Sort = fileConfig.Sort
	}

	if len(fileConfig.Module) > 0 && fileConfig.Module != defaultConfig.Module {
		defaultConfig.Module = fileConfig.Module
	}

	if len(fileConfig.Format) > 0 && fileConfig.Format != defaultConfig.Format {
		defaultConfig.Format = fileConfig.Format
	}

	if len(fileConfig.UserName) > 0 && fileConfig.UserName != defaultConfig.UserName {
		defaultConfig.UserName = fileConfig.UserName
	}

	if len(fileConfig.List) > 0 && fileConfig.List != defaultConfig.List {
		defaultConfig.List = fileConfig.List
	}

	if fileConfig.ID == "" && fileConfig.ID != defaultConfig.ID {
		defaultConfig.ID = fileConfig.ID
	}

	if fileConfig.PerPage > 0 && fileConfig.PerPage != defaultConfig.PerPage {
		defaultConfig.PerPage = fileConfig.PerPage
	}

	// Override with values from flagConfig, if present
	boolValue, err := strconv.ParseBool(flagConfig["v"])
	if err == nil {
		defaultConfig.Verbose = boolValue
	}

	if flagset["c"] && flagConfig["c"] != "" {
		defaultConfig.ConfigPath = flagConfig["c"]
	} else if fileConfig.ConfigPath != "" {
		defaultConfig.ConfigPath = fileConfig.ConfigPath
	}

	if flagset["o"] && flagConfig["o"] != "" {
		defaultConfig.Output = flagConfig["o"]
	} else if fileConfig.Output != "" {
		defaultConfig.Output = fileConfig.Output
	}
	if flagset["t"] && flagConfig["t"] != "" {
		defaultConfig.Type = flagConfig["t"]
	} else if fileConfig.Type != "" {
		defaultConfig.Type = fileConfig.Type
	}
	if flagset["f"] && flagConfig["f"] != "" {
		defaultConfig.Format = flagConfig["f"]
	} else if fileConfig.Format != "" {
		defaultConfig.Format = fileConfig.Format
	}

	if flagset["u"] && flagConfig["u"] != "" {
		defaultConfig.UserName = flagConfig["u"]
	} else if fileConfig.UserName != "" {
		defaultConfig.UserName = fileConfig.UserName
	}

	if flagset["l"] && flagConfig["l"] != "" {
		defaultConfig.List = flagConfig["l"]
	} else if fileConfig.List != "" {
		defaultConfig.List = fileConfig.List
	}

	if flagset["i"] && flagConfig["i"] != "" {
		defaultConfig.ID = flagConfig["i"]
	} else if fileConfig.ID != "" {
		defaultConfig.ID = fileConfig.ID
	}

	if flagset["m"] && flagConfig["m"] != "" {
		defaultConfig.Module = flagConfig["m"]
	} else if fileConfig.Module != "" {
		defaultConfig.Module = fileConfig.Module
	}

	if flagset["a"] && flagConfig["a"] != "" {
		defaultConfig.Action = flagConfig["a"]
	} else if fileConfig.Action != "" {
		defaultConfig.Action = fileConfig.Action
	}

	if flagset["s"] && flagConfig["s"] != "" {
		defaultConfig.Sort = flagConfig["s"]
	} else if fileConfig.Sort != "" {
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
	case configPath != "":
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
		ClientID:     "",
		ClientSecret: "",
		RedirectURI:  "",
		WarningCode:  0,
		ErrorCode:    0,
		Verbose:      false,
		TokenPath:    "",
		ConfigPath:   buildDefaultConfigPath(),
		Output:       "",
		Format:       "imdb",
		Module:       "history",
		Action:       "",
		Type:         "movies",
		SearchIDType: "trakt",
		SearchType:   []string{},
		SearchField:  []string{},
		Sort:         "rank",
		List:         "history",
		UserName:     "me",
		ID:           "",
		PerPage:      100,
	}
}

func normalizeConfig(config *Config) error {

	if len(config.ClientID) == 0 || len(config.ClientSecret) == 0 {
		return fmt.Errorf("client_id and client_secret are required fields, update your config file")
	}

	if len(config.TokenPath) == 0 || (config.TokenPath != "" && !strings.HasSuffix(config.TokenPath, "json")) {
		return fmt.Errorf("token_path should be json file, update your config file")
	}

	return nil

}

func expandTilde(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return "", err
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
