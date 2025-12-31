package scworkflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"uqtu/mediator/apiclient"
)

type WorkflowTriggers struct {
	WorkflowTriggers struct {
		WorkflowTrigger []*WorkflowTrigger `json:"workflow_trigger,omitempty"`
	} `json:"workflow_triggers,omitempty"`
}
type WorkflowTrigger struct {
	ID       int                     `json:"id,omitempty"`
	Name     string                  `json:"name,omitempty"`
	Executer WorkflowTriggerExecuter `json:"executer,omitempty"`
	Triggers []*WorkflowTriggerGroup `json:"triggers,omitempty"`
}

type WorkflowTriggerExecuter struct {
	Type      string `json:"@xsi.type,omitempty"`
	Path      string `json:"path,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type WorkflowTriggerGroup struct {
	Name     string            `json:"name,omitempty"`
	Workflow WorkflowTriggerWF `json:"workflow,omitempty"`
	Events   []string          `json:"events,omitempty"`
}

type WorkflowTriggerWF struct {
	Name             string `json:"name,omitempty"`
	ParentWorkflowID int    `json:"parent_workflow_id"`
}

func (wtwf *WorkflowTriggerWF) IsRelatedToWorkflow(w *WorkflowXML) bool {
	return wtwf.Name == w.Name
}
func (wtg *WorkflowTriggerGroup) IsRelatedToWorkflow(w *WorkflowXML) bool {
	return wtg.Workflow.IsRelatedToWorkflow(w)
}

func (wt *WorkflowTrigger) IsRelatedToWorkflow(w *WorkflowXML) bool {
	return slices.ContainsFunc(wt.Triggers, func(wtg *WorkflowTriggerGroup) bool {
		return wtg.IsRelatedToWorkflow(w)
	})
}

// check if trigger is related to any of the wf in list
func (wt *WorkflowTrigger) IsTriggerRelatedToWorkflowInList(l []*WorkflowXML) bool {
	return slices.ContainsFunc(
		l,
		func(w *WorkflowXML) bool {
			return wt.IsRelatedToWorkflow(w)
		},
	)
}

// check if a trigger is already in list
func (wt *WorkflowTrigger) IsTriggerAlreadyInList(l []*WorkflowTrigger) bool {
	return slices.ContainsFunc(
		l,
		func(existing *WorkflowTrigger) bool {
			return wt.Equals(existing)
		},
	)
}

func (wt *WorkflowTrigger) Equals(other *WorkflowTrigger) bool {
	// we assume that if they have exactly the same name, they are doing the same job
	if wt.Name == other.Name {
		return true
	}

	if !wt.Executer.Equals(&other.Executer) {
		return false
	}

	for _, t := range wt.Triggers {
		for _, other_t := range other.Triggers {
			if t.Equals(other_t) {
				return true
			}
		}
	}
	return false
}

func (wte *WorkflowTriggerExecuter) Equals(other *WorkflowTriggerExecuter) bool {
	if wte.Arguments != other.Arguments {
		return false
	}

	if wte.Path != other.Path {
		return false
	}
	return true
}

func (wtg *WorkflowTriggerGroup) Equals(other *WorkflowTriggerGroup) bool {
	// we assume that if they have exactly the same name, they are doing the same job
	if wtg.Name == other.Name {
		return true
	}
	if wtg.Workflow.Name != other.Workflow.Name {
		// not same Wf, they are different
		return false
	}

	// check if events are the same.
	if len(wtg.Events) != len(other.Events) {
		return false
	}
	return StringSlicesAreTheSame(wtg.Events, other.Events)
}

func StringSlicesAreTheSame(a, b []string) bool {
	my_events := make([]string, len(a))
	copy(my_events, a)
	other_events := make([]string, len(b))
	copy(other_events, b)
	return slices.Equal(my_events, other_events)
}

func GetSecurechangeWorkflowTriggers(username, pwd, host string) (*WorkflowTriggers, error) {
	var (
		wf_triggers WorkflowTriggers
	)
	c := getSCclient(host)

	// get data from securechange
	if _, err := c.RunGETwithCredentials("/triggers", username, pwd, "json", &wf_triggers); err != nil {
		return nil, err
	} else {
		return &wf_triggers, nil
	}
}

func CreateSecurechangeWorkflowTriggers(wf_triggers *WorkflowTriggers, username, pwd, host string) error {
	var (
		buff bytes.Buffer
	)
	c := getSCclient(host)
	enc := json.NewEncoder(&buff)
	if err := enc.Encode(wf_triggers); err != nil {
		return err
	}

	// post data to securechange
	if _, err := c.RunPOSTwithCredentials("/triggers", username, pwd, &buff, "json", nil); err != nil {
		return err
	} else {
		return nil
	}
}

func DeleteSecurechangeWorkflowTriggers(wf_trigger_id int, username, pwd, host string) error {
	c := getSCclient(host)

	if _, err := c.RunDELETEwithCredentials(fmt.Sprintf("/triggers/%d", wf_trigger_id), username, pwd, "json", nil); err != nil {
		return err
	} else {
		return nil
	}
}

func getSCclient(host string) *apiclient.APIclientHelper {
	var (
		scURL string
	)
	if strings.Contains(host, "securechangeworkflow/api/securechange") {
		scURL = host
	} else {
		scURL = fmt.Sprintf("https://%s/securechangeworkflow/api/securechange/", host)
	}
	c := apiclient.GetHelper(scURL, true)
	return c
}

func GetSecurechangeWorkflowTriggerByID(id int, username, pwd, host string) (*WorkflowTrigger, error) {

	if list, err := GetSecurechangeWorkflowTriggers(username, pwd, host); err != nil {
		return nil, err
	} else if list != nil {

		for _, w := range list.WorkflowTriggers.WorkflowTrigger {
			if w.ID == id {
				return w, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
