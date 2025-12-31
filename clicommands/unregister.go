package clicommands

import (
	"fmt"
	"net/http"
	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var (
	ignore_unregister_notfound_error bool
	UnregisterTriggerCmd             = &cobra.Command{
		Use:   "unregister [scriptname]",
		Short: "Unregister one or all Trigger script.",
		Long: `Run this command to unregister one or all Trigger script.
If no script name is provided, all Trigger script will be unregistered.

Unregistered scripts will no longer be available for Mediator to use.

This command will not delete references to this file in Mediator configuration.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return unregisterScript(mediatorscript.ScriptTrigger, args[0])
			} else {
				return unregisterScript(mediatorscript.ScriptTrigger, "")
			}
		},
	}
	UnregisterAllCmd = &cobra.Command{
		Use:   "unregister",
		Short: "Unregister all the scripts Mediator client can use.",
		Long: `Run this command to unregister all registered scripts.

Unregistered scripts will no longer be available for Mediator to use.
		
This command will not delete references to this file in Mediator configuration.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unregisterScript(mediatorscript.ScriptAll, "")
		},
	}
)

func getUnregisterCmd(script_type mediatorscript.ScriptType) *cobra.Command {
	cmd := cobra.Command{
		Use:  "unregister",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unregisterScript(script_type, "")
		},
	}
	cmd.Short = fmt.Sprintf("Unregister %s", script_type)
	cmd.Long = fmt.Sprintf(`Run this command to unregister %s.
It will no longer be available for Mediator to use.

This command will not delete references to this file in Mediator configuration.`, script_type)
	return &cmd
}

func unregisterScript(script_type mediatorscript.ScriptType, name string) error {
	var endpoint string
	if script_type == mediatorscript.ScriptAll {
		endpoint = "unregister-all"
	} else if name != "" {
		endpoint = fmt.Sprintf("unregister/%s/%s", script_type.Slug(), name)
	} else {
		endpoint = fmt.Sprintf("unregister/%s", script_type.Slug())
	}

	if _, err := BackendClient.RunDELETEwithToken(endpoint, "json", nil); err != nil {
		if ignore_unregister_notfound_error && BackendClient.GetLastRequestStatusCode() == http.StatusNotFound {
			fmt.Println("[WARNING] Script was not found.")
			return nil
		}
		return err

	} else if name != "" {
		fmt.Printf("Script '%s' has been unregistered.\n", name)
		return nil

	} else if script_type == mediatorscript.ScriptAll {
		fmt.Println("All scripts have been unregistered.")
		return nil

	} else {
		fmt.Printf("%s has been unregistered.\n", script_type)
		return nil
	}
}

func init() {
	UnregisterTriggerCmd.Flags().BoolVar(&ignore_unregister_notfound_error, "ignore-not-found", false, "Do not return an error is the provided script was not found in registry.")
}
