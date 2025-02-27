// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
)

// CertificationsCmd create or delete active checkins.
var CertificationsCmd = &Command{
	Name:    "certifications",
	Usage:   "",
	Summary: "Certifications list",
	Help:    `certifications command`,
}

func certificationsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.CertificationsHandler
	switch options.Type {
	case "movies", "shows":
		handler = handlers.CertificationsTypesHandler{}
	default:
		printer.Println("possible type: movies,shows")
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Type+":%s", err)
	}

	return nil
}

var (
	certificationsDumpTemplate = ``
)

func init() {
	CertificationsCmd.Run = certificationsFunc
}
