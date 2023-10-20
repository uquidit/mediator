package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/mediatorscript"
	"uqtu/mediator/totp"

	"github.com/sirupsen/logrus"
)

func runInteractiveScripts(scriptedCondition bool, preAssignment bool, scriptedTask bool, datafilenameFlag string, conf *mediatorscript.MediatorConfiguration) {
	if (scriptedCondition && (preAssignment || scriptedTask)) ||
		(preAssignment && (scriptedCondition || scriptedTask)) ||
		(scriptedTask && (scriptedCondition || preAssignment)) {
		logrus.Fatalf("Only one of the --scripted-condition --pre-assignment and --scripted-task flag can be used at a time.")
	}
	var (
		reqBody    io.Reader
		endpoint   string
		currScript string
		err        error
	)

	client := apiclient.NewClient(conf.Configuration.BackendURL, "", "", conf.Configuration.SSLSkipVerify)
	client.Token, err = totp.GetKey()
	if err != nil {
		logrus.Fatal(err)
	}

	args := flag.Args()
	switch {
	case scriptedCondition:
		currScript = "Scripted Condition"
		// get ticket ID from args
		if flag.NArg() != 1 {
			logrus.Fatalf("wrong number of positional arguments : one expected when using --scripted-condition flag, got %d: %v ", flag.NArg(), args)
		}
		endpoint = fmt.Sprintf("execute-scripted-condition/%s", flag.Args()[0]) // arg can be a ticket ID or the "test" keyword

	case preAssignment:
		currScript = "Pre-Assignment"
		// get ticket ID from args
		if flag.NArg() != 0 && !(flag.NArg() == 1 && args[0] == "") {
			logrus.Fatalf("wrong number of positional arguments : zero expected when using --pre-assignment flag, got %d: %v", flag.NArg(), args)
		}

		endpoint = "execute-pre-assignment"
		var err error
		if reqBody, err = getInputSource(datafilenameFlag); err != nil {
			logrus.Fatal(err)

		}

	case scriptedTask:
		currScript = "Scripted Task"
		// get ticket ID from args
		if flag.NArg() != 1 {
			logrus.Fatalf("wrong number of positional arguments : one expected when using --scripted-task flag, got %d: %v ", flag.NArg(), args)
		}
		endpoint = fmt.Sprintf("execute-scripted-task/%s", flag.Args()[0]) // arg can be a ticket ID or the "test" keyword
		var err error
		if reqBody, err = getInputSource(datafilenameFlag); err != nil {
			logrus.Fatal(err)

		}
	default:
		logrus.Fatalf("runSpecial has been called when no special script has been provided")
	}

	logrus.Infof("mediator-client is sending request to backend end-point: %s", endpoint)
	res := mediatorscript.RunResponse{}

	if r, err := client.NewPOSTwithToken(endpoint, reqBody, "json"); err != nil {
		logrus.Fatal(err)
	} else if err := r.Run(&res); err != nil {

		if res.Error != "" {
			logrus.Warningf("mediator-client received an error from backend when trying to run %s script: %v", currScript, err)
			logrus.Fatalf("detailed error: %v", res.Error)
		} else {
			logrus.Fatalf("mediator-client received an error from backend when trying to run %s script: %v", currScript, err)
		}
	} else if res.Error != "" {
		// just in case we get an error in response
		logrus.Fatalf("mediator-client received an error from backend when trying to run %s script: %v", currScript, res.Error)

	} else if len(res.RunResults) > 0 {
		// at this point, there may still be an error if script returned an error
		for _, r := range res.RunResults {
			// there should only be one item really...

			// any internal error at script level?
			if r.InternalError != "" {
				logrus.Fatalf("an internal error occured in backend when trying to run %s script: %v", currScript, r.InternalError)
			}

			// log whatever we got from back
			logrus.Infof("%s script returned:", currScript)
			logrus.Infof(" - script output: %s", r.StdOut)
			logrus.Infof(" - script error: %s", r.StdErr)
			logrus.Infof(" - execution error: %s", r.ScriptError)
			logrus.Infof(" - exit code: %d", r.ExitCode)

			// send what we got back to SecureChange
			os.Stdout.WriteString(r.StdOut)
			os.Stderr.WriteString(r.StdErr)
			os.Exit(r.ExitCode)

		}

	} else {
		logrus.Fatalf("mediator-client received no result from script execution")
	}

}
