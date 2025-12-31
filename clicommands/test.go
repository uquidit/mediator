package clicommands

import (
	"fmt"
	"uqtu/mediator/mediatorscript"

	"github.com/spf13/cobra"
)

func getTestCmd(script_type mediatorscript.ScriptType) *cobra.Command {
	cmd := cobra.Command{}

	switch script_type {
	case mediatorscript.ScriptTrigger:
		cmd.Short = "Test one or all Trigger scripts"
		cmd.Use = "test [script name]"
		cmd.Long = `If no argument is provided, test all registered trigger scripts.
If a script name is provided, test only that trigger script.

Equivalent to SecureChange TEST button.`
		cmd.Args = cobra.MaximumNArgs(1)
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return testScript(script_type, args[0])
			} else {
				return testScript(script_type, "")
			}
		}

	case mediatorscript.ScriptAll:
		cmd.Short = "Test all scripts"
		cmd.Use = "test"
		cmd.Long = `Test all registered scripts:
* Trigger scripts
* Scripted Condition script
* Scripted Task script
* Pre-Assignment script

Equivalent to SecureChange TEST button.`
		cmd.Args = cobra.MaximumNArgs(0)
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return testScript(script_type, "")
		}
	default:
		cmd.Short = fmt.Sprintf("Test %s", script_type)
		cmd.Use = "test"
		cmd.Long = fmt.Sprintf("Test %s. Equivalent to SecureChange TEST button", script_type)
		cmd.Args = cobra.MaximumNArgs(0)
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return testScript(script_type, "")
		}
	}

	return &cmd
}

func testScript(script_type mediatorscript.ScriptType, name string) error {

	var endpoint string
	if script_type == mediatorscript.ScriptAll {
		endpoint = "test-all"
	} else if name != "" {
		endpoint = fmt.Sprintf("test/%s/%s", script_type.Slug(), name)
	} else {
		endpoint = fmt.Sprintf("test/%s", script_type.Slug())
	}

	results := mediatorscript.RunResponse{}
	if _, err := BackendClient.RunPOSTwithToken(endpoint, nil, "json", &results); err == nil {

		// double check for errors
		if results.Error != "" {
			return fmt.Errorf("internal error: %s", results.Error)
		}

		fmt.Println("Test result: ")
		nberrors := 0
		for name, res := range results.RunResults {
			fmt.Printf("* %s:\n", name)

			// check for any errors before script execution
			if res.InternalError != "" {
				fmt.Printf("   - Internal error: %s\n", res.InternalError)
				fmt.Printf("   - test %s: FAILED\n", name)
				nberrors += 1
				continue
			}

			if res.ScriptError != "" {
				fmt.Printf("   - Script error: %s\n", res.ScriptError)

			}

			if res.StdOut != "" {
				fmt.Printf("   - script output: %s\n", res.StdOut)
			}
			if res.StdErr != "" {
				fmt.Printf("   - script error: %s\n", res.StdErr)
			}
			if res.ScriptError != "" {
				fmt.Printf("   - execution error: %s\n", res.ScriptError)
			}

			switch res.Type {
			case mediatorscript.ScriptTrigger, mediatorscript.ScriptAssignment:
				if res.ExitCode == 0 {
					fmt.Printf("   - test %s: OK\n", name)
				} else {
					fmt.Printf("   - test %s: FAILED\n", name)
					nberrors += 1
				}
			case mediatorscript.ScriptCondition, mediatorscript.ScriptTask, mediatorscript.RiskAnalysis:
				if res.ExitCode == 0 && res.StdOut == "<response><condition_result> true </condition_result></response>" {
					fmt.Printf("   - test %s: OK\n", name)
				} else {
					fmt.Printf("   - test %s: FAILED\n", name)
					nberrors += 1
				}
			default:
				fmt.Printf("unexecpected type: %s", res.Type)
				return fmt.Errorf("unexecpected type: %s", res.Type)

			}
		}
		if nberrors == 0 {
			fmt.Println("All tests passed")
		} else {
			fmt.Printf("%d/%d test(s) failed\n", nberrors, len(results.RunResults))
		}
		return nil

	} else {
		return err
	}

}
