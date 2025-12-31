package mediatorscript

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func TestAllScripts(c echo.Context) error {
	var res RunResponse
	TestAllScriptsByTypeAndName(ScriptAll, "", &res)
	return res.SendResponse(c)
}

func TestScript(c echo.Context) error {
	var (
		rr RunResponse
	)

	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		rr.err = fmt.Errorf("%w: %v", ErrUnknownScriptType, slug)
		rr.statusCode = http.StatusBadRequest

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		rr.err = fmt.Errorf("could not get script type from slug '%s': %w", slug, err)
		rr.statusCode = http.StatusInternalServerError

	} else if scriptname := c.Param("script"); scriptname != "" {
		// a name was specified.
		// let's check that script exists AND has the provided type
		if script, err := GetScriptByName(scriptname); err != nil {
			// could not find any script with that name: complain
			if errors.Is(err, ErrScriptNotFound) {
				rr.err = fmt.Errorf("%v: '%s'", err, scriptname)
				rr.statusCode = http.StatusNotFound

			} else {
				rr.err = fmt.Errorf("error while searching script '%s': %v", scriptname, err)
				rr.statusCode = http.StatusInternalServerError
			}

		} else if script.Type != t {
			// ok, we've got a script but it's not the expected type: complain
			rr.err = fmt.Errorf("script '%s' is not a %s", scriptname, t)
			rr.statusCode = http.StatusBadRequest

		} else {
			// will execute script in test mode and populate res
			// with execution results
			TestAllScriptsByTypeAndName(t, scriptname, &rr)
		}

	} else {
		// will execute scripts in test mode and populate res
		// with execution results
		TestAllScriptsByTypeAndName(t, "", &rr)
	}

	return rr.SendResponse(c)

}

func ExecuteScript(c echo.Context) error {
	var (
		ti  TicketInfo
		res RunResponse
	)
	scriptname := c.Param("script")

	if err := c.Bind(&ti); err != nil {
		res.Error = fmt.Sprintf("error while processing parameters: %v", err)
		return c.JSON(http.StatusBadRequest, res)

	} else if script, err := GetScriptByName(scriptname); err != nil {
		res.Error = fmt.Sprintf("cannot execute script '%s': %v", scriptname, err)
		logrus.Error(res.Error)
		return c.JSON(http.StatusBadRequest, res)

	} else if err := script.AsyncRun(&ti); err != nil {
		res.Error = fmt.Sprintf("error while executing script '%s': %v", scriptname, err)
		logrus.Error(res.Error)
		return c.JSON(http.StatusBadRequest, res)

	} else {
		return c.NoContent(http.StatusNoContent)
	}
}

func ExecuteScriptedCondition(c echo.Context) error {
	var rr RunResponse
	if ticket_id := c.Param("id"); ticket_id == "" {
		rr.err = ErrMissingTicketID
		rr.statusCode = http.StatusBadRequest
	} else {
		genericHandler(ScriptCondition, ticket_id, &rr, c)
	}
	return rr.SendResponse(c)
}

func ExecuteScriptedTask(c echo.Context) error {
	var rr RunResponse
	if ticket_id := c.Param("id"); ticket_id == "" {
		rr.err = ErrMissingTicketID
		rr.statusCode = http.StatusBadRequest
	} else {
		genericHandler(ScriptTask, ticket_id, &rr, c)
	}
	return rr.SendResponse(c)
}

func ExecutePreAssignment(c echo.Context) error {
	var rr RunResponse
	genericHandler(ScriptAssignment, "", &rr, c)
	return rr.SendResponse(c)
}

func ExecuteRiskAnalysis(c echo.Context) error {
	var rr RunResponse
	genericHandler(RiskAnalysis, "", &rr, c)

	return rr.SendResponse(c)
}
