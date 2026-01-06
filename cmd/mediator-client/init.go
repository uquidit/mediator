package main

import (
	"fmt"
	"mediator/mediatorscript"

	"github.com/sirupsen/logrus"
)

func InitFromJSONIfNeeded(conf *mediatorscript.MediatorLegacyConfiguration, folder string) error {
	storageFileName := fmt.Sprintf("%s/mediator-client.json", folder)

	// check mediator JSON data store file
	if ok := checkStorage(storageFileName); !ok {
		logrus.Infof("mediator-client could not find JSON file '%s'. Using YAML file.", storageFileName)
		return nil
	}

	logrus.Infof("mediator-client will use workflows from JSON file: %s", storageFileName)

	// remove unwanted YAML workflows to use JSON workflows
	conf.Workflows = nil

	// load in memory
	if err := loadMediatorConfiguration(conf, storageFileName); err != nil {
		return err
	}

	return nil
}
