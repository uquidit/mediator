package main

import (
	"encoding/json"
	"errors"
	"os"
	"uqtu/mediator/logger"
	"uqtu/mediator/mediatorscript"
)

func checkStorage(filename string) bool {
	_, err := os.Open(filename)
	return !errors.Is(err, os.ErrNotExist)
}

func loadMediatorConfiguration(conf *mediatorscript.MediatorConfiguration, filename string) error {

	if content, err := os.ReadFile(filename); err != nil {
		logger.Warningf("Meditator configuration init: failed reading storage file: %v", err)
		return err
	} else if err := json.Unmarshal(content, &conf.Workflows); err != nil {
		logger.Warningf("Meditator configuration init: failed loading JSON data: %v", err)
		return err
	}

	return nil
}
