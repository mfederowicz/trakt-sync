// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CertificationsTypesHandler interface to handle certifications types
type CertificationsTypesHandler struct{}

// Handle to handle calendars: shows action
func (CertificationsTypesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("certifications handler:" + options.Type)

	certifications, _, err := fetchCertifications(client, options)
	if err != nil {
		return fmt.Errorf("fetch certifications error:%w", err)
	}

	printer.Print("Found " + options.Type + " data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(certifications, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchCertifications(client *internal.Client, options *str.Options) (*str.Certifications, *str.Response, error) {
	results, resp, err := client.Certifications.GetCertifications(client.BuildCtxFromOptions(options), &options.Type)

	return results, resp, err
}
