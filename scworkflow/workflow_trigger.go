package scworkflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"uqtu/mediator/apiclient"
)

type WorkflowTriggers struct {
	WorkflowTriggers struct {
		WorkflowTrigger []WorkflowTrigger `json:"workflow_trigger,omitempty"`
	} `json:"workflow_triggers,omitempty"`
}
type WorkflowTrigger struct {
	ID       int                     `json:"id,omitempty"`
	Name     string                  `json:"name,omitempty"`
	Executer WorkflowTriggerExecuter `json:"executer,omitempty"`
	Triggers []WorkflowTriggerGroup  `json:"triggers,omitempty"`
}

type WorkflowTriggerExecuter struct {
	Type      string `json:"@xsi.type,omitempty"`
	Path      string `json:"arguments,omitempty"`
	Arguments string `json:"path,omitempty"`
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

func GetSecurechangeWorkflowTriggers(username, pwd, host string) (*WorkflowTriggers, error) {
	var (
		wf_triggers WorkflowTriggers
		scURL       string
	)
	if strings.Contains(host, "securechangeworkflow/api/securechange") {
		scURL = host
	} else {
		scURL = fmt.Sprintf("https://%s/securechangeworkflow/api/securechange/", host)
	}
	// get data from securechange
	if client := apiclient.NewClient(scURL, username, pwd, true); client == nil {
		return nil, fmt.Errorf("cannot get API client")

	} else if req, err := client.NewGETwithBasicAuth("/triggers", "json"); err != nil {
		return nil, err

	} else if err := req.Run(&wf_triggers); err != nil {
		return nil, err

	} else {
		return &wf_triggers, nil
	}
}

func CreateSecurechangeWorkflowTriggers(wf_triggers *WorkflowTriggers, username, pwd, host string) error {
	var (
		scURL string
		buff  bytes.Buffer
	)
	if strings.Contains(host, "securechangeworkflow/api/securechange") {
		scURL = host
	} else {
		scURL = fmt.Sprintf("https://%s/securechangeworkflow/api/securechange/", host)
	}

	enc := json.NewEncoder(&buff)
	if err := enc.Encode(wf_triggers); err != nil {
		return err
	}

	// post data to securechange
	if client := apiclient.NewClient(scURL, username, pwd, true); client == nil {
		return fmt.Errorf("cannot get API client")

	} else if req, err := client.NewPOSTwithBasicAuth("/triggers", &buff, "json"); err != nil {
		return err

	} else if _, err := req.RunWithoutDecode(); err != nil {
		return err

	} else {
		return nil
	}
}

func DeleteSecurechangeWorkflowTriggers(wf_trigger_id int, username, pwd, host string) error {
	var (
		scURL string
	)
	if strings.Contains(host, "securechangeworkflow/api/securechange") {
		scURL = host
	} else {
		scURL = fmt.Sprintf("https://%s/securechangeworkflow/api/securechange/", host)
	}

	// send delete request to securechange (undocumented API endpoint)
	if client := apiclient.NewClient(scURL, username, pwd, true); client == nil {
		return fmt.Errorf("cannot get API client")

	} else if req, err := client.NewDELETEwithBasicAuth(fmt.Sprintf("/triggers/%d", wf_trigger_id), "json"); err != nil {
		return err

	} else if _, err := req.RunWithoutDecode(); err != nil {
		return err

	} else {
		return nil
	}
}

func GetSecurechangeWorkflowTriggerByID(id int, username, pwd, host string) (*WorkflowTrigger, error) {

	if list, err := GetSecurechangeWorkflowTriggers(username, pwd, host); err != nil {
		return nil, err
	} else if list != nil {

		for _, w := range list.WorkflowTriggers.WorkflowTrigger {
			if w.ID == id {
				return &w, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
