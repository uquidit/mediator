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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error while processing parameters: %w", err))
	}
	if err := s.Save(); err != nil {
		if registerErrorIsBadRequest(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		} else if errors.Is(err, ErrRegisterAlreadyExistWithDifferentType) {
			return echo.NewHTTPError(http.StatusConflict, err)
		} else if errors.Is(err, ErrRegisterAlreadyExist) {
			return c.NoContent(http.StatusAlreadyReported)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("cannot save script: %w", err))
		}

	}
	return c.NoContent(http.StatusNoContent)
}

func UnregisterScript(c echo.Context) error {
	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%w: %s", ErrUnknownScriptType, slug))

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		msg := fmt.Sprintf("could not get script type from slug '%s': %v", slug, err)
		logrus.Warning(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New(msg))

	} else if scriptname := c.Param("script"); scriptname != "" {
		// a name was specified.
		// let's check that script exists AND has the provided type
		if script, err := GetScriptByName(scriptname); err != nil {
			// could not find any script with that name: complain
			if errors.Is(err, ErrScriptNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("%v: '%s'", err, scriptname))
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while searching script '%s': %w", scriptname, err))
			}

		} else if script.Type != t {
			// ok, we've got a script but it's not the expected type: complain
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Script '%s' is not a %s", scriptname, t))
		}

		if err := RemoveScriptByName(scriptname); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error while unregistering script named %s: %w", scriptname, err))
		} else {
			return c.NoContent(http.StatusNoContent)
		}

	} else if err := RemoveScriptByType(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error while unregistering %s: %w", t, err))

	} else {
		return c.NoContent(http.StatusNoContent)
	}
}

func UnregisterAll(c echo.Context) error {
	if err := RemoveScriptByType(ScriptAll); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error unregistering all scripts: %w", err))

	} else {
		return c.NoContent(http.StatusNoContent)
	}
}

func GetAll(c echo.Context) error {
	list := GetScriptByType(ScriptAll)
	return c.JSON(http.StatusOK, list)
}

func GetAllByType(c echo.Context) error {
	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%w: %v", ErrUnknownScriptType, slug))

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		msg := fmt.Sprintf("could not get script type from slug '%s': %v", slug, err)
		logrus.Warning(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New(msg))

	} else {

		list := GetScriptByType(t)
		return c.JSON(http.StatusOK, list)
	}
}

func RefreshScript(c echo.Context) error {
	//check slug is valid, complain otherwise
	if slug := c.Param("slug"); !IsScriptTypeSlug(slug) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%w: %v", ErrUnknownScriptType, slug))

	} else if t, err := GetTypeFromSlug(slug); err != nil { // get type from slug
		// this should never happen so return a internal server error
		msg := fmt.Sprintf("could not get script type from slug '%s': %v", slug, err)
		logrus.Warning(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New(msg))

	} else if scriptname := c.Param("script"); scriptname != "" {
		// a name was specified.
		// let's check that script exists AND has the provided type
		if script, err := GetScriptByName(scriptname); err != nil {
			// could not find any script with that name: complain
			if errors.Is(err, ErrScriptNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("%w: '%s'", err, scriptname))
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while searching script '%s': %w", scriptname, err))
			}

		} else if script.Type != t {
			// ok, we've got a script but it's not the expected type: complain
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Script '%s' is not a %s", scriptname, t))

		} else if err := script.Refresh(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while refreshing %s: %w", script, err))
		}

	} else {
		l := GetScriptByType(t)
		for _, s := range l {
			if err := s.Refresh(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while refreshing %s: %w", s, err))
			}
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func RefreshAllScript(c echo.Context) error {
	for _, s := range allScripts {
		if err := s.Refresh(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while refreshing %s: %w", s, err))
		}
	}
	return c.NoContent(http.StatusNoContent)
}
