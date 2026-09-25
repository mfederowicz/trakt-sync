// Package cmds used for commands modules
package cmds

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"

	"github.com/spf13/afero"
)

type TestsList struct {
	Fs     afero.Fs
	Desc   string
	Client *internal.Client
	Config *cfg.Config
	Args   []string
	Regex  []string
}

func TestHelp(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	tmpPath := "/tmp-iofs/"
	AppFs.MkdirAll(tmpPath, consts.X755)

	afero.WriteFile(AppFs, tmpPath+"/token.json", []byte("{}"), consts.X644)
	afero.WriteFile(AppFs, tmpPath+"/user_settings.json", []byte("{}"), consts.X644)

	c := cfg.DefaultConfig()
	c.ClientID = "a"
	c.ClientSecret = "b"
	c.TokenPath = tmpPath + "/token.json"
	c.SettingsPath = tmpPath + "/user_settings.json"
	tests := genTestsList(c)
	processTests(t, tests)
}

func processTests(t *testing.T, tests []TestsList) {
	for _, test := range tests {
		buf := new(bytes.Buffer)
		stdout = buf
		err := HelpCmd.Exec(test.Fs, test.Client, test.Config, test.Args)
		if err != nil {
			t.Errorf("%s", err)
		}

		out := buf.String()
		processTestRegexp(test, t, out)
	}
}

func processTestRegexp(test TestsList, t *testing.T, out string) {
	for i, r := range test.Regex {
		matched, err := regexp.MatchString(r, out)
		if err != nil {
			t.Errorf("%s: %q: %s", r, err, r)
		}
		if !matched {
			t.Errorf("%s: regexp[%d] failed: %q\nOutput:\n%s", test.Desc, i, r, out)
		}
	}
}

func genTestsList(c *cfg.Config) []TestsList {
	return []TestsList{
		{
			Fs:     AppFs,
			Desc:   "Usage",
			Client: &internal.Client{},
			Config: c,
			Args:   []string{},
			Regex: []string{
				`^trakt-sync is a tool`,
				`\shelp.*Help on the trakt-sync`,
			},
		},
		{
			Fs:     AppFs,
			Desc:   "No Command",
			Client: &internal.Client{},
			Config: c,
			Args:   []string{"frobber"},
			Regex: []string{
				`error: unknown.*frobber`,
			},
		},
		{
			Fs:     AppFs,
			Desc:   "Help help",
			Client: &internal.Client{},
			Config: c,
			Args:   []string{"help"},
			Regex: []string{
				`[options] [command]`,
				`trakt-sync.*subcommand`,
				`--godoc`,
			},
		},
	}
}

// TestHelpCommandListAligned checks that every summary in the command list starts in the same column,
// also after a command name longer than the others (social_recommendations).
func TestHelpCommandListAligned(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := render(buf, generalHelp, Commands); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(buf.String(), "\n")
	column := -1
	for _, cmd := range Commands {
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 2 || fields[0] != cmd.Name {
				continue
			}
			start := strings.Index(line, cmd.Summary)
			if start < 0 {
				t.Fatalf("summary of %q not found in %q", cmd.Name, line)
			}
			if column == -1 {
				column = start
			}
			if start != column {
				t.Errorf("summary of %q starts at column %d, want %d:\n%s", cmd.Name, start, column, line)
			}
		}
	}
	if column == -1 {
		t.Fatal("command list not found in help output")
	}
}
