package mediatorscript

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RegisterScript(c echo.Context) error {
	var s Script
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error while processing parameters: %s", err.Error())})
	}
	if err := s.Save(); err != nil {
		if registerErrorIsBadRequest(err) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("cannot save script: %s", err.Error())})
		} else if errors.Is(err, ErrRegisterAlreadyExist) {
			return c.NoContent(http.StatusAlreadyReported)
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
		logrus.Warningf("could not get script type from slug '%s': %v", slug, err)
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
		logrus.Warningf("could not get script type from slug '%s': %v", slug, err)
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
