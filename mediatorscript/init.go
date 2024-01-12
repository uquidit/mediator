package mediatorscript

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	allScripts = make(map[string]*Script)
}

func Init(storage string) error {
	if storage == "" {
		return ErrInitNoFileName
	} else {
		scriptStorageFilename = storage
		logrus.Infof("Mediatorscript package will use storage file '%s'", scriptStorageFilename)
		allScripts = make(map[string]*Script)
		if content, err := os.ReadFile(scriptStorageFilename); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		} else if len(content) == 0 {
			return fmt.Errorf("cannot read mediatorscript scripts: file '%s' is empty", scriptStorageFilename)
		} else if err := json.Unmarshal(content, &allScripts); err != nil {
			return fmt.Errorf("cannot read mediatorscript scripts from file '%s': %w", scriptStorageFilename, err)
		}
	}
	return nil
}
