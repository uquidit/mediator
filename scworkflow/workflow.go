package scworkflow

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type WorkflowXML struct {
	Id          int              `xml:"id" json:"id"`
	Name        string           `xml:"name" json:"name"`
	Steps       WorkflowStepsXML `xml:"steps>step" json:"steps"`
	HasSettings bool
}
type WorkflowStepsXML []struct {
	Name     string `xml:"name" json:"name"`
	IsActive bool   `xml:"is_active" json:"is_active"`
}
type Workflows struct {
	XMLName   xml.Name       `xml:"workflows" json:"-"`
	Workflows []*WorkflowXML `xml:"workflow" json:"workflow"`
}

type WorkflowsStepsList map[string][]string

func (scws Workflows) GetWorkflowByID(id int) *WorkflowXML {
	for _, w := range scws.Workflows {
		if w.Id == id {
			return w
		}
	}
	return nil
}

func (scws Workflows) GetWorkflowsSteps() WorkflowsStepsList {
	workflows_steps := WorkflowsStepsList{}
	for _, w := range scws.Workflows {
		workflows_steps[w.Name] = w.GetSteps()
	}
	return workflows_steps
}

func (w WorkflowXML) GetSteps() []string {
	if len(w.Steps) == 0 {
		return nil
	}
	steps := []string{}
	for _, s := range w.Steps {
		if !s.IsActive {
			continue
		}
		steps = append(steps, s.Name)
	}
	return steps
}
func (w WorkflowXML) GetLabel() string {
	if w.HasSettings {
		return w.Name
	}
	return fmt.Sprintf("%s !", w.Name)
}
func (w WorkflowXML) GetValue() int {
	return w.Id
}

func GetSecurechangeWorkflows(username, pwd, host string, get_steps bool) (*Workflows, error) {
	rr := credentials_requester{
		username: username,
		password: pwd,
	}
	if strings.HasPrefix(host, "http") {
		rr.url = host
	} else {
		rr.url = fmt.Sprintf("https://%s/securechangeworkflow/api/securechange/", host)
	}
	return GetSecurechangeWorkflowsUsingRequester(rr, get_steps)
}

func GetSecurechangeWorkflowsUsingRequester(rr Requester, get_steps bool) (*Workflows, error) {
	var (
		workflowsSecureChange Workflows
	)

	// get data from securechange
	if req, err := rr.GetRequest("/workflows/active_workflows"); err != nil {
		return nil, err

	} else if err := req.Run(&workflowsSecureChange); err != nil {
		return nil, err

	} else if get_steps {

		// create buffer channels so they're non-blocking
		ch := make(chan *WorkflowXML, len(workflowsSecureChange.Workflows))
		err_ch := make(chan error, len(workflowsSecureChange.Workflows))

		// call go routine for each workflow
		for _, element := range workflowsSecureChange.Workflows {
			go func(id int, c chan *WorkflowXML, ec chan error) {

				// get detailed data for each workflow
				endpoint := fmt.Sprintf("/workflows?id=%d", id)
				var workflow WorkflowXML
				if req, err := rr.GetRequest(endpoint); err != nil {
					ec <- err //buffered channel: non blocking

				} else if err := req.Run(&workflow); err != nil {
					ec <- err

				} else {
					c <- &workflow
				}
			}(element.Id, ch, err_ch)
		}

		// store detailed list in place of light list
		l := []*WorkflowXML{}
		for {
			select {
			case err := <-err_ch:
				// return at first error.
				// hopefully any remaining routines will end nicely (buffered channels)
				return nil, err
			case wf := <-ch:
				l = append(l, wf)
			default:
				if len(l) == len(workflowsSecureChange.Workflows) {
					// we got them all
					workflowsSecureChange.Workflows = l
					return &workflowsSecureChange, nil
				}
			}
		}

	}
	return &workflowsSecureChange, nil

}
