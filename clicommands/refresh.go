package clicommands

import (
	"fmt"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var (
	RefreshTriggerCmd = &cobra.Command{
		Use:   "refresh [scriptname]",
		Short: "Refresh one or all Trigger script registrations",
		Long: `Run this command to refresh one or all Trigger script registration.
	
Provide a script name to refresh only this script registration.
If no script name is provided, all Trigger script registrations will be refreshed.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return refreshScript(mediatorscript.ScriptTrigger, args[0])
			} else {
				return refreshScript(mediatorscript.ScriptTrigger, "")
			}
		},
	}
	RefreshAllCmd = &cobra.Command{
		Use:   "refresh",
		Short: "Refresh all script registrations",
		Long:  `Run this command to refresh all script registrations.`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return refreshScript(mediatorscript.ScriptAll, "")
		},
	}
)

func getRefreshCmd(script_type mediatorscript.ScriptType) *cobra.Command {
	cmd := cobra.Command{
		Use:  "refresh",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return refreshScript(script_type, "")
		},
	}
	cmd.Short = fmt.Sprintf("Refresh %s registration", script_type)
	cmd.Long = fmt.Sprintf("Run this command to refresh %s registration after the file has been modified.", script_type)
	return &cmd
}

func refreshScript(script_type mediatorscript.ScriptType, name string) error {
	var endpoint string

	if script_type == mediatorscript.ScriptAll {
		endpoint = "refresh-all"
	} else if name != "" {
		endpoint = fmt.Sprintf("refresh/%s/%s", script_type.Slug(), name)
	} else {
		endpoint = fmt.Sprintf("refresh/%s", script_type.Slug())
	}

	if _, err := apiclient.RunPOSTwithToken(endpoint, nil, "json", nil); err != nil {
		return err

	} else if name != "" {
		fmt.Printf("Script '%s' has been refreshed.\n", name)
		return nil

	} else if script_type == mediatorscript.ScriptAll {
		fmt.Println("All scripts have been refreshed.")
		return nil

	} else {
		fmt.Printf("%s has been refreshed.\n", script_type)
		return nil
	}

}
