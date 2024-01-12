package mediatorsettings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func ReadWorkflowsSettings(filename string) (MediatorSettingsMap, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return MediatorSettingsMap{}, nil
		}
		return nil, err
	}

	data := MediatorSettingsMap{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteWorkflowsSettings(wf_settings MediatorSettings, filename string) error {
	if buffer, err := json.MarshalIndent(wf_settings.GetMap(), "", "  "); err != nil {
		return fmt.Errorf("error encoding: %w", err)
	} else if err := os.WriteFile(filename, buffer, 0644); err != nil {
		return fmt.Errorf("error while writing data to file: %w", err)
	}
	return nil
}

func WriteWorkflowsSettingsFromMap(wf_settings MediatorSettingsMap, filename string) error {
	if buffer, err := json.MarshalIndent(wf_settings, "", "  "); err != nil {
		return fmt.Errorf("error encoding: %w", err)
	} else if err := os.WriteFile(filename, buffer, 0644); err != nil {
		return fmt.Errorf("error while writing data to file: %w", err)
	}
	return nil
}

func UploadSettingsFileToSecurechange(upload_script, filename string) error {
	if upload_script == "" {
		return ErrNoUploadScript
	}
	cmd := exec.Command("/usr/bin/sudo", upload_script, filename)
	// get stderr in a buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	logrus.Infof("Upload mediator-client settings to Securechange using command: %s %s", upload_script, filename)
	if err := cmd.Run(); err != nil {

		errmsg := stderr.String()
		if errmsg != "" {
			logrus.Warningf("Upload mediator-client settings returned an error: %s", errmsg)
		}

		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code is %d\n", exitError.ExitCode())
		}
		return err
	}
	return nil
}

func DownloadSettingsFileFromSecurechange(download_script, filename string) error {
	if download_script == "" {
		return ErrNoDownloadScript
	}
	cmd := exec.Command("/usr/bin/sudo", download_script, filename)
	// get stderr in a buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	logrus.Infof("Download mediator-client settings from Securechange using command: /usr/bin/sudo %s %s", download_script, filename)
	if err := cmd.Run(); err != nil {

		errmsg := stderr.String()
		if errmsg != "" {
			logrus.Warningf("Download mediator-client settings returned an error: %s", errmsg)
		}

		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code is %d\n", exitError.ExitCode())
		}
		return err
	}
	return nil
}
