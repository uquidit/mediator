package securechangeapi

import (
	"fmt"
	"mediator/console"
	"mediator/scworkflow"
	"slices"
	"sort"

	"github.com/spf13/cobra"
)

type my_trigger struct {
	id        int
	name      string
	path      string
	arguments string
}

func (mt my_trigger) GetLabel() string {
	return fmt.Sprintf("%s: mediator '%s', arguments '%s' (#%d)", mt.name, mt.path, mt.arguments, mt.id)
}
func (mt my_trigger) GetValue() int {
	return mt.id
}

func selectWorkflow(msg string) (*scworkflow.WorkflowXML, error) {
	wf_list, err := getWorkflowsFromSecurechangeIfNeeded()
	if err != nil {
		return nil, err
	}

	items := []string{}
	for _, w := range wf_list {
		items = append(items, w.Name)
	}
	sort.Strings(items)
	items = append(items, " - Exit -")
	for {
		if wf_name, err := console.SelectFromList(msg, items, nil); err == nil && wf_name != "" {
			if wf_name == " - Exit -" {
				return nil, nil
			}
			for _, w := range wf_list {
				if w.Name == wf_name {
					return w, nil
				}
			}
		}
	}
}

var (
	scWorkflows []*scworkflow.WorkflowXML
)

func getWorkflowsFromSecurechangeIfNeeded() ([]*scworkflow.WorkflowXML, error) {
	if scWorkflows == nil {
		// send request to SC
		if l, err := Manager.GetSecurechangeWorkflows(false); err != nil {
			return nil, err
		} else if len(l.Workflows) == 0 {
			return nil, ErrNoSecurechangeWorkflows
		} else {
			scWorkflows = l.Workflows
		}
	}
	return scWorkflows, nil
}

func getWorkflowsFromFlag(cmd *cobra.Command) []*scworkflow.WorkflowXML {
	WF_to_process, err := cmd.Flags().GetStringSlice("workflow")
	if err != nil || len(WF_to_process) == 0 {
		return nil
	}
	wf_list, err := getWorkflowsFromSecurechangeIfNeeded()
	if len(wf_list) == 0 || err != nil {
		return nil
	}
	l := []*scworkflow.WorkflowXML{}
	for _, w := range wf_list {
		if slices.Contains(WF_to_process, w.Name) {
			l = append(l, w)
		}
	}
	return l
}

// check if trigger is related to any of the wf in list
func isTriggerRelatedToWorkflowInList(l []*scworkflow.WorkflowXML, t *scworkflow.WorkflowTrigger) bool {
	return t.IsTriggerRelatedToWorkflowInList(l)
}
