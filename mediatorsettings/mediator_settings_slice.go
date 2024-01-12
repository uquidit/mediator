package mediatorsettings

import "uqtu/mediator/scworkflow"

type MediatorSettings []*WFSettings

func (ms MediatorSettings) SetPreviousStep(workflows_steps scworkflow.WorkflowsStepsList) []error {
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

func (ms MediatorSettings) SetNextStep(workflows_steps scworkflow.WorkflowsStepsList) []error {
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

func (ms MediatorSettings) GetMap() MediatorSettingsMap {
	msm := MediatorSettingsMap{}
	for _, s := range ms {
		msm[s.WFname] = s
	}
	return msm
}
