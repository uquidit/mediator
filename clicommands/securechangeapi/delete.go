package securechangeapi

import (
	"fmt"
	"strconv"
	"uqtu/mediator/console"
	"uqtu/mediator/scworkflow"

	"github.com/spf13/cobra"
)

var (
	force_delete                     bool
	MediatorSecurechangeAPIDeleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm"},
		Short:   "Delete SecurechangeAPI configuration",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err               error
				id                int
				id_is_provided    bool = false
				confirmation      bool
				trigger_to_delete *scworkflow.WorkflowTrigger
			)
			if len(args) > 0 && args[0] != "" {
				if id, err = strconv.Atoi(args[0]); err != nil {
					return err
				}
				id_is_provided = true
			}

			if !id_is_provided {
				// send request to SC to get list of triggers
				if conf, err := Manager.GetSecurechangeWorkflowTriggers(); err != nil {
					return err
				} else {
					items := []console.Item{}
					for _, c := range conf.WorkflowTriggers.WorkflowTrigger {
						item := my_trigger{
							id:        c.ID,
							name:      c.Name,
							path:      c.Executer.Path,
							arguments: c.Executer.Arguments,
						}
						items = append(items, item)
					}
					for {
						if id, err = console.SelectFromItemList("Select the trigger you want to delete:", items, nil); err == nil {
							break
						}
					}
					for _, c := range conf.WorkflowTriggers.WorkflowTrigger {
						if id == c.ID {
							trigger_to_delete = &c
						}
					}
				}
			} else if trigger_to_delete, err = Manager.GetSecurechangeWorkflowTriggerByID(id); err != nil {
				return err
			}

			if !force_delete {
				for {
					if confirmation, err = console.GetBoolean(fmt.Sprintf("Are you sure you want to permanently delete Securechange API trigger %s", trigger_to_delete.Name), console.GetBooleanDefault_No); err == nil {
						break
					}
				}
			} else {
				confirmation = true
			}

			if confirmation {
				if err = Manager.DeleteSecurechangeWorkflowTriggers(id); err != nil {
					return err
				} else {
					fmt.Printf("SecurechangeAPI trigger id %d was deleted.\n", id)
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
