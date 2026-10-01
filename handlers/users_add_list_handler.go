// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersAddListHandler struct for handler
type UsersAddListHandler struct{ common CommonLogic }

// Handle to handle users: lists action
func (u UsersAddListHandler) Handle(options *str.Options, client *trakt.Client) error {
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

	jsonDataResult, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write result to:" + options.Output)
	writer.WriteJSON(options, jsonDataResult)
	return nil
}
