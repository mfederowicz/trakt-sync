// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsCertificationsHandler struct for handler
type ShowsCertificationsHandler struct{}

// Handle to handle shows: certifications action
func (m ShowsCertificationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all content certifications for a show, including the country.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsCertifications(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found aliases for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsCertificationsHandler) fetchShowsCertifications(client *trakt.Client, options *str.Options) ([]*str.Certification, *str.Response, error) {
	certifications, resp, err := client.Shows.GetAllShowCertifications(
		cli.ContextFromOptions(options),
		options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return certifications, resp, nil
}
