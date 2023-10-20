package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/console"
	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type WorkflowXML struct {
	Id    string `xml:"id"`
	Name  string `xml:"name"`
	Steps []struct {
		Name string `xml:"name"`
	} `xml:"steps>step"`
}
type Workflows struct {
	XMLName   xml.Name      `xml:"workflows"`
	Workflows []WorkflowXML `xml:"workflow"`
}

// generateConfigCmd represents the generateConfig command
var (
	username, pwd, script, conf_filename string
	generateConfigCmd                    = &cobra.Command{
		Use:   "generate-config <SecureChange API URL>",
		Short: "Generate a template configuration file for mediator-client",
		Long: `Generate a template configuration file for mediator-client.
It will be populated with workflows and step configured in SecureChange.
If provided, a default script will be associated to all the steps.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scURL := args[0]
			var (
				workflowsSecureChange Workflows
				err                   error
			)

			// ask for user name if not provided via dedicated flag
			if username == "" {
				if username, err = console.GetText("SecureChange Username"); err != nil {
					return err
				}
			}

			// ask for password if not provided via dedicated flag
			if pwd == "" {
				if pwd, err = console.GetPassword("SecureChange Password"); err != nil {
					return err
				} else {
					fmt.Println("---")
				}
			}

			// get data from securechange
			if client := apiclient.NewClient(scURL, username, pwd, true); client == nil {
				return fmt.Errorf("cannot get API client")

			} else if req, err := client.NewGETwithBasicAuth("/workflows/active_workflows", "xml"); err != nil {
				return err

			} else if err := req.Run(&workflowsSecureChange); err != nil {
				return err

			} else {
				// initialize conf object with collected data
				var conf mediatorscript.MediatorConfiguration
				conf.Configuration.BackendURL = URL
				conf.Configuration.Logfile = "/var/log/mediator-client.log"
				conf.Configuration.SSLSkipVerify = false
				conf.Workflows = make([]mediatorscript.Workflow, len((workflowsSecureChange.Workflows)))

				for i, element := range workflowsSecureChange.Workflows {
					// get detailed data for each workflow
					url := fmt.Sprintf("/workflows?id=%s", element.Id)
					var workflow WorkflowXML
					if req, err := client.NewGETwithBasicAuth(url, "xml"); err != nil {
						return err

					} else if err := req.Run(&workflow); err != nil {
						return err

					} else {
						// fill current workflow in
						fmt.Printf("Found workflow '%s'. It has %d steps:\n", element.Name, len(workflow.Steps))
						conf.Workflows[i].Name = element.Name
						conf.Workflows[i].Steps = make([]mediatorscript.Steps, len(workflow.Steps))
						// add steps
						for j, s := range workflow.Steps {
							fmt.Printf("  - %s\n", s.Name)
							conf.Workflows[i].Steps[j].Name = s.Name
							// use default value for script. Can be "" if not provided by flag
							conf.Workflows[i].Steps[j].Script = script
						}
					}
				}

				// dump data in YAML format
				if data, err := yaml.Marshal(conf); err != nil {
					return err
				} else if err := os.WriteFile(conf_filename, data, 0644); err != nil {
					return err
				} else {
					fmt.Printf("Configuration file %s was successfully created.\n", conf_filename)
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(generateConfigCmd)

	generateConfigCmd.Flags().StringVar(&username, "username", "", "SecureChange user name. Will be prompted if not provided.")
	generateConfigCmd.Flags().StringVar(&pwd, "password", "", "SecureChange password. Will be prompted if not provided.")
	generateConfigCmd.Flags().StringVar(&script, "script", "", "Default script name. Used for all steps of all workflows")
	generateConfigCmd.Flags().StringVar(&conf_filename, "output", "mediator-client.yml", "Name or path of file where generated configuration will be dumped.")
}
