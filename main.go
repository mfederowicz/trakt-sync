// Package main github.com/mfederowicz/trakt-sync.
package main

import (
	"flag"
	"os"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/cmds"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

var (
	options      = &str.Options{}
	_verbose     = flag.Bool("v", false, cmds.VerboseUsage)
	_version     = flag.Bool("version", false, cmds.VersionUsage)
	_config_path = flag.String("c", cfg.DefaultConfig().ConfigPath, cmds.ConfigUsage)
)

func main() {

	config := cfg.InitConfig()
	client := internal.NewClient(nil)
	fs := afero.NewOsFs()
	options := cfg.OptionsFromConfig(fs, config)
	client.UpdateHeaders(options.Headers)

	flag.Usage = func() {
		cmds.HelpFunc(cmds.HelpCmd)
	}
	flag.Parse()

	if *_version {
		cli.GenAppVersion()
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	if !cli.ValidAccessToken(config, client.Oauth) {
		cli.PoolNewDeviceCode(config, client.Oauth)
	}

	cmds.ModulesRuntime(args, config, client, fs)
	os.Exit(0)

}
