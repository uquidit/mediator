package mediatorscript

import (
	"crypto/hmac"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Script struct {
	Fullpath string     `mapstructure:"fullpath" json:"fullpath"`
	Name     string     `mapstructure:"name" json:"name"`
	Hash     []byte     `mapstructure:"hash" json:"hash"`
	Type     ScriptType `mapstructure:"type" json:"type"`
}

type ScriptList []*Script

var (
	allScripts map[string]*Script
	filename   string
	logger     MsLogger
)

func GetScriptByName(name string) (*Script, error) {
	if s, exist := allScripts[name]; !exist {
		return nil, ErrScriptNotFound
	} else {
		return s, nil
	}
}

// return all the scripts of the given type in a slice
// if given type is ScriptAll, return all script
func GetScriptByType(t ScriptType) ScriptList {
	l := ScriptList{}
	for _, s := range allScripts {
		if s.Type == t || t == ScriptAll {
			l = append(l, s)
		}
	}
	return l
}

func IsEmpty(t ScriptType) bool {
	for _, s := range allScripts {
		if s.Type == t {
			return false
		}
	}
	return true
}

func RemoveScriptByName(name string) error {
	if _, exist := allScripts[name]; !exist {
		return fmt.Errorf("script '%s' does not exist", name)
	} else {
		delete(allScripts, name)
	}
	return save()
}

func RemoveScriptByType(t ScriptType) error {
	if t == ScriptAll && len(allScripts) > 0 {
		allScripts = make(map[string]*Script)
	} else {
		for _, s := range allScripts {
			if s.Type == t {
				delete(allScripts, s.Name)
			}
		}
	}
	return save()
}

func GetAllScriptNames() []string {
	keys := make([]string, 0, len(allScripts))
	for k, s := range allScripts {
		keys = append(keys, fmt.Sprintf("%s: %s", k, s.Fullpath))
	}
	return keys

}

func (s *Script) String() string {
	return fmt.Sprintf("%s '%s'", s.Type, s.Name)
}

func (s *Script) Save() error {
	var err error
	if s.Fullpath == "" {
		return ErrRegisterNoFilename
	}
	if s.Name == "" {
		return ErrRegisterNoName
	}
	if s.Name == "test" {
		return ErrRegisterNameNotAllowed
	}

	if s.Hash, err = s.computeHash(); err != nil {
		return err
	}

	if err := safeAdd(s.Name, s); err != nil {
		return err
	}

	return save()
}
func (s *Script) Refresh() error {
	var err error
	logger.Infof("Refreshing %s", s)
	if s.Hash, err = s.computeHash(); err != nil {
		return err
	}
	return save()
}

func safeAdd(name string, item *Script) error {
	if _, exist := allScripts[name]; exist {
		return fmt.Errorf("script '%s' already exist in list", name)

	} else if item.Type == ScriptTrigger || IsEmpty(item.Type) {
		// we can have several trigger scripts but only one of the other type
		// so now we know the name is not in use, just make sure the slot is empty if type is not trigger
		// If so, append the new script
		allScripts[name] = item
		return nil

	} else {
		return fmt.Errorf("%w: %s", ErrScriptExistForType, item.Type)

	}
}

func save() error {

	// marshall list into JSON
	if content, err := json.MarshalIndent(allScripts, "", " "); err != nil {
		return err
	} else {
		// save string to file
		return os.WriteFile(filename, content, 0644)
	}
}

func (s *Script) checkHash() error {
	if hash, err := s.computeHash(); err != nil {
		return err
	} else if !hmac.Equal(hash, s.Hash) {
		return fmt.Errorf("script hash does not match for %s script %s (%s)", s.Type, s.Name, s.Fullpath)
	}
	return nil
}

func (s *Script) AsyncRun(ti *TicketInfo) error {
	if err := s.checkHash(); err != nil {
		return err

	} else {
		f := s.getRunFunction()
		if data, err := xml.Marshal(ti); err != nil {
			return err
		} else {
			logger.Infof("running script %s (%s) with data: %s", s.Name, s.Fullpath, data)
			go f(data, "")
		}

		return nil
	}

}

// Run script in test mode
// return two error values:
// * first is any error returned by the script
// * second is any error occuring in that function
// So:
// - if first error is not nil, test failed
// - if second is not nil, run failed
func (s *Script) Test() *SyncRunResponse {
	var (
		res    SyncRunResponse
		stdout string
		stderr string
		err    error
	)
	res.Type = s.Type

	if err = s.checkHash(); err != nil {
		res.statusCode = http.StatusForbidden

	} else {
		f := s.getRunFunction()
		input := []byte("<ticket_info/>")
		arg := ""
		if s.Type == ScriptCondition || s.Type == ScriptTask {
			arg = "test"
		}
		if stdout, stderr, err = f(input, arg); err == nil {
			res.statusCode = http.StatusOK
			res.StdOut = stdout

		} else if errorIsScriptFailure(err) {
			// script failure. not an internal error so return OK status
			res.statusCode = http.StatusOK
			res.ExitCode = getExitCodeFromError(err)
			res.StdOut = stdout
			res.StdErr = stderr

		} else {
			//error is not a script failure
			res.statusCode = http.StatusInternalServerError

		}
	}
	if err != nil {
		res.Error = err.Error()
	}
	return &res
}

func (s *Script) getRunFunction() func([]byte, string) (string, string, error) {
	return func(input []byte, arg string) (string, string, error) {
		var (
			stdin          io.WriteCloser
			stdout, stderr strings.Builder
			err            error
			cmd            *exec.Cmd
		)
		if arg == "" {
			cmd = exec.Command(s.Fullpath)
		} else {
			cmd = exec.Command(s.Fullpath, arg)
		}

		// warm stdin up if we need to send data
		if input != nil {
			if stdin, err = cmd.StdinPipe(); err != nil {
				return "", "", err
			}
		}

		// initialize vars to stdout and stderr
		// so we can get whatever is sent by the script
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// start the script
		logger.Infof("Starting %s", s)
		if err := cmd.Start(); err != nil {
			out := strings.TrimSpace(stdout.String())
			er := strings.TrimSpace(stderr.String())
			logger.Warningf("stdout: %s", out)
			logger.Warningf("stderr: %s", er)
			logger.Warningf("err: %v", err)
			return out, er, err
		}

		// write input to stdin
		if _, err := stdin.Write(input); err != nil {
			out := strings.TrimSpace(stdout.String())
			er := strings.TrimSpace(stderr.String())
			logger.Warningf("stdout: %s", stdout.String())
			logger.Warningf("stderr: %s", stderr.String())
			logger.Warningf("err: %v", err)
			return out, er, err
		}

		// force pipe to close so script can run freely
		if err := stdin.Close(); err != nil {
			out := strings.TrimSpace(stdout.String())
			er := strings.TrimSpace(stderr.String())
			logger.Warningf("stdout: %s", stdout.String())
			logger.Warningf("stderr: %s", stderr.String())
			logger.Warningf("err: %v", err)
			return out, er, err
		}

		if err := cmd.Wait(); err != nil {
			out := strings.TrimSpace(stdout.String())
			er := strings.TrimSpace(stderr.String())
			logger.Warningf("stdout: %s", stdout.String())
			logger.Warningf("stderr: %s", stderr.String())
			logger.Warningf("err: %v", err)
			return out, er, err
		}

		out := strings.TrimSpace(stdout.String())
		er := strings.TrimSpace(stderr.String())
		logger.Infof("script %s was run successfully", s.Fullpath)
		logger.Infof("stdout: %s", stdout.String())
		logger.Infof("stderr: %s", stderr.String())
		return out, er, nil
	}
}

func (s *Script) SyncRun(input []byte, arg string) (string, string, error) {
	if err := s.checkHash(); err != nil {
		return "", "", err
	}
	if s.Type == ScriptTrigger {
		logger.Warningf("Trigger Script '%s' is run synchronously. Such scripts are usually run asynchronously.")
	}

	f := s.getRunFunction()
	return f(input, arg)
}
