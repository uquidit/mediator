package clicommands

import (
	"fmt"
	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
)

var (
	TriggerCmd = &cobra.Command{
		Use:   "trigger <sub-command>",
		Short: "List all trigger scripts mediator-client can use.",
		Long: `Trigger scripts are the scripts called by the mediator-client when a ticket reaches a given step.

There can be as many trigger scripts as needed. They need to be associated to an existing workflow step via UQTU front-end.

SecureChange API also needs to be configured so mediator-client is called when actions are triggered for a ticket.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listScripts(mediatorscript.ScriptTrigger)
		},
	}
)

func getInteractiveScriptCmd(script_type mediatorscript.ScriptType) *cobra.Command {
	cmd := cobra.Command{
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listScripts(script_type)
		},
	}
	switch script_type {
	case mediatorscript.ScriptAssignment:
		cmd.Use = "assignment"
	case mediatorscript.ScriptCondition:
		cmd.Use = "condition"
	case mediatorscript.ScriptTask:
		cmd.Use = "task"
	case mediatorscript.RiskAnalysis:
		cmd.Use = "risk-analysis"
	}

	cmd.Short = fmt.Sprintf("Show %s", script_type)
	cmd.Long = fmt.Sprintf(`Show %s.
This script will be called when Mediator is used in a worflow step as %s.

There can be only one %s.`, script_type, script_type, script_type)

	cmd.AddCommand(getRegisterCommand(script_type))
	cmd.AddCommand(getUnregisterCmd(script_type))
	cmd.AddCommand(getTestCmd(script_type))
	cmd.AddCommand(getRefreshCmd(script_type))

	return &cmd
}

func init() {
	TriggerCmd.AddCommand(getRegisterCommand(mediatorscript.ScriptTrigger))
	TriggerCmd.AddCommand(UnregisterTriggerCmd)
	TriggerCmd.AddCommand(getTestCmd(mediatorscript.ScriptTrigger))
	TriggerCmd.AddCommand(RefreshTriggerCmd)
}
