// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mfederowicz/trakt-sync/cli"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersAddListHandler struct for handler
type UsersAddListHandler struct{ common CommonLogic }

// Handle to handle users: lists action
func (u UsersAddListHandler) Handle(options *str.Options, client *internal.Client) error {
	input, err := u.common.ReadInput(*options)
	if err != nil {
		return err
	}
	result, resp, err := u.common.UsersAddPersonalList(client, options, input.List)
	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}
	if err != nil {
		return fmt.Errorf("add personal list error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("new personal list created for:%s\n", options.UserName)
	}

	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}
