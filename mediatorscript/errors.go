package mediatorscript

import (
	"errors"
	"os/exec"
)

var (
	ErrEmptyScriptList                       = errors.New("empty script list: no script match criteria")
	ErrScriptNotFound                        = errors.New("script was not found")
	ErrUnknownScriptType                     = errors.New("unknown script type")
	ErrScriptExistForType                    = errors.New("a script has already been registered for that type")
	ErrRegisterNoFilename                    = errors.New("cannot register script: no filename")
	ErrRegisterNoName                        = errors.New("cannot register script: no name")
	ErrRegisterNameNotAllowed                = errors.New("cannot register script: 'test' is not an allowed name")
	ErrRegisterAlreadyExist                  = errors.New("script with same name already exists in registry")
	ErrRegisterAlreadyExistWithDifferentType = errors.New("script with same name BUT with different type already exists in registry")
	ErrInitNoFileName                        = errors.New("cannot init mediatorscript package: no file name")
	ErrInitNoLogger                          = errors.New("cannot init mediatorscript package: no logger")
	ErrMissingTicketID                       = errors.New("ticket ID is missing")
	ErrNoRequest                             = errors.New("no request")
	ErrHashMismatch                          = errors.New("script hash does not match")
	ErrExitCode                              = errors.New("script returned a non-zero exit code")
	ErrLastStep                              = errors.New("ticket has reached workflow last step. Cannot get next one")
	ErrScriptFileIsNotNormal                 = errors.New("script file is not a normal file or symlink")
	ErrScriptFileIsNotExecutable             = errors.New("script file is not executable")
	ErrScriptFileIsNotExecutableByBack       = errors.New("script file cannot be executed by back-end")
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
		errors.Is(err, ErrScriptExistForType) ||
		errors.Is(err, ErrScriptFileIsNotNormal) ||
		errors.Is(err, ErrScriptFileIsNotExecutable) ||
		errors.Is(err, ErrScriptFileIsNotExecutableByBack)
}
