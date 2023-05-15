package mediatorscript

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TestAllScriptsByTypeAndName(script_type ScriptType, script_name string, c echo.Context) error {
	var (
		list ScriptList
	)

	if list = GetScriptByType(script_type); len(list) == 0 {
		if script_type == ScriptAll {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No script was found"})
		} else {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("No %v was found", script_type)})
		}

	}

	results := map[string]*SyncRunResponse{}

	for _, script := range list {
		if script_name != "" && script.Name != script_name {
			continue
		}

		logger.Infof("Testing %s", script)
		results[script.String()] = script.Test()
	}
	// return OK status because everythin went well on our side
	// we're not responsible for test failure
	return c.JSON(http.StatusOK, results)
}
