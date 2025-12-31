package clicommands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var (
	name_flg string
)

// return a Cobra "register" sub-command for provided script type.
// This command should be added to parent command.
func getRegisterCommand(script_type mediatorscript.ScriptType) *cobra.Command {
	var (
		long string = `Script path can be full or relative.
	
The provided script will be registered using the name of the file.
Example: the script /a/b/c/script.sh will be registered under the name "script.sh"

Script name must be unique. If, for some reason, you need to register more than one
script with the same name, use the --name flag to provide a different name.

If the file changes, the script must be refreshed using the "refresh" command.`
	)

	cmd := cobra.Command{
		Use:   "register <script path>",
		Short: fmt.Sprintf("Register a %s to be used by Mediator client.", script_type),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return registerScript(script_type, args[0])
		},
	}

	switch script_type {
	case mediatorscript.ScriptTrigger:
		cmd.Long = fmt.Sprintf("Register a new Trigger script.\n\n%s\n\nScript name must be exactly the same as what is provided in MediatorScript configuration.", long)

	default:
		txt := fmt.Sprintf("Only one script can be register for that purpose. If a %s has previously been registered, any subsequent registration will overwrite the previous one.", script_type)
		cmd.Long = fmt.Sprintf("Register the script that will be called when Mediator is used in a worflow step as %s.\n\n%s\n\n%s",
			script_type,
			txt,
			long,
		)
	}

	// add --name flag
	cmd.Flags().StringVarP(&name_flg, "name", "n", "", "Script name")

	return &cmd
}

// Common function to send a register request to back-end.
// If no name has been provided via the --name flag,
// create a name using the file name from the path.
// Keep the extension so less risk of collision.
func registerScript(script_type mediatorscript.ScriptType, path string) error {
	if fp, err := filepath.Abs(path); err != nil {
		return err
	} else {
		s := mediatorscript.Script{
			Fullpath: fp,
			Type:     script_type,
		}

		// If no name has been provided via a flag, create a name
		// Use the file name. Get it from the path
		// Keep the extension so less risk of collision
		if name_flg == "" {
			s.Name = filepath.Base(path)
		} else {
			s.Name = name_flg
		}

		if jsoninput, err := json.Marshal(s); err != nil {
			return err
		} else if _, err := BackendClient.RunPOSTwithToken("register", bytes.NewBuffer(jsoninput), "json", nil); err != nil {
			return err

		} else {
			fmt.Printf("%s '%s' has been registered and linked to file %s\n", script_type, s.Name, s.Fullpath)
		}
	}
	return nil

}
