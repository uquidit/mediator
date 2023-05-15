package mediatorscript

import (
	"fmt"
	"io"

	"net/http"

	"github.com/labstack/echo/v4"
)

type SyncRunResponse struct {
	statusCode int
	Error      string     `json:"error"`
	ExitCode   int        `json:"exitcode"`
	StdOut     string     `json:"stdout"`
	StdErr     string     `json:"stderr"`
	Type       ScriptType `json:"type"`
}

// execute special script
// this function will only return an error if
// - inputs are not ok (bad request)
// - script could not be executed (bad config or system error)
// when the script fails, it will return an OK result
// script error will be in response body.
// it's the caller responsibility to check if run was ok
func genericHandler(t ScriptType, arg string, c echo.Context) error {
	var (
		res SyncRunResponse
	)

	if l := GetScriptByType(t); len(l) == 0 {
		res.statusCode = http.StatusBadRequest
		res.Error = fmt.Sprintf("%s script was not found: register it and try again", t)

	} else if len(l) > 1 {
		res.statusCode = http.StatusBadRequest
		res.Error = fmt.Sprintf("too many %s scripts have been found: make sure only one is registered and try again", t)

	} else {
		script := l[0]

		logger.Infof("Executing synchronously %s '%s' with arg '%s'", script.Type, script.Fullpath, arg)

		if req := c.Request(); req == nil {
			res.statusCode = http.StatusInternalServerError
			res.Error = "no request"

		} else if b, err := io.ReadAll(req.Body); err != nil {
			res.statusCode = http.StatusBadRequest
			res.Error = fmt.Sprintf("cannot read request body: %v", err)

		} else if stdout, stderr, err := script.SyncRun(b, arg); err != nil {

			if errorIsScriptFailure(err) {
				// script failure. not an internal error so return OK status
				res.statusCode = http.StatusOK
				res.Error = err.Error()
				res.ExitCode = getExitCodeFromError(err)
				res.StdOut = stdout
				res.StdErr = stderr
			} else {
				//error is not a script failure
				res.statusCode = http.StatusInternalServerError
				res.Error = err.Error()
			}

		} else {
			res.statusCode = http.StatusOK
			res.StdOut = stdout
		}
	}

	return c.JSON(res.statusCode, res)
}
