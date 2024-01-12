package securechangeapi

import (
	"fmt"
	"uqtu/mediator/scworkflow"

	"github.com/spf13/cobra"
)

var (
	Manager                        SecurchangeAPIManager
	SCapi_wf_triggers              *scworkflow.WorkflowTriggers
	WF_to_process                  []string
	MediatorSecurechangeAPIShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show SecurechangeAPI configuration",
		Args:  cobra.MaximumNArgs(0),
		RunE:  MediatorSecurechangeAPICmd.RunE,
	}
	MediatorSecurechangeAPICmd = &cobra.Command{
		Use:     "securechange-api",
		Aliases: []string{"scapi", "sc", "api"},
		Short:   "Show and manage SecurechangeAPI configuration",
		Long:    "If no subcommand is provided, show SecurechangeAPI configuration.",
		Args:    cobra.MaximumNArgs(0),

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if SCapi_wf_triggers, err = Manager.GetSecurechangeWorkflowTriggers(); err != nil {
				return err
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {

			// send request to SC
			for _, c := range SCapi_wf_triggers.WorkflowTriggers.WorkflowTrigger {
				if len(WF_to_process) > 0 {
					found := false
					for _, w := range WF_to_process {
						for _, t := range c.Triggers {
							if t.Workflow.Name == w {
								found = true
								break
							}
						}
						if found {
							break
						}
					}
					if !found {
						continue
					}

				}
				fmt.Printf("\n*** #%d %s\n", c.ID, c.Name)
				fmt.Printf("  - Execute: %s \"%s\"\n", c.Executer.Path, c.Executer.Arguments)
				fmt.Println("  - Trigger groups:")
				for _, t := range c.Triggers {
					fmt.Printf("    - Name: %s\n", t.Name)
					fmt.Printf("    - Workflow: %s\n", t.Workflow.Name)
					fmt.Printf("    - Triggers: %v\n", t.Events)
				}
			}

			return nil
		},
	}
)

func init() {
	localManager := scManager{}
	Manager = &localManager
	MediatorSecurechangeAPICmd.Flags().StringSliceVarP(&WF_to_process, "workflow", "w", nil, "Comma separated list of workflow names. Only show SecurchangeAPI configuration that includes a workflow in the provided list. Can also be used multiple times.")
	MediatorSecurechangeAPIShowCmd.Flags().StringSliceVarP(&WF_to_process, "workflow", "w", nil, "Comma separated list of workflow names. Only show SecurchangeAPI configuration that includes a workflow in the provided list. Can also be used multiple times.")

	MediatorSecurechangeAPICmd.PersistentFlags().StringVarP(&localManager.SC_host, "host", "H", "", "SecureChange host. Will be prompted if not provided.")
	MediatorSecurechangeAPICmd.PersistentFlags().StringVarP(&localManager.SC_username, "username", "U", "", "SecureChange user name. Will be prompted if not provided.")
	MediatorSecurechangeAPICmd.PersistentFlags().StringVarP(&localManager.SC_pwd, "password", "P", "", "SecureChange password. Will be prompted if not provided.")

	MediatorSecurechangeAPICmd.AddCommand(MediatorSecurechangeAPIDeleteCmd)
	MediatorSecurechangeAPICmd.AddCommand(MediatorSecurechangeAPICreateCmd)
	MediatorSecurechangeAPICmd.AddCommand(MediatorSecurechangeAPIShowCmd)

}
