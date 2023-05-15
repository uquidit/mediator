package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"uqtu/mediator/mediatorscript"
	"uqtu/mediator/totp"

	"uqtu/mediator/configparser"

	"uqtu/mediator/apiclient"
	"uqtu/mediator/logger"
)

var (
	Version = "develop"
)

func main() {
	var (
		datafilenameFlag  string
		useNextStep       bool
		scriptedCondition bool
		preAssignment     bool
		scriptedTask      bool
	)
	flag.StringVar(&datafilenameFlag, "file", "", "Read data from file instead of stdin.")
	flag.BoolVar(&useNextStep, "next-step", false, "Tell mediator-client to act as if the ticket is not in provided step but in following step in list. In that case, the script associated with the next step will be executed.")
	flag.BoolVar(&scriptedCondition, "scripted-condition", false, "Tell mediator-client to request back-end to run special 'Scripted Condition' script.")
	flag.BoolVar(&preAssignment, "pre-assignment", false, "Tell mediator-client to request back-end to run special 'Pre-Assignment' script.")
	flag.BoolVar(&scriptedTask, "scripted-task", false, "Tell mediator-client to request back-end to run special 'Scripted Task' script.")
	flag.Parse()
	if flag.NArg() != 1 && !(scriptedCondition || preAssignment || scriptedTask) {
		log.Fatal("ERROR: wrong number of positional arguments: one expected, got ", flag.NArg())
	}

	var conf mediatorscript.MediatorConfiguration

	// get current executable so we can figure out what is current path
	// conf file is in the same folder as configuration file
	if ex, err := os.Executable(); err != nil {
		log.Fatal(err)
	} else {
		// get current folder from executable absolute path
		currPath := filepath.Dir(ex)

		// forge config name
		absPath := fmt.Sprintf("%s/mediator-client.yml", currPath)
		if err := configparser.ReadConfAbsolutePath(absPath, &conf, nil); err != nil {
			log.Fatal(err)
		}

		if err := logger.InitFullPath(conf.Configuration.Logfile, true, true, true); err != nil {
			log.Fatal(err)
		}
		defer logger.CloseAllLogfiles()

		if err := InitFromJSONIfNeeded(&conf, currPath); err != nil {
			logger.Error(err)
		}
		// initialize mediatorscript logger with rootlogger so all logs will go to the same location
		if rootlogger, err := logger.GetRootFileStdoutLogger(); err != nil {
			log.Fatalf("error while getting root logger: %v", err)
		} else {
			if err := mediatorscript.Init(
				"/dev/null", // no storage file. only needed in server
				rootlogger,
			); err != nil {
				logger.Warningf("error while loading scripts for mediator list: %v", err)
			}
		}

	}

	if len(conf.Configuration.BackendURL) == 0 {
		logger.Errorf("Back-end URL is empty. Stop.")
	}

	if conf.Configuration.BackendURL[:1] != "/" {
		conf.Configuration.BackendURL += "/"
	}
	logger.Infof("mediator-client will send requests to backend @ %s", conf.Configuration.BackendURL)

	if scriptedCondition || preAssignment || scriptedTask {
		logger.Infof("Starting mediator-client for Interactive scripts")

		// this function will terminate current process
		// it will deal with all error cases and log accordingly
		// actually, the code is in a separate function in an effort to keep main() short
		runInteractiveScripts(scriptedCondition, preAssignment, scriptedTask, datafilenameFlag, &conf)

	}

	if useNextStep {
		logger.Infof("Starting mediator-client in 'next step' mode")
	} else {
		logger.Infof("Starting mediator-client in normal mode")
	}

	if source, err := getInputSource(datafilenameFlag); err != nil {
		logger.Error(err)

	} else if xmlData, err := io.ReadAll(source); err != nil {
		logger.Error(err)

	} else {
		// check for test run
		// we need to check the system is working fine when user clicks on "Test" button in SC
		// in that case, input is "<ticket_info/>"
		if bytes.Contains(xmlData, []byte("<ticket_info/>")) {
			logger.Infof("mediator-client is in test mode")

			// forward test request to backend:
			client := apiclient.NewClient(conf.Configuration.BackendURL, "", "", conf.Configuration.SSLSkipVerify)
			client.Token = totp.GetKey()

			res := map[string]string{}
			logger.Infof("mediator-client is sending test resquest")
			if r, err := client.NewPOSTwithToken("test-all", nil, "json"); err != nil {
				logger.Error(err)
				os.Exit(2)
			} else if err := r.Run(&res); err != nil {
				logger.Infof("mediator-client received an error from backend: %v", err)
				os.Exit(2)
			} else {

				logger.Infof("%d tests performed. Results are:", len(res))
				nberrors := 0
				for name, result := range res {
					if result == "OK" {
						logger.Infof(" - %s: OK", name)
					} else {
						logger.Infof(" - %s: %s", name, result)
						nberrors += 1
					}
				}
				if nberrors == 0 {
					logger.Infof("All tests passed")
					os.Exit(0)
				} else {
					logger.Warningf("%d/%d test(s) failed", nberrors, len(res))
					os.Exit(1)
				}
			}

		}
		// real mode

		// get workflow
		current_workflow := flag.Args()[0]

		workflow, err := conf.GetWorkflow(current_workflow)
		if err != nil {
			logger.Error(err)
		} else {
			logger.Infof("mediator-client called on worflow '%s'", current_workflow)

		}

		// get ticket info
		var data mediatorscript.TicketInfo
		if err := xml.Unmarshal(xmlData, &data); err != nil {
			logger.Errorf("cannot parse XML input #%s#: %v", string(xmlData), err)
		} else {
			logger.Infof("mediator-client received XML data: %s", string(xmlData))

			// work out what is current step
			// we look for step name in completion_data and in current_stage.
			// they can't be both null neither both non null
			// if completion data is provided, get from there. ignore next step
			// otherwise, get from current status. In that case, check next-step flag as we may need to get next step.
			currentStep := ""
			switch {
			case data.CompletionData != nil && data.CurrentStage == nil:
				currentStep = data.CompletionData.Name
			case data.CompletionData == nil && data.CurrentStage != nil && !useNextStep:
				currentStep = data.CurrentStage.Name
			case data.CompletionData == nil && data.CurrentStage != nil && useNextStep:
				var err error
				if currentStep, err = workflow.GetNextStep(data.CurrentStage.Name); err != nil {
					logger.Error(err)
				}
			default:
				logger.Errorf("unexpected XML data: both 'current_stage' and 'completion_data' tags are missing. Cannot get ticket step. Stop.")
			}

			if script := workflow.GetScriptForStep(currentStep); script != "" {
				logger.Infof("mediator-client found a script for ticket '%s' (ID=%d) in step %s: %s. Trigger action.", data.Subject, data.ID, currentStep, script)

				client := apiclient.NewClient(conf.Configuration.BackendURL, "", "", conf.Configuration.SSLSkipVerify)
				client.Token = totp.GetKey()

				script_url := fmt.Sprintf("execute/%s", script)
				logger.Infof("mediator-client is sending resquest to entry point: %s", script_url)

				if jsonData, err := json.Marshal(data); err != nil {
					logger.Error(err)
				} else if r, err := client.NewPOSTwithToken(script_url, bytes.NewBuffer(jsonData), "json"); err != nil {
					logger.Error(err)
				} else if response, err := r.RunWithoutDecode(); err != nil {
					if response != nil {
						b, _ := io.ReadAll(response)
						logger.Infof("mediator-client received an error: %v. response body is %s", err, string(b))
					}
					logger.Error(err)
				} else {
					logger.Infof("sent json data: %s", string(jsonData))
					if response != nil {
						b, _ := io.ReadAll(response)
						if len(b) > 0 {
							logger.Infof("mediator-client received a valid response: %s", string(b))
						} else {
							logger.Infof("mediator-client received an empty OK response")
						}
					} else {
						logger.Infof("mediator-client received no response but no error!")
					}
					os.Exit(0)

				}

			}
			logger.Infof("mediator-client found no script for ticket '%s' (ID=%d) in step '%s'. Do nothing.\n", data.Subject, data.ID, currentStep)
		}
	}
}

func getInputSource(datafilenameFlag string) (*os.File, error) {
	if datafilenameFlag != "" {
		logger.Infof("mediator-client gets ticket data from file '%s'", datafilenameFlag)

		return os.Open(datafilenameFlag)
	} else {
		return os.Stdin, nil
	}
}
