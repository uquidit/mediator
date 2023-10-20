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

func Init(fname string) error {
	if fname == "" {
		return ErrInitNoFileName
	} else {
		filename = fname
		logrus.Infof("Mediatorscript package will use storage file '%s'", filename)
		allScripts = make(map[string]*Script)
		if content, err := os.ReadFile(filename); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		} else if len(content) == 0 {
			return fmt.Errorf("cannot read mediatorscript scripts: file '%s' is empty", filename)
		} else if err := json.Unmarshal(content, &allScripts); err != nil {
			return fmt.Errorf("cannot read mediatorscript scripts from file '%s': %w", filename, err)
		}
	}
	return nil
}
