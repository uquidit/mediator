package mediatorsettings

import (
	"fmt"
	"uqtu/mediator/mediatorscript"
	"uqtu/mediator/scworkflow"
)

// Stores a list of step <=> script association
// index is step name
// value is script name
type Rule struct {
	Trigger string  `json:"trigger,omitempty"`
	Script  string  `json:"script,omitempty"`
	Step    *string `json:"step,omitempty"`
	Comment string  `json:"comment"`
}
type RulesSlice []*Rule

func (r Rule) String() string {
	if r.Script == "" {
		return ""
	}
	var s string
	if r.Step != nil && *r.Step != "" {
		s = fmt.Sprintf("Trigger %s on step %s fires script %s.", r.Trigger, *r.Step, r.Script)
	} else {
		s = fmt.Sprintf("Trigger %s fires script %s.", r.Trigger, r.Script)
	}
	if r.Comment != "" {
		s = fmt.Sprintf("%s Comment: %s", s, r.Comment)
	}
	return s
}

// Check if a rule is valid:
// - trigger is set
// - step is set if required by trigger
// - script is set and is a trigger script
func (r Rule) isValid() error {
	return r.isValidInner(mediatorscript.GetScriptByName)
}

func (r Rule) isValidInner(getScriptByName func(name string) (*mediatorscript.Script, error)) error {
	trigger := scworkflow.GetTriggerFromString(r.Trigger)
	if trigger == scworkflow.NO_TRIGGER {
		return ErrNoTriggerInRule
	}

	if trigger.NeedStepToGetScript() && (r.Step == nil || *r.Step == "") {
		return ErrMissingStepInRule
	}

	script, err := getScriptByName(r.Script)
	if err != nil {
		return fmt.Errorf("%w: '%s'", ErrUnknownScript, r.Script)
	}

	if script.Type != mediatorscript.ScriptTrigger {
		return fmt.Errorf("%w: '%s'", ErrScriptIsNotTriggerScript, script.Name)
	}
	return nil
}
