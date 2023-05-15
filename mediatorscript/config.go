package mediatorscript

import (
	"fmt"
	"strings"
)

type Steps struct {
	Name   string `json:"name,omitempty"`
	Script string `json:"script,omitempty"`
}

type Workflow struct {
	Name  string  `json:"name,omitempty"`
	Steps []Steps `json:"steps,omitempty"`
}
type MediatorConfiguration struct {
	Configuration struct {
		BackendURL    string `json:"backend_url,omitempty"`
		Logfile       string `json:"logfile,omitempty"`
		SSLSkipVerify bool   `json:"ssl_skip_verify,omitempty"`
	} `json:"configuration,omitempty"`
	Workflows []Workflow `json:"workflows,omitempty"`
}

// Looks for a workflow with the provided name in current configuration.
// If more that one workflow are defined with that name, silently returns the first one.
// Sends an error if no worflow with that name could be found.
func (c MediatorConfiguration) GetWorkflow(name string) (Workflow, error) {
	for _, w := range c.Workflows {
		if strings.EqualFold(strings.TrimSpace(w.Name), strings.TrimSpace(name)) {
			return w, nil
		}
	}
	return Workflow{}, fmt.Errorf("workflow '%s' was not found in configuration", name)
}

// Returns the name of the script attached to provided step
// Returns an empty string if no script has been attached to the step
// or if the step is not defined in the workflow,
func (w Workflow) GetScriptForStep(s string) string {
	for _, st := range w.Steps {
		if strings.EqualFold(strings.TrimSpace(st.Name), strings.TrimSpace(s)) {
			// can be an empty string if no script has been attached to this step
			return st.Script
		}
	}
	// if we get here, the step was not found in this workflow.
	// this may be a config error so dump a warning
	logger.Warningf("Step %s was not found in configuration for workflow %s", s, w.Name)
	// return an empty string as if no script were found
	return ""
}

// Returns the name of the step following the step corresponding to given name.
// Returns an error if:
// - the step corresponding to given name could not be found or
// - it's the last step in list
func (w Workflow) GetNextStep(s string) (string, error) {
	stepwasfound := false
	for _, st := range w.Steps {
		if stepwasfound {
			return st.Name, nil
		}
		stepwasfound = strings.EqualFold(strings.TrimSpace(st.Name), strings.TrimSpace(s))
	}

	// if we reach that part, return an error
	if stepwasfound {
		return "", fmt.Errorf("step '%s' was last in workflow '%s' step list. Cannot get next one", s, w.Name)
	}
	return "", fmt.Errorf("step '%s' was not found in workflow '%s' step list. Cannot get next one", s, w.Name)

}
