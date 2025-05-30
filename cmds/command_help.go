// Package cmds used for commands modules
package cmds

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/spf13/afero"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Command vars
var (
	AppFs = afero.NewOsFs()
)

// HelpCmd shows help on the trakt-sync command and subcommands.
var HelpCmd = &Command{
	Name:    "help",
	Usage:   "[<commands>]",
	Summary: "Help on the trakt-sync command and subcommands.",
}

var helpDump = HelpCmd.Flag.Bool("godoc", false, "Dump the godoc output for the command(s)")

// HelpFunc shows help message for command
func HelpFunc(_ *Command, args ...string) error {
	selected := getSelectedFromArgs(args)

	switch {
	case *helpDump:
		return render(stdout, docTemplate, Commands)
	case len(selected) < len(args):
		printer.Fprintf(stdout, "error: unknown command %q\n", args[0])
		return render(stdout, helpTemplate, HelpCmd)
	case len(selected) == consts.ZeroValue:
		return render(stdout, usageTemplate, Commands)
	case len(selected) == consts.OneValue:
		return render(stdout, helpTemplate, selected[0])
	}
	return nil
}

func getSelectedFromArgs(args []string) []*Command {
	var selected []*Command

	if len(args) > consts.ZeroValue {
		want := strings.ToLower(args[0])
		for _, cmd := range Commands {
			if cmd.Name == want {
				selected = append(selected, cmd)
			}
		}
	}

	return selected
}

func init() {
	HelpCmd.Run = HelpFunc
}

func tabify(w io.Writer) *tabwriter.Writer {
	const (
		writerMinWidth = 0
		writerTabWidth = 0
		writerPadding  = 1
		writerPadChar  = ' '
		writerFlags    = 0
	)

	return tabwriter.NewWriter(w, writerMinWidth, writerTabWidth, writerPadding, writerPadChar, writerFlags)
}

var templateFuncs = template.FuncMap{
	"flags": func(indent int, args ...any) string {
		b := new(bytes.Buffer)
		prefix := strings.Repeat(" ", indent)
		w := tabify(b)
		visit := func(f *flag.Flag) {
			dash := "--"
			if len(f.Name) == 1 {
				dash = "-"
			}
			eq := "= " + f.DefValue
			switch typeName := fmt.Sprintf("%T", f.Value); {
			case typeName == "*flag.stringValue":
				// TODO(kevlar): make my own stringValue type so as to not depend on this?
				eq = fmt.Sprintf("= %q", f.DefValue)
			case f.DefValue == "":
				eq = ""
			}
			printer.Fprintf(w, "%s%s%s\t%s\t   %s\n", prefix, dash, f.Name, eq, f.Usage)
		}
		if len(args) == 0 {
			flag.VisitAll(visit)
		} else {
			args[0].(*Command).Flag.VisitAll(visit)
		}
		_ = w.Flush()
		if b.Len() == 0 {
			return ""
		}
		return fmt.Sprintf("\nOptions:\n%s", b)
	},
	"title": func(s string) string {
		titleCaser := cases.Title(language.Und) // "Und" stands for undetermined language
		return titleCaser.String(s + " command")
	},
	"trim": func(s string) string {
		return strings.TrimSpace(s)
	},
}

var stdout io.Writer = tabConverter{os.Stdout}

type tabConverter struct{ io.Writer }

func (t tabConverter) Write(p []byte) (int, error) {
	p = bytes.ReplaceAll(p, []byte{'\t'}, []byte{' ', ' ', ' ', ' '})
	return t.Writer.Write(p)
}

func render(w io.Writer, tpl string, data any) error {
	t := template.New("help")
	t.Funcs(templateFuncs)
	if err := template.Must(t.Parse(tpl)).Execute(w, data); err != nil {
		return fmt.Errorf("render error:%w", err)
	}
	return nil
}

var generalHelp = `	trakt-sync [<options>] [<command> [<suboptions>] [<arguments> ...]]
{{flags 2}}
Commands:{{range .}}
	{{.Name | printf "%-16s"}} {{.Summary}}{{end}}

Use "trakt-sync help <command>" for more help with a command.
`

var usageTemplate = `trakt-sync is a tool to sync data from your trakt account.

Usage:
` + generalHelp

var helpTemplate = `Usage: trakt-sync {{.Name}} [options]{{with .Usage}} {{.}}{{end}}{{if .Abbrev}}
       trakt-sync {{.Abbrev}} [options]{{with .Usage}} {{.}}{{end}}
{{end}}{{flags 2 .}}
{{.Summary}}
{{if .Help}}
{{.Help | trim}}{{end}}
`

var docTemplate = `
/*

Warning

This tool is under heavy development.  Don't depend on commands, options, or
pretty much anything else being stable yet.

Installation

As usual, the trakt-sync tool can be installed or upgraded via the "go" tool:
	go get -u github.com/mfederowicz/trakt-sync

General Usage

The trakt-sync command is composed of numerous sub-commands.
Sub-commands can be abbreviated to any unique prefix on the command-line.
The general usage is:

` + generalHelp + `

See below for a description of the various sub-commands understood by trakt-sync.
{{range .}}
{{.Name | title}}

{{.Summary | trim}}

Usage:
	trakt-sync {{.Name}} {{.Usage}}
{{flags 2 .}}
{{.Help | trim}}
{{end}}
*/
package main
`
