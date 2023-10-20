package mediatorscript

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func TestAllScriptsByTypeAndName(script_type ScriptType, script_name string, res *RunResponse) {
	var (
		list ScriptList
	)

	if list = GetScriptByType(script_type); len(list) == 0 {
		res.statusCode = http.StatusNotFound
		res.err = ErrEmptyScriptList
		return
	}

	if res.RunResults == nil {
		res.RunResults = make(SyncRunResponsesMap)
	}

	for _, script := range list {
		if script_name != "" && script.Name != script_name {
			continue
		}

		logrus.Infof("Testing %s", script)
		res.RunResults[script.Name] = script.Test()
	}

	// return OK status because everythin went well on our side
	// we're not responsible for test failure
	res.statusCode = http.StatusOK
}
