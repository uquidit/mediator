package scworkflow

import (
	"testing"
)

func TestGetTriggerFromString(t *testing.T) {
	tests := []struct {
		name string
		t    string
		want SecurechangeTrigger
	}{
		{
			name: "empty",
			t:    "",
			want: NO_TRIGGER,
		},
		{
			name: "unknown",
			t:    "unknown trigger",
			want: NO_TRIGGER,
		},
		{
			name: "Create",
			t:    "Create",
			want: CREATE,
		},
		{
			name: "Close",
			t:    "Close",
			want: CLOSE,
		},
		{
			name: "Cancel",
			t:    "Cancel",
			want: CANCEL,
		},
		{
			name: "Reject",
			t:    "Reject",
			want: REJECT,
		},
		{
			name: "Advance",
			t:    "Advance",
			want: ADVANCE,
		},
		{
			name: "Redo",
			t:    "Redo",
			want: REDO,
		},
		{
			name: "Resubmit",
			t:    "Resubmit",
			want: RESUBMIT,
		},
		{
			name: "Reopen",
			t:    "Reopen",
			want: REOPEN,
		},
		{
			name: "Resolve",
			t:    "Resolve",
			want: RESOLVE,
		},
		{
			name: "Automatic step failed",
			t:    "Automatic step failed",
			want: AUTOMATION_FAILED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTriggerFromString(tt.t); got != tt.want {
				t.Errorf("GetTriggerFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurechangeTrigger_UseNextStep(t *testing.T) {
	tests := []struct {
		name string
		tr   SecurechangeTrigger
		want bool
	}{
		{
			name: "Create",
			tr:   CREATE,
			want: false, //CREATE trigger cannot use next step: this trigger does not use steps
		},
		{
			name: "Close",
			tr:   CLOSE,
			want: false,
		},
		{
			name: "Cancel",
			tr:   CANCEL,
			want: false,
		},
		{
			name: "Reject",
			tr:   REJECT,
			want: false,
		},
		{
			name: "Advance",
			tr:   ADVANCE,
			want: true,
		},
		{
			name: "Redo",
			tr:   REDO,
			want: false,
		},
		{
			name: "Resubmit",
			tr:   RESUBMIT,
			want: false,
		},
		{
			name: "Reopen",
			tr:   REOPEN,
			want: false,
		},
		{
			name: "Resolve",
			tr:   RESOLVE,
			want: false,
		},
		{
			name: "Automatic step failed",
			tr:   AUTOMATION_FAILED,
			want: false,
		},
		{
			name: "unkwnown",
			tr:   NO_TRIGGER,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.UseNextStep(); got != tt.want {
				t.Errorf("SecurechangeTrigger.UseNextStep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurechangeTrigger_NeedStepToGetScript(t *testing.T) {
	tests := []struct {
		name string
		tr   SecurechangeTrigger
		want bool
	}{
		{
			name: "Create",
			tr:   CREATE,
			want: false,
		},
		{
			name: "Close",
			tr:   CLOSE,
			want: false,
		},
		{
			name: "Cancel",
			tr:   CANCEL,
			want: false,
		},
		{
			name: "Reject",
			tr:   REJECT,
			want: false,
		},
		{
			name: "Advance",
			tr:   ADVANCE,
			want: true,
		},
		{
			name: "Redo",
			tr:   REDO,
			want: true,
		},
		{
			name: "Resubmit",
			tr:   RESUBMIT,
			want: false,
		},
		{
			name: "Reopen",
			tr:   REOPEN,
			want: true,
		},
		{
			name: "Resolve",
			tr:   RESOLVE,
			want: false,
		},
		{
			name: "Automatic step failed",
			tr:   AUTOMATION_FAILED,
			want: true,
		},
		{
			name: "unkwnown",
			tr:   NO_TRIGGER,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.NeedStepToGetScript(); got != tt.want {
				t.Errorf("SecurechangeTrigger.NeedStepToGetScript() = %v, want %v", got, tt.want)
			}
		})
	}
}
