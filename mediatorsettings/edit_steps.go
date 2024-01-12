package mediatorsettings

import (
	"errors"
	"fmt"
	"uqtu/mediator/scworkflow"

	"github.com/labstack/echo/v4"
)

type msMapOrSlice interface {
	MediatorSettings | MediatorSettingsMap
	SetPreviousStep(workflows_steps scworkflow.WorkflowsStepsList) []error
	SetNextStep(workflows_steps scworkflow.WorkflowsStepsList) []error
}

func editSteps[dataType msMapOrSlice](data dataType, set_previous bool, c echo.Context) error {
	if sc_username := c.QueryParam("sc_username"); sc_username == "" {
		return errors.New("SC username is missing")
	} else if sc_password := c.QueryParam("sc_password"); sc_password == "" {
		return errors.New("SC password is missing")
	} else if sc_host := c.QueryParam("sc_host"); sc_host == "" {
		return errors.New("SC host is missing")
	} else if workflows, err := scworkflow.GetSecurechangeWorkflows(sc_username, sc_password, sc_host, true); err != nil {
		return fmt.Errorf("cannot get SC workflows: %w", err)
	} else {
		workflows_steps := workflows.GetWorkflowsSteps()
		if set_previous {
			if errs := data.SetPreviousStep(workflows_steps); len(errs) > 0 {
				return fmt.Errorf("error while setting previous step: %v", errs)
			}
		} else {
			if errs := data.SetNextStep(workflows_steps); len(errs) > 0 {
				return fmt.Errorf("error while setting next step: %v", errs)
			}
		}
	}
	return nil
}
