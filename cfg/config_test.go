package cfg

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	AppFs              = afero.NewMemMapFs()
	homeDirPath        = "/tmp-iofs/home/tester"
	configFileName     = "/trakt-sync.toml"
	configDirPath      = "/tmp-iofs/home/tester/.config"
	tokenFile          = "/token.json"
	bufferClientID     = "client_id = \"xxxxxxxxxxxxxxxxxxx\"\n"
	bufferClientSecret = "client_secret = \"xxxxxxxxxxxxxxxxxxxxxxxxxxxx\"\n"
	bufferTokenPath    = "token_path = \"\"\n"
	bufferPerPage      = "per_page = 50\n"
)

func commandLineArg(filename string) []string {
	return []string{"cmd", "-c=" + filename}
}

func addCFlag() {
	flag.String("c", "", "c")
}

func addDaysFlag() {
	flag.String("days", "", "days")
}

func emptyFlagset() *flag.FlagSet {
	return flag.NewFlagSet(consts.EmptyString, flag.ExitOnError)
}

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
	flag.CommandLine = emptyFlagset()
	flagset := make(map[string]bool)
	got := GenUsedFlagMap()
	if !test.MapsStringBoolEqual(got, flagset) {
		t.Errorf("maps should be equal")
	}

	// Recreate flags as needed
	os.Args = []string{"cmd", "--days=1"}
	addDaysFlag()
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
	flag.CommandLine = emptyFlagset()
	AppFs.MkdirAll(homeDirPath, consts.X755)
	filename := homeDirPath + configFileName

	os.Args = commandLineArg(filename)
	addCFlag()
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
	flag.CommandLine = emptyFlagset()
	AppFs.MkdirAll(homeDirPath, consts.X755)
	filename := homeDirPath + configFileName
	afero.WriteFile(AppFs, filename, []byte(consts.EmptyString), consts.X644)

	os.Args = commandLineArg(filename)
	addCFlag()
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
	flag.CommandLine = emptyFlagset()

	AppFs.MkdirAll(homeDirPath, consts.X755)
	filename := homeDirPath + configFileName
	afero.WriteFile(AppFs, filename, []byte("..."), consts.X644)

	os.Args = commandLineArg(filename)
	addCFlag()
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
	flag.CommandLine = emptyFlagset()

	AppFs.MkdirAll(homeDirPath, consts.X755)
	AppFs.MkdirAll(configDirPath, consts.X755)
	filenameToken := configDirPath + tokenFile

	afero.WriteFile(AppFs, filenameToken, []byte("\n"), consts.X644)

	filename := homeDirPath + configFileName
	var buffer bytes.Buffer

	// Write each line individually
	buffer.WriteString(bufferClientID)
	buffer.WriteString(bufferClientSecret)
	buffer.WriteString("token_path = \"" + filenameToken + "\"\n")
	buffer.WriteString(bufferPerPage)

	// Convert the buffer to a []byte
	data := buffer.Bytes()

	afero.WriteFile(AppFs, filename, data, consts.X644)
	os.Args = []string{consts.CMD, "-c=" + homeDirPath + configFileName}
	addCFlag()
	flag.Parse()
	c, _ := InitConfig(AppFs)
	assert.Equal(t, c.PerPage, consts.PerPage)
}

func TestInitConfigNoClient(t *testing.T) {
	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	// Reset the flag.CommandLine and reinitialize your flags
	flag.CommandLine = emptyFlagset()

	AppFs.MkdirAll(homeDirPath, consts.X755)

	AppFs.MkdirAll(configDirPath, consts.X755)

	filenameToken := configDirPath + tokenFile 

	afero.WriteFile(AppFs, filenameToken, []byte("\n"), consts.X644)

	filename := homeDirPath + configFileName
	var buffer bytes.Buffer

	// Write each line individually
	buffer.WriteString("client_id = \"\"\n")
	buffer.WriteString(bufferClientSecret)
	buffer.WriteString("token_path = \"" + filenameToken + "\"\n")
	buffer.WriteString(bufferPerPage)

	// Convert the buffer to a []byte
	data := buffer.Bytes()

	afero.WriteFile(AppFs, filename, data, consts.X644)
	os.Args = commandLineArg(homeDirPath + configFileName)
	addCFlag()
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
	flag.CommandLine = emptyFlagset()

	AppFs.MkdirAll(homeDirPath, consts.X755)

	AppFs.MkdirAll(configDirPath, consts.X755)

	filenameToken := configDirPath + tokenFile

	afero.WriteFile(AppFs, filenameToken, []byte(consts.NewLine), consts.X644)

	filename := homeDirPath + configFileName
	var buffer bytes.Buffer

	// Write each line individually
	buffer.WriteString(bufferClientID)
	buffer.WriteString(bufferClientSecret)
	buffer.WriteString(bufferTokenPath)
	buffer.WriteString(bufferPerPage)

	// Convert the buffer to a []byte
	data := buffer.Bytes()

	afero.WriteFile(AppFs, filename, data, consts.X644)
	os.Args = commandLineArg(homeDirPath + configFileName)
	addCFlag()
	flag.Parse()
	_, err := InitConfig(AppFs)
	assert.Contains(t, err.Error(), "token_path should be json file")
}
