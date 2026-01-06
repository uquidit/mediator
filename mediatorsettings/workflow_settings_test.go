package mediatorsettings

import (
	"mediator/scworkflow"
	"slices"
	"testing"

	"github.com/go-test/deep"
)

var steps []string = []string{"step1", "step2", "step3", "step4"}

func TestWFSettings_SetPreviousStep(t *testing.T) {

	tests := []struct {
		name      string
		rules     RulesSlice
		wantRules RulesSlice
		wantErr   bool
	}{
		{
			name:      "no rules",
			rules:     nil,
			wantRules: nil,
			wantErr:   false,
		},
		{
			name:      "empty rules",
			rules:     []*Rule{},
			wantRules: []*Rule{},
			wantErr:   false,
		},
		{
			name: "2 simple triggers: no change",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Close",
					Script:  "script2.py",
					Step:    nil,
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Close",
					Script:  "script2.py",
					Step:    nil,
				},
			},
			wantErr: false,
		},
		{
			name: "1 simple trigger, 1 using steps: no change",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[1],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[1],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger using first step: no change",
			rules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger using last step: no change",
			rules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger (next step) using last step: change",
			rules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[len(steps)-2],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger (next step) using first step: error",
			rules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WFSettings{
				WFname: "WF",
				Rules:  tt.rules,
			}
			if err := w.SetPreviousStep(steps); (err != nil) != tt.wantErr {
				t.Errorf("WFSettings.SetPreviousStep() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := deep.Equal(w.Rules, tt.wantRules); diff != nil {
				t.Errorf("WFSettings.SetPreviousStep() error: returned %v, want %v. Diffs are %v", w.Rules, tt.wantRules, diff)
			}
		})
	}
}

func TestWFSettings_SetNextStep(t *testing.T) {

	tests := []struct {
		name      string
		rules     RulesSlice
		wantRules RulesSlice
		wantErr   bool
	}{
		{
			name:      "no rules",
			rules:     nil,
			wantRules: nil,
			wantErr:   false,
		},
		{
			name:      "empty rules",
			rules:     []*Rule{},
			wantRules: []*Rule{},
			wantErr:   false,
		},
		{
			name: "2 simple triggers: no change",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Close",
					Script:  "script2.py",
					Step:    nil,
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Close",
					Script:  "script2.py",
					Step:    nil,
				},
			},
			wantErr: false,
		},
		{
			name: "1 simple trigger, 1 using steps: no change",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[1],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Create",
					Script:  "script.py",
					Step:    nil,
				},
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[1],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger using first step: no change",
			rules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger using last step: nochange",
			rules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Redo",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantErr: false,
		},
		{
			name: "1 trigger (next step) using last step: error",
			rules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[len(steps)-1],
				},
			},
			wantErr: true,
		},
		{
			name: "1 trigger (next step) using first step: change",
			rules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[0],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "script2.py",
					Step:    &steps[1],
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WFSettings{
				WFname: "WF",
				Rules:  tt.rules,
			}
			if err := w.SetNextStep(steps); (err != nil) != tt.wantErr {
				t.Errorf("WFSettings.SetNextStep() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := deep.Equal(w.Rules, tt.wantRules); diff != nil {
				t.Errorf("WFSettings.SetNextStep() error: returned %v, want %v. Diffs are %v", w.Rules, tt.wantRules, diff)
			}
		})
	}
}

func TestWFSettings_Clean(t *testing.T) {
	tests := []struct {
		name      string
		rules     RulesSlice
		wantRules RulesSlice
	}{
		{
			name:      "no rule",
			rules:     nil,
			wantRules: nil,
		},
		{
			name:      "empty rules",
			rules:     []*Rule{},
			wantRules: nil,
		},
		{
			name: "3 ok rules: no change",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "toto",
				},
				{
					Trigger: "Close",
					Script:  "titi",
				},
				{
					Trigger: "Advance",
					Script:  "tata",
					Step:    &steps[2],
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Create",
					Script:  "toto",
				},
				{
					Trigger: "Close",
					Script:  "titi",
				},
				{
					Trigger: "Advance",
					Script:  "tata",
					Step:    &steps[2],
				},
			},
		},
		{
			name: "1 ok rules, 4 empty: 1 left",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "",
				},
				nil,
				{
					Trigger: "Close",
					Script:  "titi",
				},
				{
					Trigger: "Advance",
					Script:  "tata",
					Step:    nil,
				},
				{
					Trigger: "",
					Script:  "ssss",
					Step:    nil,
				},
			},
			wantRules: []*Rule{
				{
					Trigger: "Close",
					Script:  "titi",
				},
			},
		},
		{
			name: "3 empty: return nil",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "",
				},
				{
					Trigger: "Advance",
					Script:  "tata",
					Step:    nil,
				},
				{
					Trigger: "",
					Script:  "ssss",
					Step:    nil,
				},
			},
			wantRules: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WFSettings{
				Rules: tt.rules,
			}
			w.Clean()
			if diff := deep.Equal(w.Rules, tt.wantRules); diff != nil {
				t.Errorf("WFSettings.Clean() error: returned %v, want %v. Diffs are %v", w.Rules, tt.wantRules, diff)
			}
		})
	}
}

