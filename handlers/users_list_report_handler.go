// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersListReportHandler struct for handler
type UsersListReportHandler struct{ common CommonLogic }

// Handle to handle sync: list_report action
func (m UsersListReportHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.ValidReason(options)
	if err != nil {
		return err
	}
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if len(options.Msg) == consts.ZeroValue {
		return errors.New(consts.EmptyReportMsg)
	}

	resp, err := m.usersListReport(client, options)

	if err != nil {
		return fmt.Errorf("list report error:%w", err)
	}
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("list report success for list:", options.ListItemID)
	}

	return nil
}

func (UsersListReportHandler) usersListReport(client *internal.Client, options *str.Options) (*str.Response, error) {
	report := new(str.ListReport)
	report.Reason = &options.Reason
	report.Message = &options.Msg

	result, resp, err := client.Users.ListReport(
		client.BuildCtxFromOptions(options),
		&options.UserName,
		&options.ID,
		report)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("list not found:%s", options.ID)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("reason error:%s", *result.Message)
	}

	if resp.StatusCode == http.StatusConflict {
		return nil, fmt.Errorf("reason error:%s", *result.Message)
	}
	if err != nil {
		return nil, fmt.Errorf("report error:%w", err)
	}

	return resp, nil
}
