package mediatorscript

import (
	"fmt"
	"io"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// execute special script
// this function will only return an error if
// - inputs are not ok (bad request)
// - script could not be executed (bad config or system error)
// when the script fails, it will return an OK result
// script error will be in response body.
// it's the caller responsibility to check if run was ok
func genericHandler(t ScriptType, arg string, rr *RunResponse, c echo.Context) {

	if rr.RunResults == nil {
		rr.RunResults = make(SyncRunResponsesMap)
	}

	if l := GetScriptByType(t); len(l) == 0 {
		rr.statusCode = http.StatusBadRequest
		rr.err = fmt.Errorf("%s script was not found: register it and try again", t)

	} else if len(l) > 1 {
		rr.statusCode = http.StatusBadRequest
		rr.err = fmt.Errorf("too many %s scripts have been found: make sure only one is registered and try again", t)

	} else {
		script := l[0]

		logrus.Infof("Executing synchronously %s '%s' with arg '%s'", script.Type, script.Fullpath, arg)

		if req := c.Request(); req == nil {
			rr.statusCode = http.StatusInternalServerError
			rr.err = ErrNoRequest

		} else if b, err := io.ReadAll(req.Body); err != nil {
			rr.statusCode = http.StatusBadRequest
			rr.err = fmt.Errorf("cannot read request body: %v", err)

		} else {
			// execute script and store results in map
			rr.RunResults[script.Name] = script.SyncRun(b, arg)

		}
	}

}
