// Package main github.com/mfederowicz/trakt-sync.
package main

import (
	"flag"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/cmds"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"

	"github.com/spf13/afero"
)

var (
	options     = &str.Options{}
	_verbose    = flag.Bool("v", false, consts.VerboseUsage)
	_version    = flag.Bool("version", false, consts.VersionUsage)
	_configPath = flag.String("c", cfg.DefaultConfig().ConfigPath, consts.ConfigUsage)
)

func main() {
	fs := afero.NewOsFs()
	config, err := cfg.InitConfig(fs)
	if err != nil {
		printer.Printf("Error: %v\n", err)
		return
	}
	client := internal.NewClient(nil)
	options, err := cfg.OptionsFromConfig(fs, config)
	if err != nil {
		printer.Printf("Error: %v\n", err)
		return
	}

	args, noflags := handleArgs()
	if noflags {
		return
	}

	client.UpdateHeaders(options.Headers)
	cli.HandleToken(fs, config, client, options)
	cmds.ModulesRuntime(args, fs, config, client)
}

func handleArgs() ([]string, bool) {
	flag.Usage = func() {
		cmds.HelpFunc(cmds.HelpCmd)
	}
	flag.Parse()

	if *_version {
		printer.Println(cli.GenAppVersion())
		return []string{}, true
	}

	args := flag.Args()
	if len(args) == consts.ZeroValue {
		flag.Usage()
		return []string{}, true
	}

	return args, false
}
