package mediatorsettings

import "errors"

var (
	ErrNoUploadScript   error = errors.New("no upload script: upload settings to Securechange feature is disabled")
	ErrNoDownloadScript error = errors.New("no download script: download settings from Securechange feature is disabled")
	ErrNoSettingsFile   error = errors.New("no settings file for mediatorscript")
	ErrInvalidRule      error = errors.New("invalid rule")
)
