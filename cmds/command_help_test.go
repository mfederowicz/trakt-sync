// Package cmds used for commands modules
package cmds

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"

	"github.com/spf13/afero"
)

func TestHelp(t *testing.T) {
	
	AppFs = afero.NewMemMapFs()
	tmpPath := "/tmp-iofs/"
	AppFs.MkdirAll(tmpPath, consts.X755)

	afero.WriteFile(AppFs, tmpPath+"/token.json", []byte("{}"), consts.X644)

	c := cfg.DefaultConfig()
	c.ClientID = "a"
	c.ClientSecret = "b"
	c.TokenPath = tmpPath + "/token.json"
	tests := []struct {
		Fs     afero.Fs
		Desc   string
		Client *internal.Client
		Config *cfg.Config
		Args   []string
		Regex  []string
	}{
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

	for _, test := range tests {
		buf := new(bytes.Buffer)
		stdout = buf
		HelpCmd.Exec(test.Fs, test.Client, test.Config, test.Args)
		out := buf.String()
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
}
