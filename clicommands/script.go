package clicommands

import (
	"fmt"
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

			// list all scripts
			var list []*mediatorscript.Script
			if _, err := client.RunGETwithToken("", "json", &list); err != nil {
				return err
			}

			lines := map[mediatorscript.ScriptType][]string{}
			for _, s := range list {
				line := fmt.Sprintf("  - %s: %s\n", s.Name, s.Fullpath)

				if _, ok := lines[s.Type]; !ok {
					lines[s.Type] = make([]string, 1)
					lines[s.Type][0] = line
				} else {
					lines[s.Type] = append(lines[s.Type], line)
				}

			}

			fmt.Printf("Nb of scripts: %d\n", len(lines))
			for _, t := range []mediatorscript.ScriptType{
				mediatorscript.ScriptTrigger,
				mediatorscript.ScriptCondition,
				mediatorscript.ScriptTask,
				mediatorscript.ScriptAssignment,
			} {
				fmt.Printf("\n* %s: %d\n", t, len(lines[t]))
				for _, s := range lines[t] {
					fmt.Print(s)
				}
			}

			return nil
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
