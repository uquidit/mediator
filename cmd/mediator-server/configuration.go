package main

import (
	"uqtu/mediator/configparser"

	"github.com/sirupsen/logrus"
)

type Configurations struct {
	Server         ServerConfigurations
	Mediatorscript MediatorConfigurations
}

type MediatorConfigurations struct {
	ScriptStorage string `json:"scriptstorage"`
}

type ServerConfigurations struct {
	Port   uint
	Host   string
	Log    LogConfigurations
	Secret string
	Ssl    SslConfigurations
}

type LogConfigurations struct {
	Access string
	Error  string
}

type SslConfigurations struct {
	Certificate string
	Key         string
	Enabled     bool
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
