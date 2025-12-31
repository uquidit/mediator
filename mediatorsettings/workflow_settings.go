package mediatorsettings

import (
	"fmt"
	"uqtu/mediator/scworkflow"
)

type WFSettings struct {
	WFname      string     `json:"wf_name,omitempty"`
	WFid        int        `json:"wf_id"`
	Rules       RulesSlice `json:"settings"`
	Description string     `json:"description"`
}

// Checks if settings are valid:
// - wf name is not empty
// - wf ID is set
// - all rules are valid
func (w *WFSettings) isValid() error {
	if w.WFname == "" {
		return ErrNoWorkflowName
	}
	if w.WFid == 0 {
		return ErrNoWorkflowID
	}
	for _, r := range w.Rules {
		if err := r.isValid(); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidSettings, err)
		}
	}
	return nil
}

// Returns the names of the scripts attached to provided step and trigger
// Returns an empty string if no script has been attached to the step
// or if the step is not defined in the workflow,
func (w *WFSettings) GetScriptsForTriggerAndStep(t scworkflow.SecurechangeTrigger, step string) []string {
	scripts_list := []string{}
	for _, rule := range w.Rules {
		if rule.Script == "" {
			continue
		}
		var step_match bool
		if rule.Step == nil {
			// if no step, we don't care about what was given in arguments list
			// always a match
			step_match = true
		} else {
			step_match = *rule.Step == step
		}

		if rule.Trigger == t.String() && step_match {
			scripts_list = append(scripts_list, rule.Script)
		}
	}
	return scripts_list
}

// Returns the name of all the script attached to the workflow for the given trigger.
// This is useful when testing worflow scripts.
// Return a slice of strings containing script names.
func (w WFSettings) GetAllScripts(t scworkflow.SecurechangeTrigger) []string {
	scripts_map := map[string]struct{}{}
	trigger_name := t.String()
	for _, rule := range w.Rules {
		if rule.Script == "" {
			continue
		}
		if rule.Trigger != trigger_name {
			continue
		}
		scripts_map[rule.Script] = struct{}{}
	}

	scripts_list := []string{}
	for s := range scripts_map {
		scripts_list = append(scripts_list, s)
	}
	return scripts_list

}

// look for rules involving triggers using next step
// and change the step to the previous one
// according to provided step list
func (w *WFSettings) SetPreviousStep(steps []string) error {
	for i, setting := range w.Rules {
		trigger := scworkflow.GetTriggerFromString(setting.Trigger)
		if trigger.UseNextStep() {
			for j, step := range steps {
				if step == *setting.Step {
					if j-1 >= 0 {
						s := steps[j-1]
						w.Rules[i].Step = &s
						break
					} else {
						return fmt.Errorf("%w in '%s' settings: cannot run a script on first step %s for trigger %s", ErrInvalidRule, w.WFname, *setting.Step, setting.Trigger)
					}
				}
			}
		}
	}
	return nil
}

// look for rules involving triggers using next step
// and change the step to the next one
// according to provided step list
func (w *WFSettings) SetNextStep(steps []string) error {
	for i, setting := range w.Rules {
		trigger := scworkflow.GetTriggerFromString(setting.Trigger)
		if trigger.UseNextStep() {
			for j, step := range steps {
				if step == *setting.Step {
					if j+1 < len(steps) {
						s := steps[j+1]
						w.Rules[i].Step = &s
						break
					} else {
						return fmt.Errorf("%w in '%s' settings: cannot run a script on last step %s for trigger %s", ErrInvalidRule, w.WFname, *setting.Step, setting.Trigger)
					}
				}
			}
		}
	}
	return nil
}

// Remove useless or empty rules
// set Rules to nil if no rules left
func (w *WFSettings) Clean() {
	if len(w.Rules) == 0 {
		w.Rules = nil
		return
	}
	// remove empty rules
	cleaned_settings := RulesSlice{}
	for _, rule := range w.Rules {
		if rule == nil {
			continue
		}
		trigger := scworkflow.GetTriggerFromString(rule.Trigger)
		if rule.Script == "" ||
			trigger == scworkflow.NO_TRIGGER ||
			(trigger.NeedStepToGetScript() && rule.Step == nil) ||
			(trigger.NeedStepToGetScript() && rule.Step != nil && *rule.Step == "") {
			continue
		}
		cleaned_settings = append(cleaned_settings, rule)
	}
	if len(cleaned_settings) > 0 {
		w.Rules = cleaned_settings
	} else {
		w.Rules = nil
	}

}
