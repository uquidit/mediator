package main

import (
	"fmt"
	"io"
	"os"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/mediatorscript"
	"uqtu/mediator/totp"

	"github.com/sirupsen/logrus"
)

func runInteractiveScripts(args arguments, conf *mediatorscript.MediatorLegacyConfiguration) {
	switch {
	case args.scriptedCondition:
		runScriptedConditionScript(args.positional, args.data_filename, conf)
	case args.scriptedTask:
		runScriptedTaskScript(args.positional, args.data_filename, conf)
	case args.preAssignment:
		runPreAssignmentScript(args.data_filename, conf)
	case args.riskAnalysis:
		runRiskAnalysisScript(args.data_filename, conf)
	}
}

func runScriptedConditionScript(args []string, datafilenameFlag string, conf *mediatorscript.MediatorLegacyConfiguration) {
	currScript := "Scripted Condition"
	// get ticket ID from args
	if len(args) != 1 {
		logrus.Fatalf("wrong number of positional arguments : one expected when using --scripted-condition flag, got %d: %v ", len(args), args)
	}
	endpoint := fmt.Sprintf("execute-scripted-condition/%s", args[0]) // arg can be a ticket ID or the "test" keyword

	if reqBody, err := getInputSource(datafilenameFlag); err != nil {
		logrus.Fatal(err)

	} else {
		requestInteractiveScriptExecution(endpoint, currScript, reqBody, conf)
	}
}

func runPreAssignmentScript(datafilenameFlag string, conf *mediatorscript.MediatorLegacyConfiguration) {
	currScript := "Pre-Assignment"
	if reqBody, err := getInputSource(datafilenameFlag); err != nil {
		logrus.Fatal(err)

	} else {
		requestInteractiveScriptExecution("execute-pre-assignment", currScript, reqBody, conf)
	}
}

func runScriptedTaskScript(args []string, datafilenameFlag string, conf *mediatorscript.MediatorLegacyConfiguration) {
	currScript := "Scripted Task"
	// get ticket ID from args
	if len(args) != 1 {
		logrus.Fatalf("wrong number of positional arguments : one expected when using --scripted-task flag, got %d: %v ", len(args), args)
	}
	endpoint := fmt.Sprintf("execute-scripted-task/%s", args[0]) // arg can be a ticket ID or the "test" keyword

	if reqBody, err := getInputSource(datafilenameFlag); err != nil {
		logrus.Fatal(err)

	} else {
		requestInteractiveScriptExecution(endpoint, currScript, reqBody, conf)
	}
}

func runRiskAnalysisScript(datafilenameFlag string, conf *mediatorscript.MediatorLegacyConfiguration) {
	currScript := "Risk Analysis"
	if reqBody, err := getInputSource(datafilenameFlag); err != nil {
		logrus.Fatal(err)

	} else {
		requestInteractiveScriptExecution("execute-risk-analysis", currScript, reqBody, conf)
	}
}

func requestInteractiveScriptExecution(endpoint string, currScript string, reqBody io.Reader, conf *mediatorscript.MediatorLegacyConfiguration) {
	var (
		err error
	)

	client := apiclient.NewClient(conf.Configuration.BackendURL, "", "", conf.Configuration.SSLSkipVerify)
	client.Token, err = totp.GetKey()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("mediator-client is sending request to backend end-point: %s", endpoint)
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
