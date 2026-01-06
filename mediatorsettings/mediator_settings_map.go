package mediatorsettings

import (
	"fmt"
	"mediator/scworkflow"
)

type MediatorSettingsMap map[string]*WFSettings

func (ms MediatorSettingsMap) SetNextStep(workflows_steps scworkflow.WorkflowsStepsList) []error {
	errs := []error{}

	for _, s := range ms {
		steps := workflows_steps[s.WFname]
		if err := s.SetNextStep(steps); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	} else {
		return errs
	}
}

func (ms MediatorSettingsMap) SetPreviousStep(workflows_steps scworkflow.WorkflowsStepsList) []error {
	errs := []error{}

	for _, s := range ms {
		steps := workflows_steps[s.WFname]
		if err := s.SetPreviousStep(steps); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	} else {
		return errs
	}
}

func (msm MediatorSettingsMap) GetSlice() MediatorSettings {
	ms := MediatorSettings{}
	for _, s := range msm {
		ms = append(ms, s)
	}
	return ms

}

// Looks for a workflow with the provided name in current settings.
// If more that one workflow are defined with that name, silently returns the first one.
// Sends an error if no worflow with that name could be found.
func (wm MediatorSettingsMap) GetWorkflowSettings(name string) (*WFSettings, error) {
	if w, ok := wm[name]; ok {
		return w, nil
	} else {
		return nil, fmt.Errorf("workflow '%s' was not found in settings", name)
	}
}

func (msm MediatorSettingsMap) Clean() {
	for wf_name, ms := range msm {
		if len(ms.Rules) == 0 {
			delete(msm, wf_name)
		}
	}
}
