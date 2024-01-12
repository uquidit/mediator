package clicommands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/console"
	"uqtu/mediator/mediatorsettings"
	"uqtu/mediator/scworkflow"

	"github.com/spf13/cobra"
)

var (
	Settings_endpoint   string = "/settings"
	SettingsFile        string
	WFsettings          mediatorsettings.MediatorSettingsMap
	SCworkflows         *scworkflow.Workflows
	workflows_steps     scworkflow.WorkflowsStepsList
	SC_username         string
	SC_pwd              string
	SC_host             string
	SC_querystring      string
	trigger_script_list []string
	save_on_exit        bool = true
	MediatorSettingsCmd      = &cobra.Command{
		Use:   "settings",
		Short: "Update Mediator settings",
		Long: `Interactively update Mediator client settings.
	
These settings tells Mediator client which script should be run for a given ticket and trigger.
This is an interactive command: required information will be prompted to you.`,
		Args: cobra.MaximumNArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if SCworkflows == nil {
				fmt.Println("Get fresh list of workflows from SC...")

				// ask for user name and password if not provided via dedicated flag
				if err = setSCcredentials(); err != nil {
					return err
				}

				// send request to SC !!
				// TODO: send request to backend instead
				SCworkflows, err = scworkflow.GetSecurechangeWorkflows(SC_username, SC_pwd, SC_host, true)
				if err != nil {
					return err
				}
				fmt.Println("OK !")
			}
			workflows_steps = SCworkflows.GetWorkflowsSteps()

			// get settings from server of a file if provided
			if SettingsFile == "" {
				fmt.Print("Getting settings from server...")
				wf_settings_slice := mediatorsettings.MediatorSettings{}

				// build a querystring with SC credentials
				if err = setSCcredentials(); err != nil {
					return err
				}

				endpoint := fmt.Sprintf("%s?%s", Settings_endpoint, SC_querystring)

				if _, err := apiclient.RunGETwithToken(endpoint, "json", &wf_settings_slice); err != nil {
					return err
				} else {
					WFsettings = wf_settings_slice.GetMap()
				}
				fmt.Println("    OK !")

			} else {
				fmt.Printf("Reading settings file %s...", SettingsFile)
				if WFsettings, err = mediatorsettings.ReadWorkflowsSettings(SettingsFile); err != nil {
					fmt.Println("")
					return fmt.Errorf("error while reading configuration file %s: %w", SettingsFile, err)
				}

				if errs := WFsettings.SetNextStep(workflows_steps); len(errs) > 0 {
					return fmt.Errorf("cannot fix next-step: %v", errs)
				}
				fmt.Println("    OK !")
			}

			// Get list of registered trigger scripts from backend
			if trigger_script_list, err = GetTriggerScripts(); err != nil {
				return err
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {

			if SCworkflows == nil {
				return fmt.Errorf("init error: Securechange workflow list is empty. Check at least one workflow is activated in Securechange")
			}

			// flag workflows that has no settings
			for key, scwf := range SCworkflows.Workflows {
				_, SCworkflows.Workflows[key].HasSettings = WFsettings[scwf.Name]
			}

			//********* MAIN LOOP *************
			// walk the workflow list given by Securechange
			// For every wf, ask use if it requires a settings
			// if so, user will choose from script list for every trigger
			for {
				var scwf *scworkflow.WorkflowXML
				fmt.Println("\n******************************************")
				scwf, save_on_exit = selectWorkflow(SCworkflows)
				if scwf == nil {
					fmt.Println("Exiting...")
					break
				}

				var (
					found bool
				)
				if _, found = WFsettings[scwf.Name]; !found {
					WFsettings[scwf.Name] = &mediatorsettings.WFSettings{
						WFname: scwf.Name,
						WFid:   scwf.Id,
						Rules:  mediatorsettings.RulesSlice{},
					}

				}

				// Any settings for this WF ?
				for {
					index, new, exit := selectRule(WFsettings[scwf.Name].Rules)
					if exit {
						break
					}

					if new {
						rule := getNewRule(scwf.GetSteps())
						WFsettings[scwf.Name].Rules = append(WFsettings[scwf.Name].Rules, rule)
					} else {
						rule := WFsettings[scwf.Name].Rules[index]
						delete := editRule(rule, scwf.GetSteps())
						if delete {
							WFsettings[scwf.Name].Rules[index] = nil
						}
					}

					// remove empty settings
					WFsettings[scwf.Name].Clean()
				}

				// remove WF settings from map if no setting
				WFsettings.Clean()
			}

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if !save_on_exit {
				return nil
			}
			if SettingsFile == "" {
				fmt.Print("Sending settings to backend for upload to Securechange...")
				// ask for user name and password if not provided via dedicated flag
				if err := setSCcredentials(); err != nil {
					return err
				}

				endpoint := fmt.Sprintf("%s?%s", Settings_endpoint, SC_querystring)

				if jsoninput, err := json.Marshal(WFsettings.GetSlice()); err != nil {
					return err
				} else if _, err := apiclient.RunPOSTwithToken(endpoint, bytes.NewBuffer(jsoninput), "json", nil); err != nil {
					return err
				}
				fmt.Println("    OK !")

			} else {
				// a setting file was provided
				// save updated data in that file
				// do not try to upload settings

				fmt.Printf("Writing new settings to file %s...", SettingsFile)
				settings_slice := WFsettings.GetSlice()

				if errs := settings_slice.SetPreviousStep(workflows_steps); len(errs) > 0 {
					return fmt.Errorf("cannot fix next-step: %v", errs)
				}

				if err := mediatorsettings.WriteWorkflowsSettings(settings_slice, SettingsFile); err != nil {
					fmt.Println("")
					return fmt.Errorf("error while writing setting file %s: %w", SettingsFile, err)
				}
				fmt.Println("    OK !")
			}
			return nil
		},
	}
)

func init() {
	MediatorSettingsCmd.Flags().StringVarP(&SettingsFile, "settings", "s", "", "Path to local Mediator client settings file.")
	MediatorSettingsCmd.Flags().StringVarP(&SC_host, "host", "H", "", "SecureChange host. Will be prompted if not provided.")
	MediatorSettingsCmd.Flags().StringVarP(&SC_username, "username", "U", "", "SecureChange user name. Will be prompted if not provided.")
	MediatorSettingsCmd.Flags().StringVarP(&SC_pwd, "password", "P", "", "SecureChange password. Will be prompted if not provided.")
}

// build a querystring with SC credentials
// Ask user to provide missing credentials
// do it only once: return SC_querystring if already set.
// this function has side effects: modifies global SC_* variables
func setSCcredentials() error {
	if SC_querystring != "" {
		return nil
	}
	var err error
	// ask for user name if not provided via dedicated flag
	for {
		if SC_username != "" {
			break
		}
		if SC_username, err = console.GetText("SecureChange Username"); err != nil {
			return err
		}
	}

	// ask for password if not provided via dedicated flag
	for {
		if SC_pwd != "" {
			break
		}
		if SC_pwd, err = console.GetPassword("SecureChange Password"); err != nil {
			return err
		} else {
			fmt.Println("---")
		}
	}

	// ask for host if not provided via dedicated flag
	for {
		if SC_host != "" {
			break
		}
		if SC_host, err = console.GetPassword("SecureChange Host"); err != nil {
			return err
		} else {
			fmt.Println("---")
		}
	}

	v := url.Values{}
	v.Add("sc_username", SC_username)
	v.Add("sc_password", SC_pwd)
	v.Add("sc_host", SC_host)

	SC_querystring = v.Encode()
	return nil
}
