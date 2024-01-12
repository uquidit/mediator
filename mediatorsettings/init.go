package mediatorsettings

import "fmt"

var (
	settings_filename               string
	upload_to_securechange_script   string
	download_to_securechange_script string
)

const (
	DEFAULT_SETTINGS_FILENAME = "/tmp/mediator-client-settings.json"
)

func Init(settings_file, upl_script, dl_script string) []error {
	errs := []error{}

	if upl_script == "" {
		errs = append(errs, ErrNoUploadScript)
	} else {
		upload_to_securechange_script = upl_script
	}
	if dl_script == "" {
		errs = append(errs, ErrNoDownloadScript)
	} else {
		download_to_securechange_script = dl_script
	}

	if settings_file != "" {
		settings_filename = settings_file
	} else {
		errs = append(errs, fmt.Errorf("%w: will use %s", ErrNoSettingsFile, DEFAULT_SETTINGS_FILENAME))
		settings_filename = DEFAULT_SETTINGS_FILENAME
	}
	return errs
}
