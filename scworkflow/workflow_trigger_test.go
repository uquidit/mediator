package scworkflow

import (
	"testing"
)

func TestWorkflowTrigger_Equals(t *testing.T) {
	type fields struct {
		ID       int
		Name     string
		Executer WorkflowTriggerExecuter
		Triggers []*WorkflowTriggerGroup
	}
	type args struct {
		other *WorkflowTrigger
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "same name",
			fields: fields{
				ID:       0,
				Name:     "coucou",
				Executer: WorkflowTriggerExecuter{},
				Triggers: []*WorkflowTriggerGroup{},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:       0,
					Name:     "coucou",
					Executer: WorkflowTriggerExecuter{},
					Triggers: []*WorkflowTriggerGroup{},
				},
			},
			want: true,
		},
		{
			name: "different executer path",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script1",
					Arguments: "arg",
				},
				Triggers: []*WorkflowTriggerGroup{},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script2",
						Arguments: "arg",
					},
					Triggers: []*WorkflowTriggerGroup{},
				},
			},
			want: false,
		},
		{
			name: "different executer arg",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script",
					Arguments: "arg1",
				},
				Triggers: []*WorkflowTriggerGroup{},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script",
						Arguments: "arg2",
					},
					Triggers: []*WorkflowTriggerGroup{},
				},
			},
			want: false,
		},
		{
			name: "different name but all the same",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script",
					Arguments: "arg",
				},
				Triggers: []*WorkflowTriggerGroup{
					{
						Name:     "trigger ADVANCE",
						Workflow: WorkflowTriggerWF{Name: "workflow"},
						Events:   []string{ADVANCE.Slug()},
					},
				},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script",
						Arguments: "arg",
					},
					Triggers: []*WorkflowTriggerGroup{
						{
							Name:     "trigger ADVANCE",
							Workflow: WorkflowTriggerWF{Name: "workflow"},
							Events:   []string{ADVANCE.Slug()},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "different events 1",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script",
					Arguments: "arg",
				},
				Triggers: []*WorkflowTriggerGroup{
					{
						Name:     "trigger ADVANCE1",
						Workflow: WorkflowTriggerWF{Name: "workflow"},
						Events:   []string{ADVANCE.Slug()},
					},
				},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script",
						Arguments: "arg",
					},
					Triggers: []*WorkflowTriggerGroup{
						{
							Name:     "trigger ADVANCE2",
							Workflow: WorkflowTriggerWF{Name: "workflow"},
							Events:   []string{ADVANCE.Slug(), CREATE.Slug()},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "different events 2",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script",
					Arguments: "arg",
				},
				Triggers: []*WorkflowTriggerGroup{
					{
						Name:     "trigger ADVANCE2",
						Workflow: WorkflowTriggerWF{Name: "workflow"},
						Events:   []string{ADVANCE.Slug(), CREATE.Slug()},
					},
				},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script",
						Arguments: "arg",
					},
					Triggers: []*WorkflowTriggerGroup{
						{
							Name:     "trigger ADVANCE1",
							Workflow: WorkflowTriggerWF{Name: "workflow"},
							Events:   []string{ADVANCE.Slug()},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "different events 3",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script",
					Arguments: "arg",
				},
				Triggers: []*WorkflowTriggerGroup{
					{
						Name:     "trigger ADVANCE2",
						Workflow: WorkflowTriggerWF{Name: "workflow"},
						Events:   []string{ADVANCE.Slug(), CREATE.Slug(), CLOSE.Slug(), RESUBMIT.Slug()},
					},
				},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script",
						Arguments: "arg",
					},
					Triggers: []*WorkflowTriggerGroup{
						{
							Name:     "trigger ADVANCE1",
							Workflow: WorkflowTriggerWF{Name: "workflow"},
							Events:   []string{ADVANCE.Slug(), CREATE.Slug(), REDO.Slug(), RESUBMIT.Slug()},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "different workflow name",
			fields: fields{
				ID:   0,
				Name: "coucou1",
				Executer: WorkflowTriggerExecuter{
					Type:      "",
					Path:      "/path/to/script",
					Arguments: "arg",
				},
				Triggers: []*WorkflowTriggerGroup{
					{
						Name:     "trigger ADVANCE2",
						Workflow: WorkflowTriggerWF{Name: "workflow1"},
						Events:   []string{ADVANCE.Slug()},
					},
				},
			},
			args: args{
				other: &WorkflowTrigger{
					ID:   0,
					Name: "coucou2",
					Executer: WorkflowTriggerExecuter{
						Type:      "",
						Path:      "/path/to/script",
						Arguments: "arg",
					},
					Triggers: []*WorkflowTriggerGroup{
						{
							Name:     "trigger ADVANCE1",
							Workflow: WorkflowTriggerWF{Name: "workflow2"},
							Events:   []string{ADVANCE.Slug()},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := &WorkflowTrigger{
				ID:       tt.fields.ID,
				Name:     tt.fields.Name,
				Executer: tt.fields.Executer,
				Triggers: tt.fields.Triggers,
			}
			if got := wt.Equals(tt.args.other); got != tt.want {
				t.Errorf("WorkflowTrigger.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSlicesAreTheSame(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: true,
		},
		{
			name: "nil",
			args: args{
				a: nil,
				b: nil,
			},
			want: true,
		},
		{
			name: "1 item same",
			args: args{
				a: []string{"aaa"},
				b: []string{"aaa"},
			},
			want: true,
		},
		{
			name: "1 item diff",
			args: args{
				a: []string{"aaa"},
				b: []string{"bbb"},
			},
			want: false,
		},
		{
			name: "2 items same",
			args: args{
				a: []string{"aa", "b"},
				b: []string{"aa", "b"},
			},
			want: true,
		},
		{
			name: "2 items diff",
			args: args{
				a: []string{"aa", "c"},
				b: []string{"aa", "b"},
			},
			want: false,
		},
		{
			name: "not same length",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSlicesAreTheSame(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("StringSlicesAreTheSame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowTrigger_IsTriggerRelatedToWorkflowInList(t *testing.T) {
	type fields struct {
		ID       int
		Name     string
		Executer WorkflowTriggerExecuter
		Triggers []*WorkflowTriggerGroup
	}
	type args struct {
		l []*WorkflowXML
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				Executer: WorkflowTriggerExecuter{},
				Triggers: []*WorkflowTriggerGroup{
					{
						Workflow: WorkflowTriggerWF{Name: "wf1"},
					},
				},
			},
			args: args{
				l: []*WorkflowXML{{Name: "wf1"}, {Name: "wf2"}, {Name: "wf3"}},
			},
			want: true,
		},
		{
			name: "nok",
			fields: fields{
				Executer: WorkflowTriggerExecuter{},
				Triggers: []*WorkflowTriggerGroup{
					{
						Workflow: WorkflowTriggerWF{Name: "wf4"},
					},
				},
			},
			args: args{
				l: []*WorkflowXML{{Name: "wf1"}, {Name: "wf2"}, {Name: "wf3"}},
			},
			want: false,
		},
		{
			name: "empty WF list",
			fields: fields{
				Executer: WorkflowTriggerExecuter{},
				Triggers: []*WorkflowTriggerGroup{
					{
						Workflow: WorkflowTriggerWF{Name: "wf1"},
					},
				},
			},
			args: args{
				l: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := &WorkflowTrigger{
				ID:       tt.fields.ID,
				Name:     tt.fields.Name,
				Executer: tt.fields.Executer,
				Triggers: tt.fields.Triggers,
			}
			if got := wt.IsTriggerRelatedToWorkflowInList(tt.args.l); got != tt.want {
				t.Errorf("WorkflowTrigger.IsTriggerRelatedToWorkflowInList() = %v, want %v", got, tt.want)
			}
		})
	}
}
