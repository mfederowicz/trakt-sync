// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
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
	allHandlers := map[string]handlers.Handler{
		"movies":handlers.CertificationsTypesHandler{},
		"shows":handlers.CertificationsTypesHandler{},	
	}

	handler, err := cmd.common.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{"movies","shows"}
	if err != nil {
		cmd.common.GenTypeUsage(cmd.Name, validTypes)
		return nil
	}

	err = handler.Handle(options, client)
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
