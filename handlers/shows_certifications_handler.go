// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsCertificationsHandler struct for handler
type ShowsCertificationsHandler struct{}

// Handle to handle shows: certifications action
func (m ShowsCertificationsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all content certifications for a show, including the country.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsCertifications(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found aliases for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsCertificationsHandler) fetchShowsCertifications(client *internal.Client, options *str.Options) ([]*str.Certification, *str.Response, error) {
	certifications, resp, err := client.Shows.GetAllShowCertifications(
		context.Background(),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return certifications, resp, nil
}
