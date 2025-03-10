// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CommentsTrendingHandler struct for handler
type CommentsTrendingHandler struct{ common CommonLogic }

// Handle to handle comments: trending action
func (h CommentsTrendingHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all comments with the most likes and replies over the last 7 days.")
	
	result, err := h.common.FetchTrendingComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch comments error:%v", err)
	}
	
	if len(result) == consts.ZeroValue {
		return errors.New("empty list")
	}
	
	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.CommentTrendingItem{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	
	writer.WriteJSON(options, jsonData)
	
	return nil
}
