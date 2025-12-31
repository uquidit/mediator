package mediatorsettings

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

var mutex sync.Mutex

func GetSettings(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := DownloadSettingsFileFromSecurechange(download_to_securechange_script, settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if settings, err := ReadWorkflowsSettings(settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	} else {
		res := settings.GetSlice()
		if err := editSteps(res, false, c); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		} else {

			return c.JSON(http.StatusOK, res)
		}
	}
}

func SetSettings(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	var (
		data MediatorSettings
	)

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// check settings
	for _, settings := range data {
		if err := settings.isValid(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	if err := editSteps(data, true, c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := WriteWorkflowsSettings(data, settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := UploadSettingsFileToSecurechange(upload_to_securechange_script, settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

func SetWorkflowSettings(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	var (
		settings    MediatorSettingsMap
		wf_settings WFSettings
		err         error
	)

	// get current settings
	if err := DownloadSettingsFileFromSecurechange(download_to_securechange_script, settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if settings, err = ReadWorkflowsSettings(settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := editSteps(settings, false, c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// get wf settings from body
	if err := c.Bind(&wf_settings); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid data in body: %w", err))
	}

	// check settings
	if err := wf_settings.isValid(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	settings[wf_settings.WFname] = &wf_settings

	res := settings.GetSlice()
	if err := editSteps(res, true, c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := WriteWorkflowsSettings(res, settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := UploadSettingsFileToSecurechange(upload_to_securechange_script, settings_filename); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}
