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

// UsersReportHandler struct for handler
type UsersReportHandler struct{ common CommonLogic }

// Handle to handle users: report action
func (m UsersReportHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.ValidReason(options)
	if err != nil {
		return err
	}
	if len(options.UserName) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if len(options.Msg) == consts.ZeroValue {
		return errors.New(consts.EmptyReportMsg)
	}

	resp, err := m.usersReport(client, options)

	if err != nil {
		return fmt.Errorf("report error:%w", err)
	}
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("report success for user:", options.UserName)
	}

	return nil
}

func (UsersReportHandler) usersReport(client *internal.Client, options *str.Options) (*str.Response, error) {
	report := new(str.UserReport)
	report.Reason = &options.Reason
	report.Message = &options.Msg

	result, resp, err := client.Users.Report(
		client.BuildCtxFromOptions(options),
		&options.UserName,
		report)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found:%s", options.UserName)
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
