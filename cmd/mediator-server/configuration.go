package main

import (
	"uqtu/mediator/configparser"

	"github.com/sirupsen/logrus"
)

type Configurations struct {
	Server         ServerConfigurations   `json:"server"`
	Mediatorscript MediatorConfigurations `json:"mediatorscript"`
}
type MediatorscriptClientConfigurations struct {
	// full path of the generated configuration file (JSON format)
	SettingsFile   string `json:"settingsfile"`
	UploadScript   string `json:"uploadscript"`
	DownloadScript string `json:"downloadscript"`
}

type MediatorConfigurations struct {
	ScriptStorage       string                             `json:"scriptstorage"`
	ClientConfiguration MediatorscriptClientConfigurations `json:"clientconfiguration"`
}

type ServerConfigurations struct {
	Port   uint              `json:"port"`
	Host   string            `json:"host"`
	Log    LogConfigurations `json:"log"`
	Secret string            `json:"secret"`
	Ssl    SslConfigurations `json:"ssl"`
}

type LogConfigurations struct {
	Access string `json:"access"`
	Error  string `json:"error"`
}

type SslConfigurations struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
	Enabled     bool   `json:"enabled"`
}

var Configuration Configurations

func ReadConf(config_name string, verbose bool) {
	configparser.Verbose = verbose
	if err := configparser.ReadConf(config_name, &Configuration, nil); err != nil {
		logrus.Fatalf("Unable to decode into struct, %v", err)
	}

}

func ReadConfFromFile(abs_path string) error {
	if err := configparser.ReadConfAbsolutePath(abs_path, &Configuration, nil); err != nil {
		return err
	}
	return nil
}
