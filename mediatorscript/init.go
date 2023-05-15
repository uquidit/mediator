package mediatorscript

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func Init(fname string, l MsLogger) error {

	if fname == "" {
		return fmt.Errorf("cannot init mediatorscript package: no file name")
	} else {
		filename = fname
	}

	if l == nil {
		return fmt.Errorf("cannot init mediatorscript package: no logger")
	} else {
		logger = l
	}

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

	return nil
}
