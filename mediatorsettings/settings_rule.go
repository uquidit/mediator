package mediatorsettings

import (
	"fmt"
)

// Stores a list of step <=> script association
// index is step name
// value is script name
type Rule struct {
	Trigger string  `json:"trigger,omitempty"`
	Script  string  `json:"script,omitempty"`
	Step    *string `json:"step,omitempty"`
}
type RulesSlice []*Rule

func (r Rule) String() string {
	if r.Script == "" {
		return ""
	}
	if r.Step != nil && *r.Step != "" {
		return fmt.Sprintf("Trigger %s on step %s fires script %s", r.Trigger, *r.Step, r.Script)
	} else {
		return fmt.Sprintf("Trigger %s fires script %s", r.Trigger, r.Script)
	}
}
