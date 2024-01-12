package securechangeapi

import (
	"fmt"
	"strings"
	"uqtu/mediator/console"
	"uqtu/mediator/scworkflow"

	"github.com/spf13/cobra"
)

var (
	SCworkflows                      *scworkflow.Workflows
	create_for_all                   bool
	wf_name_flag                     string
	list_triggers                    []string
	New_wf_triggers                  scworkflow.WorkflowTriggers
	MediatorSecurechangeAPICreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new SecurechangeAPI trigger configuration",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if SCworkflows == nil {
				fmt.Println("Get fresh list of workflows from SC...")
				// send request to SC
				SCworkflows, err = Manager.GetSecurechangeWorkflows(false)
				if err != nil {
					return err
				}
				fmt.Println("OK !")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				wf_name string = wf_name_flag
				wf_id   int
			)

			// a wf was provided. Is it in list?
			if wf_name != "" {
				// make sure wf is in list
				found := false
				for _, c := range SCworkflows.Workflows {
					if c.Name == wf_name {
						found = true
						break
					}
				}
				if !found {
					fmt.Printf("[ERROR] No workflow named '%s' was found.\n", wf_name)
					wf_name = ""
				}
			}

			// select a workflow if needed
			if wf_name == "" {
				items := []string{}
				for _, c := range SCworkflows.Workflows {
					items = append(items, c.Name)
				}
				for {
					if wf_name, err = console.SelectFromList("Select the workflow for which you want to create a trigger :", items, nil); err == nil && wf_name != "" {
						for _, c := range SCworkflows.Workflows {
							if c.Name == wf_name {
								wf_id = c.Id
							}
						}
						break
					}
				}
			}

			for i := 1; i < int(scworkflow.LAST_TRIGGER); i++ {
				t := scworkflow.SecurechangeTrigger(i)
				list_triggers = append(list_triggers, t.String())
			}

			if !create_for_all {
				var tg string
				// select a trigger
				for {
					if tg, err = console.SelectFromList("Select a trigger :", list_triggers, nil); err == nil && tg != "" {
						break
					}
				}
				list_triggers = []string{tg}
			}

			New_wf_triggers.WorkflowTriggers.WorkflowTrigger = make([]scworkflow.WorkflowTrigger, len(list_triggers))
			for i, tg_name := range list_triggers {
				// create data set
				tg_slug := scworkflow.GetTriggerFromString(tg_name).Slug()
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Name = fmt.Sprintf("%s %s", wf_name, tg_slug)
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Executer.Type = "ScriptDTO"
				// swap Arguments and Path fields to fix a silly Securechange bug
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Executer.Arguments = fmt.Sprintf("/opt/tufin/data/securechange/scripts/mediator-client-%s.sh", strings.ToLower(tg_name))
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Executer.Path = wf_name

				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Triggers = make([]scworkflow.WorkflowTriggerGroup, 1)
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Triggers[0].Name = fmt.Sprintf("trigger %s", tg_name)
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Triggers[0].Workflow.Name = wf_name
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Triggers[0].Workflow.ParentWorkflowID = wf_id
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Triggers[0].Events = make([]string, 1)
				New_wf_triggers.WorkflowTriggers.WorkflowTrigger[i].Triggers[0].Events[0] = tg_slug
			}

			return nil
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			if err := Manager.CreateSecurechangeWorkflowTriggers(&New_wf_triggers); err != nil {
				return err
			} else {
				fmt.Printf("SecurechangeAPI trigger was created for trigger(s) %v.\n", list_triggers)
			}
			return nil
		},
	}
)

func init() {
	MediatorSecurechangeAPICreateCmd.Flags().BoolVar(&create_for_all, "all-triggers", false, "Create SecureChangeAPI trigger conf for all triggers.")
	MediatorSecurechangeAPICreateCmd.Flags().StringVarP(&wf_name_flag, "workflow", "w", "", "Create SecureChangeAPI trigger conf for the provided workflow.")
}
