package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"mediator/mediatorscript"
	"mediator/mediatorsettings"
	"mediator/scworkflow"
	"mediator/totp"

	"mediator/configparser"

	"mediator/apiclient"
	"mediator/logger"

	"github.com/sirupsen/logrus"
)

var (
	Version = "develop"
	trigger scworkflow.SecurechangeTrigger
)

func main() {
	args := arguments{}
	flag.StringVar(&args.data_filename, "file", "", "Read data from file instead of stdin.")
	flag.StringVar(&args.trigger, "trigger", "", "Specify Securechange trigger for current processing")
	flag.StringVar(&args.settings_filename, "settings-filename", "mediator-client.json", "Specify the name of workflow settings file.")
	flag.BoolVar(&args.scriptedCondition, "scripted-condition", false, "Tell mediator-client to request back-end to run special 'Scripted Condition' script.")
	flag.BoolVar(&args.preAssignment, "pre-assignment", false, "Tell mediator-client to request back-end to run special 'Pre-Assignment' script.")
	flag.BoolVar(&args.scriptedTask, "scripted-task", false, "Tell mediator-client to request back-end to run special 'Scripted Task' script.")
	flag.BoolVar(&args.riskAnalysis, "risk-analysis", false, "Tell mediator-client to request back-end to run special 'Risk Analysis' script.")

	// version
	versionPtr := flag.Bool("version", false, "Print version number and exit.")
	flag.BoolVar(versionPtr, "v", false, "Alias of --version")
	flag.Parse()
	if *versionPtr {
		fmt.Printf("Mediator version %s\n", Version)
		os.Exit(0)
	}
	args.positional = flag.Args()
	if t, err := checkArgumentsAndGetTrigger(args); err != nil {
		logrus.Fatal(err)
	} else {
		trigger = t
	}

	var conf mediatorscript.MediatorLegacyConfiguration

	// get current executable so we can figure out what is current path
	// conf file is in the same folder as configuration file
	ex, err := os.Executable()
	if err != nil {
		logrus.Fatal(err)
	}
	// get current folder from executable absolute path
	currPath := filepath.Dir(ex)

	if err := configparser.ReadConfAbsolutePath(fmt.Sprintf("%s/mediator-client.yml", currPath), &conf, nil); err != nil {
		logrus.Fatal(err)
	}

	if _, err := logger.InitAppLogger(
		true, // use default logger
		GetLogLevel(conf.Configuration.Log.Level),
		true, true, true, true, true, true, "",
		conf.Configuration.Log.File); err != nil {

		logrus.Fatal(err)
	}
	defer logger.CloseLogFile()

	if err := mediatorscript.Init(""); err != nil {
		// mediator-client does not need to care about this error
		// there is no file to read, required for server only
		if !errors.Is(err, mediatorscript.ErrInitNoFileName) {
			logrus.Warning(err)
		}
	}

	if len(conf.Configuration.BackendURL) == 0 {
		logrus.Fatalf("Back-end URL is empty. Stop.")
	}

	if conf.Configuration.BackendURL[:1] != "/" {
		conf.Configuration.BackendURL += "/"
	}
	logrus.Infof("mediator-client will send requests to backend @ %s", conf.Configuration.BackendURL)

	if conf.Configuration.SSLSkipVerify {
		logrus.Warningf("SSL verification will be skipped! Connection to backend is insecure!")
	}

	if args.isInteractiveScript() {
		logrus.Infof("Starting mediator-client for Interactive scripts")

		// this function will terminate current process
		// it will deal with all error cases and log accordingly
		// actually, the code is in a separate function in an effort to keep main() short
		runInteractiveScripts(args, &conf)

	}

	logrus.Infof("Starting mediator-client in normal mode. Trigger is %s", trigger)

	if wf_settings, err := mediatorsettings.ReadWorkflowsSettings(fmt.Sprintf("%s/%s", currPath, args.settings_filename)); err != nil {
		logrus.Fatal(err)

	} else if source, err := getInputSource(args.data_filename); err != nil {
		logrus.Fatal(err)

	} else if xmlData, err := io.ReadAll(source); err != nil {
		logrus.Fatal(err)

	} else {
		// get workflow
		current_workflow := flag.Args()[0]

		settings, err := wf_settings.GetWorkflowSettings(current_workflow)
		if err != nil {
			logrus.Fatal(err)
		} else {
			logrus.Infof("mediator-client called on worflow '%s'", current_workflow)
		}

		// check for test run
		// we need to check the system is working fine when user clicks on "Test" button in SC
		// in that case, input is "<ticket_info/>"
		if bytes.Contains(xmlData, []byte("<ticket_info/>")) {
			logrus.Infof("mediator-client is in test mode")

			// we will send a test request for all scripts attached to workflow
			scripts := settings.GetAllScripts(trigger)

			run_results := mediatorscript.SyncRunResponsesMap{}

			// forward test request to backend:
			// send a test request for every scripts
			// dump summary at the end
			client := apiclient.NewClient(conf.Configuration.BackendURL, "", "", conf.Configuration.SSLSkipVerify)
			for _, s := range scripts {
				var err error

				// get a new tOTP password for every request
				if client.Token, err = totp.GetKey(); err != nil {
					// this shoud never happens
					// if it does, useless to keep trying. Stop here.
					logrus.Fatalf("TOTP error: %v. Stop!", err)
				}

				// sending test request for script
				logrus.Infof("mediator-client is sending a test request for script '%s'", s)

				// build test endpoint for current script
				// it's a workflow script so it is a trigger script
				endpoint := fmt.Sprintf("test/trigger/%s", s)

				results := mediatorscript.RunResponse{}

				if r, err := client.NewPOSTwithToken(endpoint, nil, "json"); err != nil {
					// this shoud never happens
					// if it does, useless to keep trying. Stop here.
					logrus.Fatal(err)

				} else if err := r.Run(&results); err != nil {
					// something went wrong before script execution
					logrus.Warningf("mediator-client received an error from backend: %v", err)
					if results.Error != "" {
						logrus.Warningf("mediator-client received additional error information: %v", results.Error)
					}
					run_results[s] = &mediatorscript.SyncRunResponse{
						InternalError: fmt.Sprintf("mediator-client received an error from backend: %v", err),
					}
				} else if results.Error != "" {
					// just in case we get an error in response
					logrus.Warningf("mediator-client received an error: %v", results.Error)
					run_results[s] = &mediatorscript.SyncRunResponse{
						InternalError: fmt.Sprintf("mediator-client received an error: %v", err),
					}

				} else {
					// store it so we can dump results later
					run_results[s] = results.RunResults[s]
				}
			}

			// dump all results in logs
			logrus.Infof("%d tests performed. Results are:", len(run_results))
			nberrors := 0
			for name, result := range run_results {
				if result.InternalError != "" {
					logrus.Infof(" - %s: Internal error: %s", name, result.InternalError)
					nberrors += 1
				} else if result.ScriptError != "" {
					logrus.Infof(" - %s: Script error: %s", name, result.ScriptError)
					nberrors += 1
				} else {
					logrus.Infof(" - %s: OK", name)
				}
			}
			if nberrors == 0 {
				logrus.Infof("All tests passed")
				os.Exit(0)
			} else {
				logrus.Warningf("%d/%d test(s) failed", nberrors, len(run_results))
				os.Exit(1)
			}

		}

		// real mode

		// get ticket info
		var data mediatorscript.TicketInfo
		if err := xml.Unmarshal(xmlData, &data); err != nil {
			logrus.Fatalf("cannot parse XML input #%s#: %v", string(xmlData), err)
		}

		// work out what is current step
		// we look for step name in completion_data and in current_stage.
		// they can't be both null neither both non null
		// if completion data is provided, get from there. ignore next step
		// otherwise, get from current status. In that case, check next-step flag as we may need to get next step.
		currentStep := ""
		switch {
		case data.CompletionData != nil && data.CurrentStage == nil:
			currentStep = data.CompletionData.Name
			logrus.Infof("Info from XML: ticket is '%s' (ID=%d), current step is %s", data.Subject, data.ID, currentStep)
		case data.CompletionData == nil && data.CurrentStage != nil:
			currentStep = data.CurrentStage.Name
			logrus.Infof("Info from XML: ticket is '%s' (ID=%d), current step is %s", data.Subject, data.ID, currentStep)
		default:
			logrus.Warningf("mediator-client received unexpected XML data: %s", string(xmlData))
			logrus.Fatalf("unexpected XML data: both 'current_stage' and 'completion_data' tags are missing. Cannot get ticket step. Stop.")
		}

		if scripts := settings.GetScriptsForTriggerAndStep(trigger, currentStep); len(scripts) > 0 {
			logrus.Infof("mediator-client found %d script(s) for ticket '%s' (ID=%d) in step %s: %v. Trigger action.", len(scripts), data.Subject, data.ID, currentStep, scripts)

			if client, err := apiclient.NewClientWithOTP(conf.Configuration.BackendURL, conf.Configuration.SSLSkipVerify); err != nil {
				logrus.Fatalf("cannot get API client: %v", err)
			} else {
				var err error
				for _, script := range scripts {
					if script == "" {
						continue
					}
					if client.Token, err = totp.GetKey(); err != nil {
						logrus.Fatalf("TOTP error: %v. Stop!", err)
					}

					script_url := fmt.Sprintf("execute/%s", script)

					if jsonData, err := json.Marshal(data); err != nil {
						logrus.Fatal(err)
					} else if r, err := client.NewPOSTwithToken(script_url, bytes.NewBuffer(jsonData), "json"); err != nil {
						logrus.Warningf("mediator-client is sending resquest to entry point: %s", script_url)
						logrus.Fatal(err)
					} else if _, err := r.RunWithoutDecode(); err != nil { // returns 204 (no content) or an error
						// something went wrong before script execution
						// we don't need response body: error will be decoded into err
						logrus.Warningf("mediator-client sent resquest to entry point '%s'  with json data: %s", script_url, string(jsonData))
						logrus.Errorf("mediator-client received an error from backend: %v", err)
					} else {
						logrus.Infof("mediator-client received an empty OK response")
					}
				}
			}
		} else {
			logrus.Infof("mediator-client found no script for ticket '%s' (ID=%d) in step '%s'. Do nothing.", data.Subject, data.ID, currentStep)
		}
	}
}

