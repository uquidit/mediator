package mediatorsettings

import "errors"

var (
	ErrNoUploadScript           error = errors.New("no upload script: upload settings to Securechange feature is disabled")
	ErrNoDownloadScript         error = errors.New("no download script: download settings from Securechange feature is disabled")
	ErrNoSettingsFile           error = errors.New("no settings file for mediatorscript")
	ErrCannotDecodeSettingsFile error = errors.New("cannot decode settings file for mediatorscript")
	ErrCannotReadSettingsFile   error = errors.New("cannot read settings file for mediatorscript")
	ErrInvalidRule              error = errors.New("invalid rule")
	ErrInvalidSettings          error = errors.New("invalid settings")
	ErrNoWorkflowName           error = errors.New("no workflow name")
	ErrNoWorkflowID             error = errors.New("no workflow ID")
	ErrNoTriggerInRule          error = errors.New("no trigger in rule")
	ErrMissingStepInRule        error = errors.New("no step in rule but rule trigger requires a step")
	ErrUnknownScript            error = errors.New("missing or unknown script in rule")
	ErrScriptIsNotTriggerScript error = errors.New("rule script is not a trigger script")
)
