package securechangeapi

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var (
	Manager                        SecurchangeAPIManager
	all_triggers                   bool
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

		RunE: func(cmd *cobra.Command, args []string) error {

			SCapi_wf_triggers, err := Manager.GetSecurechangeWorkflowTriggers()
			if err != nil {
				return err
			}
			trigger_list := SCapi_wf_triggers.WorkflowTriggers.WorkflowTrigger
			if len(trigger_list) == 0 {
				fmt.Println("No trigger found in Securechange.")
				return nil
			}
			// sort list of triggers by WFname
			sort.SliceStable(trigger_list, func(p, q int) bool {
				return trigger_list[p].Name < trigger_list[q].Name
			})
			// get a list of WF from --workflow flag
			// if list is empty, ask user to select a WF
			wf_list := getWorkflowsFromFlag(cmd)

			for _, c := range trigger_list {
				// check if trigger is related to any of the wf in list
				if len(wf_list) > 0 && !isTriggerRelatedToWorkflowInList(wf_list, c) {
					continue
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
	MediatorSecurechangeAPICmd.PersistentFlags().StringSliceP("workflow", "w", nil, "Comma separated list of workflow names. Only apply command to SecurchangeAPI configuration that includes a workflow in the provided list. Can also be used multiple times.")
	MediatorSecurechangeAPICmd.PersistentFlags().BoolVarP(&all_triggers, "all-triggers", "a", false, "Apply command to all SecureChangeAPI triggers. When used in conjunction with --workflow flag, the command applies to all triggers of provided workflows.")
	// MediatorSecurechangeAPIShowCmd.Flags().StringSliceVarP(&WF_to_process, "workflow", "w", nil, "Comma separated list of workflow names. Only show SecurchangeAPI configuration that includes a workflow in the provided list. Can also be used multiple times.")

	MediatorSecurechangeAPICmd.PersistentFlags().StringVarP(&localManager.SC_host, "host", "H", "", "SecureChange host. Will be prompted if not provided.")
	MediatorSecurechangeAPICmd.PersistentFlags().StringVarP(&localManager.SC_username, "username", "U", "", "SecureChange user name. Will be prompted if not provided.")
	MediatorSecurechangeAPICmd.PersistentFlags().StringVarP(&localManager.SC_pwd, "password", "P", "", "SecureChange password. Will be prompted if not provided.")

	MediatorSecurechangeAPICmd.AddCommand(MediatorSecurechangeAPIDeleteCmd)
	MediatorSecurechangeAPICmd.AddCommand(MediatorSecurechangeAPICreateCmd)
	MediatorSecurechangeAPICmd.AddCommand(MediatorSecurechangeAPIShowCmd)

}
