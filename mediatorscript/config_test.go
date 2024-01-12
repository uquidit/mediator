package mediatorscript

/*
func TestWorkflow_GetNextStep(t *testing.T) {
	type fields struct {
		Steps []Steps
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil step array",
			fields: fields{
				Steps: nil,
			},
			args: args{
				s: "step",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty step array",
			fields: fields{
				Steps: []Steps{},
			},
			args: args{
				s: "step",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "one step good",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step",
						Script: "a_script.sh",
					},
				},
			},
			args: args{
				s: "step",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "one step bad",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "stepbad",
						Script: "a_script.sh",
					},
				},
			},
			args: args{
				s: "step",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "two steps - first good",
			fields: fields{
				Steps: []Steps{
					{
						Name: "step1",
					},
					{
						Name: "step2",
					},
				},
			},
			args: args{
				s: "step1",
			},
			want:    "step2",
			wantErr: false,
		},
		{
			name: "two steps - last good",
			fields: fields{
				Steps: []Steps{
					{
						Name: "step1",
					},
					{
						Name: "step2",
					},
				},
			},
			args: args{
				s: "step2",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "3 steps - first good",
			fields: fields{
				Steps: []Steps{
					{
						Name: "step1",
					},
					{
						Name: "step2",
					},
					{
						Name: "step3",
					},
				},
			},
			args: args{
				s: "step1",
			},
			want:    "step2",
			wantErr: false,
		},
		{
			name: "3 steps - second good",
			fields: fields{
				Steps: []Steps{
					{
						Name: "step1",
					},
					{
						Name: "step2",
					},
					{
						Name: "step3",
					},
				},
			},
			args: args{
				s: "step2",
			},
			want:    "step3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Workflow{
				Steps: tt.fields.Steps,
			}
			got, err := w.GetNextStep(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Workflow.GetNextStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Workflow.GetNextStep() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

/*func TestWorkflow_GetScriptForStep(t *testing.T) {
	type fields struct {
		Steps []Steps
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "1 step with script ok",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step",
						Script: "script",
					},
				},
			},
			args: args{
				s: "step",
			},
			want: "script",
		},
		{
			name: "1 step with script nok",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step",
						Script: "script",
					},
				},
			},
			args: args{
				s: "step2",
			},
			want: "",
		},
		{
			name: "2 steps with script ok",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step1",
						Script: "script1",
					},
					{
						Name:   "step2",
						Script: "script2",
					},
				},
			},
			args: args{
				s: "step1",
			},
			want: "script1",
		},
		{
			name: "2 steps with script nok",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step1",
						Script: "script1",
					},
					{
						Name:   "step2",
						Script: "script2",
					},
				},
			},
			args: args{
				s: "step3",
			},
			want: "",
		},
		{
			name: "2 steps with script ok",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step1",
						Script: "script1",
					},
					{
						Name:   "step2",
						Script: "script2",
					},
				},
			},
			args: args{
				s: "step2",
			},
			want: "script2",
		},
		{
			name: "empty step",
			fields: fields{
				Steps: []Steps{},
			},
			args: args{
				s: "s",
			},
			want: "",
		},
		{
			name: "empty input step",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "step",
						Script: "script",
					},
				},
			},
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "empty input step",
			fields: fields{
				Steps: []Steps{
					{
						Name:   "",
						Script: "scriptempty",
					},
					{
						Name:   "step",
						Script: "script",
					},
				},
			},
			args: args{
				s: "",
			},
			want: "scriptempty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Workflow{
				Steps: tt.fields.Steps,
			}
			if got := w.GetScriptForTriggerAndStep(tt.args.s); got != tt.want {
				t.Errorf("Workflow.GetScriptForStep() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