func TestWFSettings_GetAllScripts(t *testing.T) {
	tests := []struct {
		name        string
		rules       RulesSlice
		arg_trigger scworkflow.SecurechangeTrigger
		want        []string
	}{
		{
			name:        "no rules",
			rules:       nil,
			arg_trigger: scworkflow.ADVANCE,
			want:        []string{},
		},
		{
			name:        "empty rules",
			rules:       []*Rule{},
			arg_trigger: scworkflow.ADVANCE,
			want:        []string{},
		},
		{
			name: "3 steps",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "Want that one",
				},
				{
					Trigger: "Create",
					Script:  "Want that one too",
				},
				{
					Trigger: "Close",
					Script:  "Dont want",
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[2],
				},
				{
					Trigger: "Create",
					Script:  "Last one I want",
				},
			},
			arg_trigger: scworkflow.CREATE,
			want:        []string{"Want that one", "Want that one too", "Last one I want"},
		},
		{
			name: "0 steps",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "Want that one? no",
				},
				{
					Trigger: "Create",
					Script:  "",
				},
				{
					Trigger: "Close",
					Script:  "Dont want",
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[2],
				},
				{
					Trigger: "Create",
					Script:  "Last one I don't want",
				},
			},
			arg_trigger: scworkflow.REDO,
			want:        []string{},
		},
		{
			name: "0 steps no trigger",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "Want that one? no",
				},
				{
					Trigger: "Create",
					Script:  "Want that one too... nah",
				},
				{
					Trigger: "Close",
					Script:  "Dont want",
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[2],
				},
				{
					Trigger: "Create",
					Script:  "Last one I don't want",
				},
			},
			arg_trigger: 0,
			want:        []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WFSettings{
				Rules: tt.rules,
			}
			got := w.GetAllScripts(tt.arg_trigger)
			slices.Sort(tt.want)
			slices.Sort(got)
			if slices.Compare(got, tt.want) != 0 {
				t.Errorf("WFSettings.GetAllScripts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWFSettings_GetScriptsForTriggerAndStep(t *testing.T) {
	tests := []struct {
		name        string
		rules       RulesSlice
		arg_trigger scworkflow.SecurechangeTrigger
		arg_step    string
		want        []string
	}{
		{
			name:        "no rules",
			rules:       nil,
			arg_trigger: scworkflow.ADVANCE,
			want:        []string{},
		},
		{
			name:        "empty rules",
			rules:       []*Rule{},
			arg_trigger: scworkflow.ADVANCE,
			want:        []string{},
		},
		{
			name: "3 scripts",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "Want that one",
				},
				{
					Trigger: "Create",
					Script:  "Want that one too",
				},
				{
					Trigger: "Close",
					Script:  "Dont want",
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[2],
				},
				{
					Trigger: "Create",
					Script:  "Last one I want",
				},
			},
			arg_trigger: scworkflow.CREATE,
			want:        []string{"Want that one", "Want that one too", "Last one I want"},
		},
		{
			name:     "0 script",
			arg_step: "ignored",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "Want that one? no",
				},
				{
					Trigger: "Create",
					Script:  "",
				},
				{
					Trigger: "Close",
					Script:  "Dont want",
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[2],
				},
				{
					Trigger: "Create",
					Script:  "Last one I don't want",
				},
			},
			arg_trigger: scworkflow.REJECT,
			want:        []string{},
		},
		{
			name: "0 script no trigger",
			rules: []*Rule{
				{
					Trigger: "Create",
					Script:  "Want that one? no",
				},
				{
					Trigger: "Create",
					Script:  "Want that one too... nah",
				},
				{
					Trigger: "Close",
					Script:  "Dont want",
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[2],
				},
				{
					Trigger: "Create",
					Script:  "Last one I don't want",
				},
			},
			arg_trigger: 0,
			want:        []string{},
		},
		{
			name: "2 scripts: advance + step",
			rules: []*Rule{
				{
					Trigger: "Advance",
					Script:  "Want that one",
					Step:    &steps[2],
				},
				{
					Trigger: "Advance",
					Script:  "Want that one too",
					Step:    &steps[2],
				},
				{
					Trigger: "Advance",
					Script:  "Dont want",
					Step:    &steps[3],
				},
				{
					Trigger: "Advance",
					Script:  "Nope",
					Step:    &steps[1],
				},
				{
					Trigger: "Create",
					Script:  "Last one I don't want",
				},
			},
			arg_trigger: scworkflow.ADVANCE,
			arg_step:    steps[2],
			want:        []string{"Want that one", "Want that one too"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WFSettings{
				Rules: tt.rules,
			}
			got := w.GetScriptsForTriggerAndStep(tt.arg_trigger, tt.arg_step)
			slices.Sort(tt.want)
			slices.Sort(got)
			if slices.Compare(got, tt.want) != 0 {
				t.Errorf("WFSettings.GetScriptsForTriggerAndStep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWFSettings_isValid(t *testing.T) {
	type fields struct {
		WFname string
		WFid   int
		Rules  RulesSlice
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "no wf name",
			fields: fields{
				WFname: "",
				WFid:   10,
				Rules:  []*Rule{},
			},
			wantErr: ErrNoWorkflowName,
		},
		{
			name: "no wf ID",
			fields: fields{
				WFname: "no ID",
				WFid:   0,
				Rules:  []*Rule{},
			},
			wantErr: ErrNoWorkflowID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WFSettings{
				WFname: tt.fields.WFname,
				WFid:   tt.fields.WFid,
				Rules:  tt.fields.Rules,
			}
			if err := w.isValid(); err != tt.wantErr {
				t.Errorf("WFSettings.isValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
