package main

import (
	"reflect"
	"testing"
	"uqtu/mediator/scworkflow"
)

func Test_checkArgumentsAndGetTrigger(t *testing.T) {
	tests := []struct {
		name    string
		args    arguments
		want    scworkflow.SecurechangeTrigger
		wantErr bool
	}{
		{
			name: "nothing",
			args: arguments{
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "one arg, no flag",
			args: arguments{
				positional:        []string{"one-arg"},
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "4 args, no flag",
			args: arguments{
				positional:        []string{"arg1", "arg2", "arg3", "arg4"},
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "one arg, scriptedCondition",
			args: arguments{
				positional:        []string{"arg1"},
				data_filename:     "",
				scriptedCondition: true,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: false,
		},
		{
			name: "scriptedCondition and preAssignment",
			args: arguments{
				positional:        []string{"arg1"},
				data_filename:     "",
				scriptedCondition: true,
				preAssignment:     true,
				scriptedTask:      false,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "scriptedCondition and preAssignment and scriptedTask",
			args: arguments{
				positional:        []string{"arg1"},
				data_filename:     "",
				scriptedCondition: true,
				preAssignment:     true,
				scriptedTask:      true,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "preAssignment and scriptedTask",
			args: arguments{
				positional:        []string{"arg1"},
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     true,
				scriptedTask:      true,
				trigger:           "",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "unknown trigger",
			args: arguments{
				positional:        []string{"arg1"},
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "coucou",
			},
			want:    scworkflow.NO_TRIGGER,
			wantErr: true,
		},
		{
			name: "known trigger, one arg",
			args: arguments{
				positional:        []string{"arg1"},
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "advance",
			},
			want:    scworkflow.ADVANCE,
			wantErr: false,
		},
		{
			name: "known trigger, 4 args",
			args: arguments{
				positional:        []string{"arg1", "arg2", "arg3", "arg4"},
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "close",
			},
			want:    scworkflow.CLOSE,
			wantErr: true,
		},
		{
			name: "known trigger, 0 arg",
			args: arguments{
				data_filename:     "",
				scriptedCondition: false,
				preAssignment:     false,
				scriptedTask:      false,
				trigger:           "Create",
			},
			want:    scworkflow.CREATE,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkArgumentsAndGetTrigger(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkArgumentsAndGetTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkArgumentsAndGetTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}
