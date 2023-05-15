package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/logger"
	"uqtu/mediator/mediatorscript"
	"uqtu/mediator/totp"
)

func runInteractiveScripts(scriptedCondition bool, preAssignment bool, scriptedTask bool, datafilenameFlag string, conf *mediatorscript.MediatorConfiguration) {
	if (scriptedCondition && (preAssignment || scriptedTask)) ||
		(preAssignment && (scriptedCondition || scriptedTask)) ||
		(scriptedTask && (scriptedCondition || preAssignment)) {
		logger.Errorf("Only one of the --scripted-condition --pre-assignment and --scripted-task flag can be used at a time.")
	}

	client := apiclient.NewClient(conf.Configuration.BackendURL, "", "", conf.Configuration.SSLSkipVerify)
	client.Token = totp.GetKey()
	res := mediatorscript.SyncRunResponse{}
	var (
		reqBody    io.Reader
		endpoint   string
		currScript string
	)

	args := flag.Args()
	switch {
	case scriptedCondition:
		currScript = "Scripted Condition"
		// get ticket ID from args
		if flag.NArg() != 1 {
			logger.Errorf("wrong number of positional arguments : one expected when using --scripted-condition flag, got %d: %v ", flag.NArg(), args)
		}
		endpoint = fmt.Sprintf("execute-scripted-condition/%s", flag.Args()[0]) // arg can be a ticket ID or the "test" keyword

	case preAssignment:
		currScript = "Pre-Assignment"
		// get ticket ID from args
		if flag.NArg() != 0 && !(flag.NArg() == 1 && args[0] == "") {
			logger.Errorf("wrong number of positional arguments : zero expected when using --pre-assignment flag, got %d: %v", flag.NArg(), args)
		}

		endpoint = "execute-pre-assignment"
		var err error
		if reqBody, err = getInputSource(datafilenameFlag); err != nil {
			logger.Error(err)

		}

	case scriptedTask:
		currScript = "Scripted Task"
		// get ticket ID from args
		if flag.NArg() != 1 {
			logger.Errorf("wrong number of positional arguments : one expected when using --scripted-task flag, got %d: %v ", flag.NArg(), args)
		}
		endpoint = fmt.Sprintf("execute-scripted-task/%s", flag.Args()[0]) // arg can be a ticket ID or the "test" keyword
		var err error
		if reqBody, err = getInputSource(datafilenameFlag); err != nil {
			logger.Error(err)

		}
	default:
		logger.Errorf("runSpecial has been called when no special script has been provided")
	}

	logger.Infof("mediator-client is sending request to backend end-point: %s", endpoint)

	if r, err := client.NewPOSTwithToken(endpoint, reqBody, "json"); err != nil {
		logger.Error(err)
	} else if err := r.Run(&res); err != nil {

		if res.Error != "" {
			logger.Warningf("mediator-client received an error from backend when trying to run %s script: %s", currScript, err)
			logger.Errorf("detailed error: %v", res.Error)
		} else {
			logger.Errorf("mediator-client received an error from backend when trying to run %s script: %s", currScript, err)
		}

	} /*else if res.Error != "" {
		// at this point, there may still be an error if script returned an error
		logger.Warningf("%s script failed with exit code %d and error %s", currScript, res.ExitCode, res.Error)
		logger.Warningf("%s script stderr is: %s", currScript, res.StdErr)
		fmt.Print(res.StdErr)
		os.Exit(res.ExitCode)

	} else {
		logger.Infof("%s script succeded and returned '%s'", currScript, res.StdOut)
		fmt.Print(res.StdOut)
		os.Exit(0)
	}*/

	// log whatever we got from back
	logger.Infof("%s script returned:", currScript)
	logger.Infof(" - script output: %s", res.StdOut)
	logger.Infof(" - script error: %s", res.StdErr)
	logger.Infof(" - execution error: %s", res.Error)
	logger.Infof(" - exit code: %d", res.ExitCode)

	// send what we got back to SecureChange
	os.Stdout.WriteString(res.StdOut)
	os.Stderr.WriteString(res.StdErr)
	os.Exit(res.ExitCode)
}
