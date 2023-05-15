package mediatorscript

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TestAllScripts(c echo.Context) error {
	return TestAllScriptsByTypeAndName(ScriptAll, "", c)
}

func TestScript(c echo.Context) error {

	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v: %v", ErrUnknownScriptType, slug)})

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		logger.Warningf("could not get script type from slug '%s': %v", slug, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("could not get script type from slug '%s': %v", slug, err)})

	} else if scriptname := c.Param("script"); scriptname != "" {
		// a name was specified.
		// let's check that script exists AND has the provided type
		if script, err := GetScriptByName(scriptname); err != nil {
			// could not find any script with that name: complain
			if errors.Is(err, ErrScriptNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("%v: '%s'", err, scriptname)})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error while searching script '%s': %v", scriptname, err)})
			}

		} else if script.Type != t {
			// ok, we've got a script but it's not the expected type: complain
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Script '%s' is not a %s", scriptname, t)})
		}
		return TestAllScriptsByTypeAndName(t, scriptname, c)

	} else {
		return TestAllScriptsByTypeAndName(t, "", c)
	}
}

func ExecuteScript(c echo.Context) error {
	var ti TicketInfo
	scriptname := c.Param("script")

	if err := c.Bind(&ti); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while processing parameters: %v", err.Error())})
	}
	if script, err := GetScriptByName(scriptname); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while processing parameters: %v", err.Error())})

	} else if err := script.AsyncRun(&ti); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while processing parameters: %v", err.Error())})

	} else {
		return c.NoContent(http.StatusNoContent)
	}
}

func RegisterScript(c echo.Context) error {
	var s Script
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while processing parameters: %s", err.Error())})
	}
	if err := s.Save(); err != nil {
		if registerErrorIsBadRequest(err) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("cannot save script: %s", err.Error())})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("cannot save script: %s", err.Error())})
		}

	}
	return c.NoContent(http.StatusNoContent)
}

func UnregisterScript(c echo.Context) error {
	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v: %v", ErrUnknownScriptType, slug)})

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		logger.Warningf("could not get script type from slug '%s': %v", slug, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("could not get script type from slug '%s': %v", slug, err)})

	} else if scriptname := c.Param("script"); scriptname != "" {
		// a name was specified.
		// let's check that script exists AND has the provided type
		if script, err := GetScriptByName(scriptname); err != nil {
			// could not find any script with that name: complain
			if errors.Is(err, ErrScriptNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("%v: '%s'", err, scriptname)})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error while searching script '%s': %v", scriptname, err)})
			}

		} else if script.Type != t {
			// ok, we've got a script but it's not the expected type: complain
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Script '%s' is not a %s", scriptname, t)})
		}

		if err := RemoveScriptByName(scriptname); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while unregitering script named %s: %v", scriptname, err)})
		} else {
			return c.NoContent(http.StatusNoContent)
		}

	} else if err := RemoveScriptByType(t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while unregistering %s: %v", t, err)})

	} else {
		return c.NoContent(http.StatusNoContent)
	}
}

func UnregisterAll(c echo.Context) error {
	if err := RemoveScriptByType(ScriptAll); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while unregistering all scripts: %v", err)})

	} else {
		return c.NoContent(http.StatusNoContent)
	}
}

func GetAll(c echo.Context) error {
	list := make([]*Script, 0, len(allScripts))
	for _, s := range allScripts {
		list = append(list, s)
	}
	return c.JSON(http.StatusOK, list)
}

func RefreshScript(c echo.Context) error {
	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v: %v", ErrUnknownScriptType, slug)})

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		logger.Warningf("could not get script type from slug '%s': %v", slug, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("could not get script type from slug '%s': %v", slug, err)})

	} else if scriptname := c.Param("script"); scriptname != "" {
		// a name was specified.
		// let's check that script exists AND has the provided type
		if script, err := GetScriptByName(scriptname); err != nil {
			// could not find any script with that name: complain
			if errors.Is(err, ErrScriptNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("%v: '%s'", err, scriptname)})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error while searching script '%s': %v", scriptname, err)})
			}

		} else if script.Type != t {
			// ok, we've got a script but it's not the expected type: complain
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Script '%s' is not a %s", scriptname, t)})

		} else if err := script.Refresh(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error while refreshing %s: %v", script, err)})
		}

	} else {
		l := GetScriptByType(t)
		for _, s := range l {
			if err := s.Refresh(); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error while refreshing %s: %v", s, err)})
			}
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func RefreshAllScript(c echo.Context) error {
	for _, s := range allScripts {
		if err := s.Refresh(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error while refreshing %s: %v", s, err)})
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func ExecuteScriptedCondition(c echo.Context) error {
	if ticket_id := c.Param("id"); ticket_id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ticket ID is missing"})
	} else {
		return genericHandler(ScriptCondition, ticket_id, c)
	}
}

func ExecuteScriptedTask(c echo.Context) error {
	if ticket_id := c.Param("id"); ticket_id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ticket ID is missing"})
	} else {
		return genericHandler(ScriptTask, ticket_id, c)
	}
}

func ExecutePreAssignment(c echo.Context) error {
	return genericHandler(ScriptAssignment, "", c)
}

func AddMediatorscriptAPI(g *echo.Group) {
	g.GET("", GetAll)

	g.POST("/register", RegisterScript)

	g.DELETE("/unregister-all", UnregisterAll)
	g.DELETE("/unregister/:slug/:script", UnregisterScript)
	g.DELETE("/unregister/:slug", UnregisterScript)

	g.POST("/refresh-all", RefreshAllScript)
	g.POST("/refresh/:slug/:script", RefreshScript)
	g.POST("/refresh/:slug", RefreshScript)

	g.POST("/execute/:script", ExecuteScript)
	g.POST("/execute-scripted-condition/:id", ExecuteScriptedCondition)
	g.POST("/execute-scripted-task/:id", ExecuteScriptedTask)
	g.POST("/execute-pre-assignment", ExecutePreAssignment)

	g.POST("/test-all", TestAllScripts)
	g.POST("/test/:slug/:script", TestScript)
	g.POST("/test/:slug", TestScript)

}
