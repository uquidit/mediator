package securechangeapi

import (
	"fmt"
	"mediator/console"
	"mediator/scworkflow"
	"strings"

	"github.com/spf13/cobra"
)

var (
	MediatorSecurechangeAPICreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new SecurechangeAPI trigger configuration",
		Args:  cobra.ExactArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err               error
				list_all_triggers []string
				list_triggers     []string
				New_wf_triggers   scworkflow.WorkflowTriggers
			)
			verbosity, _ := cmd.Flags().GetCount("verbose")

			// get a list of WF from --workflow flag
			// if list is empty, ask user to select a WF
			wf_list := getWorkflowsFromFlag(cmd)

			if len(wf_list) == 0 {
				// select a workflow
				if w, err := selectWorkflow("Select the workflow for which you want to create a trigger :"); err != nil {
					return err
				} else if w == nil {
					// exit
					return nil
				} else {
					// add selected WF to list
					wf_list = append(wf_list, w)
				}
			}

			if verbosity > 0 {
				fmt.Println("Triggers will be created for following workflows:")
				for _, w := range wf_list {
					fmt.Print(" - ")
					fmt.Println(w.Name)
				}
			}

			for i := 1; i < int(scworkflow.LAST_TRIGGER); i++ {
				t := scworkflow.SecurechangeTrigger(i)
				list_all_triggers = append(list_all_triggers, t.String())
			}

			if !all_triggers {
				var tg string
				// select a trigger
				for {
					if tg, err = console.SelectFromList("Select a trigger :", list_all_triggers, nil); err == nil && tg != "" {
						break
					}
				}
				list_triggers = []string{tg}
			} else {
				list_triggers = list_all_triggers
			}

			if verbosity > 0 {
				fmt.Println("Following triggers will be created for above workflows:")
				for _, t := range list_triggers {
					fmt.Print(" - ")
					fmt.Println(t)
				}
			}

			// get list of existing triggers so we don't create duplicates
			if verbosity > 0 {
				fmt.Print("Getting list of existing triggers...")
			}
			existing_conf, err := Manager.GetSecurechangeWorkflowTriggers()
			if err != nil {
				return err
			}
			if verbosity > 0 {
				fmt.Printf("  %d existing triggers.\n", len(existing_conf.WorkflowTriggers.WorkflowTrigger))
			}
			New_wf_triggers.WorkflowTriggers.WorkflowTrigger = make([]*scworkflow.WorkflowTrigger, 0)
			for _, w := range wf_list {
				wf_name := w.Name
				for _, tg_name := range list_triggers {
					// create data set
					tg_slug := scworkflow.GetTriggerFromString(tg_name).Slug()
					new_trigger := scworkflow.WorkflowTrigger{}
					new_trigger.Name = fmt.Sprintf("%s %s", wf_name, tg_slug)
					new_trigger.Executer.Type = "ScriptDTO"
					new_trigger.Executer.Arguments = wf_name
					new_trigger.Executer.Path = fmt.Sprintf("/opt/tufin/data/securechange/scripts/mediator-client-%s.sh", strings.ReplaceAll(strings.ToLower(tg_name), " ", "-"))

					new_trigger_group := scworkflow.WorkflowTriggerGroup{}
					new_trigger_group.Name = fmt.Sprintf("trigger %s", tg_name)
					new_trigger_group.Workflow.Name = wf_name
					// new_trigger_group.Workflow.ParentWorkflowID = wf_id
					new_trigger_group.Events = make([]string, 1)
					new_trigger_group.Events[0] = tg_slug

					new_trigger.Triggers = make([]*scworkflow.WorkflowTriggerGroup, 1)
					new_trigger.Triggers[0] = &new_trigger_group

					if !new_trigger.IsTriggerAlreadyInList(existing_conf.WorkflowTriggers.WorkflowTrigger) {
						New_wf_triggers.WorkflowTriggers.WorkflowTrigger = append(New_wf_triggers.WorkflowTriggers.WorkflowTrigger, &new_trigger)
					} else {
						fmt.Printf("Similar trigger already exists in Securechange configuration: %s\n", new_trigger.Name)
					}
				}
			}
			if len(New_wf_triggers.WorkflowTriggers.WorkflowTrigger) == 0 {
				fmt.Println("Nothing to do.")
				return nil
			}
			// send creation request to SC
			if err := Manager.CreateSecurechangeWorkflowTriggers(&New_wf_triggers); err != nil {
				return err
			} else {
				fmt.Printf("%d SecurechangeAPI triggers were created.\n", len(New_wf_triggers.WorkflowTriggers.WorkflowTrigger))
			}
			return nil
		},
	}
)
