package clicommands

import (
	"fmt"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
)

// ScriptCmd represents the script command
var (
	ScriptCmd = &cobra.Command{
		Use:     "script",
		Aliases: []string{"scripts"},
		Short:   "List available scripts for mediator. Available alias:'scripts'",
		Long: `Give a list of all registered script the mediator-client can use.
	
Use their name in mediatorscript workflow configuration to enable execution.
Add new script using the "script register" command.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			if list_by_type, err := getAllScriptNamesByType(); err != nil {
				return err
			} else {

				fmt.Println("List of registered scripts:")
				for _, t := range []mediatorscript.ScriptType{
					mediatorscript.ScriptTrigger,
					mediatorscript.ScriptCondition,
					mediatorscript.ScriptTask,
					mediatorscript.ScriptAssignment,
				} {
					fmt.Printf("\n* %s: %d\n", t, len(list_by_type[t]))
					for _, s := range list_by_type[t] {
						fmt.Printf("  - %s: %s\n", s.Name, s.Fullpath)
					}
				}

				return nil
			}
		},
	}
	interactiveTypes []mediatorscript.ScriptType
)

func init() {
	grpall := cobra.Group{
		ID:    "all",
		Title: "Global commands to manage all scripts at once:",
	}
	grptype := cobra.Group{
		ID:    "type",
		Title: "Manage scripts by type:",
	}
	ScriptCmd.AddGroup(&grpall, &grptype)

	UnregisterAllCmd.GroupID = "all"
	RefreshAllCmd.GroupID = "all"
	ScriptCmd.AddCommand(UnregisterAllCmd)
	ScriptCmd.AddCommand(RefreshAllCmd)
	c := getTestCmd(mediatorscript.ScriptAll)
	c.GroupID = "all"
	ScriptCmd.AddCommand(c)

	TriggerCmd.GroupID = "type"
	ScriptCmd.AddCommand(TriggerCmd)

	interactiveTypes = []mediatorscript.ScriptType{
		mediatorscript.ScriptCondition,
		mediatorscript.ScriptTask,
		mediatorscript.ScriptAssignment,
	}
	for _, t := range interactiveTypes {
		c = getInteractiveScriptCmd(t)
		c.GroupID = "type"
		ScriptCmd.AddCommand(c)
	}

}

type scriptListByType map[mediatorscript.ScriptType][]*mediatorscript.Script

func getAllScriptNamesByType() (scriptListByType, error) {
	var list []*mediatorscript.Script
	if _, err := apiclient.RunGETwithToken("", "json", &list); err != nil {
		return nil, err
	}

	list_by_type := scriptListByType{}
	for _, s := range list {
		if _, ok := list_by_type[s.Type]; !ok {
			list_by_type[s.Type] = make([]*mediatorscript.Script, 1)
			list_by_type[s.Type][0] = s
		} else {
			list_by_type[s.Type] = append(list_by_type[s.Type], s)
		}

	}
	return list_by_type, nil
}
