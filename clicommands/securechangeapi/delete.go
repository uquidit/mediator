package securechangeapi

import (
	"fmt"
	"sort"
	"strconv"
	"uqtu/mediator/console"
	"uqtu/mediator/scworkflow"

	"github.com/spf13/cobra"
)

var (
	force_delete                     bool
	MediatorSecurechangeAPIDeleteCmd = &cobra.Command{
		Use:     "delete [id]",
		Aliases: []string{"del", "rm"},
		Short:   "Delete SecurechangeAPI configuration",
		Long: `Delete SecurechangeAPI configuration.

If a valid SecurechangeAPI trigger ID is provided, it will be deleted. --all-triggers and --workflow flags are ignored.
Otherwise, you will be prompted to select the trigger you want to delete.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				id_is_provided bool = false
				confirmation   bool
			)
			triggers_to_delete := []*scworkflow.WorkflowTrigger{}

			if len(args) > 0 && args[0] != "" {
				// A trigger ID has been provided as argument
				//
				// convert arg into a trigger ID. Complains if not possible.
				// then get the trigger and add it to to_delete list.
				// complain if trigger is not found.

				if id, err := strconv.Atoi(args[0]); err != nil {
					return err

				} else if trigger_to_delete, err := Manager.GetSecurechangeWorkflowTriggerByID(id); err != nil {
					return err
				} else {
					id_is_provided = true
					triggers_to_delete = append(triggers_to_delete, trigger_to_delete)
				}
			}

			if !id_is_provided {

				// get a list of WF from --workflow flag
				// if list is empty, ask user to select a WF
				wf_list := getWorkflowsFromFlag(cmd)

				if len(wf_list) == 0 { // ask user to select a WF
					if w, err := selectWorkflow("Select the workflow for which you want to delete a trigger :"); err != nil {
						return err
					} else if w == nil {
						// exit
						return nil
					} else {
						wf_list = append(wf_list, w)
					}
				}
				// send request to SC to get list of triggers
				if conf, err := Manager.GetSecurechangeWorkflowTriggers(); err != nil {
					return err
				} else {
					items := []console.Item{}
					for _, t := range conf.WorkflowTriggers.WorkflowTrigger {
						// check if trigger is related to any of the wf in list
						if !isTriggerRelatedToWorkflowInList(wf_list, t) {
							continue
						}

						if all_triggers {
							// we will delete all triggers. No need to ask user to choose
							// just add trigger to delete list
							triggers_to_delete = append(triggers_to_delete, t)
						} else {
							item := my_trigger{
								id:        t.ID,
								name:      t.Name,
								path:      t.Executer.Path,
								arguments: t.Executer.Arguments,
							}
							items = append(items, item)
						}
					}

					// if there is no trigger in list, there is nothing to do
					if len(items) == 0 && len(triggers_to_delete) == 0 {
						fmt.Println("No triggers found for selected workflow(s)")
						return nil
					}

					if !all_triggers {
						// sort list of triggers
						sort.SliceStable(items, func(p, q int) bool {
							return items[p].GetLabel() < items[q].GetLabel()
						})

						for {
							if id, err := console.SelectFromItemList("Select the trigger you want to delete:", items, nil); err == nil {
								for _, t := range conf.WorkflowTriggers.WorkflowTrigger {
									if id == t.ID {
										triggers_to_delete = append(triggers_to_delete, t)
									}
								}
								break
							}
						}
					}
				}
			}

			if len(triggers_to_delete) == 0 {
				fmt.Println("Nothing to do.")
				return nil
			}

			if !force_delete {
				msg := ""
				if len(triggers_to_delete) == 1 {
					msg = fmt.Sprintf("Are you sure you want to permanently delete Securechange API trigger '%s'", triggers_to_delete[0].Name)
				} else {
					fmt.Printf("The following %d triggers have been selected for deletion:\n", len(triggers_to_delete))
					for _, t := range triggers_to_delete {
						fmt.Printf(" - %s (#%d)\n", t.Name, t.ID)
					}
					msg = "Are you sure you want to permanently delete these Securechange API triggers"
				}
				var err error
				for {
					if confirmation, err = console.GetBoolean(msg, console.GetBooleanDefault_No); err == nil {
						break
					}
				}
			} else {
				confirmation = true
			}

			if confirmation {
				var err error
				for _, t := range triggers_to_delete {
					if err = Manager.DeleteSecurechangeWorkflowTriggers(t.ID); err != nil {
						fmt.Printf("ERROR: SecurechangeAPI trigger '%s' (id=%d) could not be deleted: %v\n", t.Name, t.ID, err)
					} else {
						fmt.Printf("SecurechangeAPI trigger '%s' (id=%d) was deleted.\n", t.Name, t.ID)
					}
				}
			} else {
				fmt.Println("Cancel.")
			}
			return nil
		},
	}
)

func init() {
	MediatorSecurechangeAPIDeleteCmd.Flags().BoolVar(&force_delete, "force", false, "Don't ask for confirmation.")
}
