package mediatorscript

import (
	"errors"
	"os/exec"
)

var (
	ErrScriptNotFound         = errors.New("script was not found")
	ErrUnknownScriptType      = errors.New("unknown script type")
	ErrScriptExistForType     = errors.New("a script has already been registered for that type")
	ErrRegisterNoFilename     = errors.New("cannot register script: no filename")
	ErrRegisterNoName         = errors.New("cannot register script: no name")
	ErrRegisterNameNotAllowed = errors.New("cannot register script: 'test' is not an allowed name")
)

func errorIsScriptFailure(err error) bool {
	_, ok := err.(*exec.ExitError)
	return ok
}

func getExitCodeFromError(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		return exiterr.ExitCode()
	} else {
		return 0
	}
}

func registerErrorIsBadRequest(err error) bool {
	return errors.Is(err, ErrRegisterNoFilename) ||
		errors.Is(err, ErrRegisterNoName) ||
		errors.Is(err, ErrRegisterNameNotAllowed) ||
		errors.Is(err, ErrScriptExistForType)
}
