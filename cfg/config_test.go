package cfg

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/mfederowicz/trakt-sync/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	AppFs = afero.NewMemMapFs()
)

func TestMain(m *testing.M) {
	//AppFs := afero.NewMemMapFs()
	os.Unsetenv("HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Unsetenv("FORK")
	AppFs = afero.NewMemMapFs()
	os.Exit(m.Run())
}

func TestGenUsedFlagMap(t *testing.T) {
	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	flagset := make(map[string]bool)
	got := GenUsedFlagMap()
	if !test.MapsStringBoolEqual(got, flagset) {
		t.Errorf("maps should be equal")
	}

	// Recreate flags as needed
	os.Args = []string{"cmd", "--days=1"}
	flag.String("days", "", "days")
	flag.Parse()
	got2 := GenUsedFlagMap()
	var flagset2 = map[string]bool{
		"d": true,
	}

	if !test.MapsStringBoolEqual(got2, flagset2) {
		t.Errorf("maps should be equal:%v, got:%v", flagset2, got2)
	}

}

func TestInitConfigCannotRead(t *testing.T) {

	AppFs = afero.NewMemMapFs()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "config file not exist")
}

func TestInitConfigCannotReadOtherFile(t *testing.T) {

	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})
	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)
	filename := homeDirPath + "/trakt-sync.toml"

	os.Args = []string{"cmd", "-c=" + filename}
	flag.String("c", "", "c")
	flag.Parse()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "cannot read the config file")

}

func TestInitConfigNoContent(t *testing.T) {

	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})
	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)
	filename := homeDirPath + "/trakt-sync.toml"
	afero.WriteFile(AppFs, filename, []byte(""), 0644)

	os.Args = []string{"cmd", "-c=" + filename}
	flag.String("c", "", "c")
	flag.Parse()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "empty file content")

}

func TestInitConfigMalformedFile(t *testing.T) {

	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})
	
	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)
	filename := homeDirPath + "/trakt-sync.toml"
	afero.WriteFile(AppFs, filename, []byte("..."), 0644)

	os.Args = []string{"cmd", "-c=" + filename}
	flag.String("c", "", "c")
	flag.Parse()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "cannot parse the config file")

}

func TestInitConfigPerPageValue(t *testing.T) {

	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)

	configDirPath := "/tmp-iofs/home/tester/.config"
	AppFs.MkdirAll(configDirPath, 0755)

	filenameToken := configDirPath + "/token.json"

	afero.WriteFile(AppFs, filenameToken, []byte("\n"), 0644)

	filename := homeDirPath + "/trakt-sync.toml"
	var buffer bytes.Buffer

	// Write each line individually
	buffer.WriteString("client_id = \"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\"\n")
	buffer.WriteString("client_secret = \"xxxxxxxxxxxxxxxxxxxxxxxxxxxx\"\n")
	buffer.WriteString("token_path = \"" + filenameToken + "\"\n")
	buffer.WriteString("per_page = 50\n")

	// Convert the buffer to a []byte
	data := buffer.Bytes()

	afero.WriteFile(AppFs, filename, data, 0644)
	os.Args = []string{"cmd", "-c=" + homeDirPath + "/trakt-sync.toml"}
	flag.String("c", "", "c")
	flag.Parse()
	c, _ := InitConfig(AppFs)
	assert.Equal(t, c.PerPage, 50)

}

func TestInitConfigNoClient(t *testing.T) {

	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)

	configDirPath := "/tmp-iofs/home/tester/.config"
	AppFs.MkdirAll(configDirPath, 0755)

	filenameToken := configDirPath + "/token.json"

	afero.WriteFile(AppFs, filenameToken, []byte("\n"), 0644)

	filename := homeDirPath + "/trakt-sync.toml"
	var buffer bytes.Buffer

	// Write each line individually
	buffer.WriteString("client_id = \"\"\n")
	buffer.WriteString("client_secret = \"xxxxxxxxxxxxxxxxxxxxxxxxxxxx\"\n")
	buffer.WriteString("token_path = \"" + filenameToken + "\"\n")
	buffer.WriteString("per_page = 50\n")

	// Convert the buffer to a []byte
	data := buffer.Bytes()

	afero.WriteFile(AppFs, filename, data, 0644)
	os.Args = []string{"cmd", "-c=" + homeDirPath + "/trakt-sync.toml"}
	flag.String("c", "", "c")
	flag.Parse()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "client_id and client_secret are required fields")

}

func TestInitConfigNoTokenPath(t *testing.T) {

	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	homeDirPath := "/tmp-iofs/home/tester"
	AppFs.MkdirAll(homeDirPath, 0755)

	configDirPath := "/tmp-iofs/home/tester/.config"
	AppFs.MkdirAll(configDirPath, 0755)

	filenameToken := configDirPath + "/token.json"

	afero.WriteFile(AppFs, filenameToken, []byte("\n"), 0644)

	filename := homeDirPath + "/trakt-sync.toml"
	var buffer bytes.Buffer

	// Write each line individually
	buffer.WriteString("client_id = \"xxxxxxxxxxxxxxxxxxx\"\n")
	buffer.WriteString("client_secret = \"xxxxxxxxxxxxxxxxxxxxxxxxxxxx\"\n")
	buffer.WriteString("token_path = \"\"\n")
	buffer.WriteString("per_page = 50\n")

	// Convert the buffer to a []byte
	data := buffer.Bytes()

	afero.WriteFile(AppFs, filename, data, 0644)
	os.Args = []string{"cmd", "-c=" + homeDirPath + "/trakt-sync.toml"}
	flag.String("c", "", "c")
	flag.Parse()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "token_path should be json file")

}
