package mediatorsettings

import (
	"errors"
	"mediator/mediatorscript"
	"testing"
)

func mock_getScriptByName(name string) (*mediatorscript.Script, error) {
	switch name {
	case "trigger":
		return &mediatorscript.Script{
			Fullpath: "/path/to/trigger_script.sh",
			Name:     name,
			Hash:     []byte{},
			Type:     mediatorscript.ScriptTrigger,
		}, nil

	case "condition":
		return &mediatorscript.Script{
			Fullpath: "/path/to/condition_script.sh",
			Name:     name,
			Hash:     []byte{},
			Type:     mediatorscript.ScriptCondition,
		}, nil
	case "pre-assignment":
		return &mediatorscript.Script{
			Fullpath: "/path/to/pre-assignment_script.sh",
			Name:     name,
			Hash:     []byte{},
			Type:     mediatorscript.ScriptAssignment,
		}, nil
	default:
		return nil, mediatorscript.ErrScriptNotFound
	}
}

func TestRule_isValidInner(t *testing.T) {
	type fields struct {
		Trigger string
		Script  string
		Step    string
		NilStep bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				Trigger: "ADVANCE",
				Script:  "trigger",
				Step:    "a good step",
			},
			wantErr: nil,
		},
		{
			name: "ok no step",
			fields: fields{
				Trigger: "Create",
				Script:  "trigger",
				NilStep: true,
			},
			wantErr: nil,
		},
		{
			name: "no trigger",
			fields: fields{
				Trigger: "",
				Script:  "trigger",
				Step:    "xxx",
			},
			wantErr: ErrNoTriggerInRule,
		},
		{
			name: "no script",
			fields: fields{
				Trigger: "redo",
				Script:  "",
				Step:    "xxx",
			},
			wantErr: ErrUnknownScript,
		},
		{
			name: "no step for step trigger",
			fields: fields{
				Trigger: "advance",
				Script:  "trigger",
				NilStep: true,
			},
			wantErr: ErrMissingStepInRule,
		},
		{
			name: "not trigger script",
			fields: fields{
				Trigger: "ADVANCE",
				Script:  "condition",
				Step:    "a good step",
			},
			wantErr: ErrScriptIsNotTriggerScript,
		},
		{
			name: "not trigger script",
			fields: fields{
				Trigger: "ADVANCE",
				Script:  "pre-assignment",
				Step:    "a good step",
			},
			wantErr: ErrScriptIsNotTriggerScript,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rule{
				Trigger: tt.fields.Trigger,
				Script:  tt.fields.Script,
			}
			if !tt.fields.NilStep {
				r.Step = &tt.fields.Step
			}
			if err := r.isValidInner(mock_getScriptByName); !errors.Is(err, tt.wantErr) {
				t.Errorf("Rule.isValidInner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
