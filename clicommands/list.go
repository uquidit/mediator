package clicommands

import (
	"fmt"
	"strings"
	"uqtu/mediator/apiclient"
	"uqtu/mediator/mediatorscript"
)

func listScripts(script_type mediatorscript.ScriptType) error {
	// list all scripts
	var list []*mediatorscript.Script
	if _, err := apiclient.RunGETwithToken("", "json", &list); err != nil {
		return err
	}

	list_lines := []string{}
	for _, s := range list {
		if s.Type == script_type {
			list_lines = append(list_lines, fmt.Sprintf("- %s: %s\n", s.Name, s.Fullpath))
		}
	}
	fmt.Printf("Nb of %s: %d\n", script_type, len(list_lines))
	fmt.Print(strings.Join(list_lines, ""))

	return nil
}