func checkArgumentsAndGetTrigger(args arguments) (scworkflow.SecurechangeTrigger, error) {
	if args.isInteractiveScript() {
		if err := args.isUniqueInteractiveScriptFlag(); err != nil {
			return scworkflow.NO_TRIGGER, err
		}
		return scworkflow.NO_TRIGGER, nil
	}

	t := scworkflow.GetTriggerFromString(args.trigger)
	if t == scworkflow.NO_TRIGGER {
		return t, fmt.Errorf("unknown or empty trigger: %s", args.trigger)
	}

	nb_args := args.NPositional()
	if nb_args != 1 {
		return t, fmt.Errorf("wrong number of positional arguments: one expected, got %d", nb_args)
	} else {
		return t, nil
	}

}

func getInputSource(datafilenameFlag string) (*os.File, error) {
	if datafilenameFlag != "" {
		logrus.Infof("mediator-client gets ticket data from file '%s'", datafilenameFlag)

		return os.Open(datafilenameFlag)
	} else {
		return os.Stdin, nil
	}
}

func GetLogLevel(l string) logrus.Level {
	if strings.EqualFold(l, "Panic") {
		return logrus.PanicLevel
	}
	if strings.EqualFold(l, "Fatal") {
		return logrus.FatalLevel
	}
	if strings.EqualFold(l, "Error") {
		return logrus.ErrorLevel
	}
	if strings.EqualFold(l, "Warn") {
		return logrus.WarnLevel
	}
	if strings.EqualFold(l, "Info") {
		return logrus.InfoLevel
	}
	if strings.EqualFold(l, "Debug") {
		return logrus.DebugLevel
	}
	if strings.EqualFold(l, "Trace") {
		return logrus.TraceLevel
	}
	return logrus.WarnLevel
}
