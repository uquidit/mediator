package mediatorscript

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Response format
type RunResponse struct {
	err        error
	Error      string              `json:"error"`
	RunResults SyncRunResponsesMap `json:"run_results"`
	statusCode int
}

type SyncRunResponse struct {
	scriptError   error
	ScriptError   string `json:"script_error"`
	internalError error
	InternalError string     `json:"internal_error"`
	ExitCode      int        `json:"exitcode"`
	StdOut        string     `json:"stdout"`
	StdErr        string     `json:"stderr"`
	Type          ScriptType `json:"type"`
}

type SyncRunResponsesMap map[string]*SyncRunResponse

func (rr *RunResponse) SendResponse(c echo.Context) error {
	if rr.err != nil {
		logrus.Warning(rr.err.Error())
		rr.Error = rr.err.Error()
	}
	for _, r := range rr.RunResults {
		if r.internalError != nil {
			r.InternalError = r.internalError.Error()
			logrus.Warning(r.InternalError)
		}
		if r.scriptError != nil {
			r.ScriptError = r.scriptError.Error()
		}
	}
	return c.JSON(rr.statusCode, rr)
}

// Return any error from SyncRunResponse
// return nil if everything went well
func (res *SyncRunResponse) GetError() error {
	if res.internalError != nil {
		return res.internalError
	}
	if res.scriptError != nil {
		return res.scriptError
	}
	if res.ExitCode != 0 {
		return ErrExitCode
	}
	return nil
}

func (res *SyncRunResponse) IsOK() bool {
	return res.GetError() == nil
}
