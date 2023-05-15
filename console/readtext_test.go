package console

import (
	"testing"
)

func Test_isYes(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 'y'",
			args: args{v: "y"},
			want: true,
		},
		{
			name: "test 'yes'",
			args: args{v: "yes"},
			want: true,
		},
		{
			name: "test 'n'",
			args: args{v: "n"},
			want: false,
		},
		{
			name: "test 'no'",
			args: args{v: "no"},
			want: false,
		},
		{
			name: "test 'Y'",
			args: args{v: "Y"},
			want: false, //case sensitive
		},
		{
			name: "test 'YES'",
			args: args{v: "YES"},
			want: false, //case sensitive
		},
		{
			name: "test 'N'",
			args: args{v: "N"},
			want: false,
		},
		{
			name: "test 'NO'",
			args: args{v: "NO"},
			want: false,
		},
		{
			name: "test empty",
			args: args{v: ""},
			want: false,
		},
		{
			name: "test foobar",
			args: args{v: "foo"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isYes(tt.args.v); got != tt.want {
				t.Errorf("isYes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNo(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 'y'",
			args: args{v: "y"},
			want: false,
		},
		{
			name: "test 'yes'",
			args: args{v: "yes"},
			want: false,
		},
		{
			name: "test 'n'",
			args: args{v: "n"},
			want: true,
		},
		{
			name: "test 'no'",
			args: args{v: "no"},
			want: true,
		},
		{
			name: "test 'Y'",
			args: args{v: "Y"},
			want: false, //case sensitive
		},
		{
			name: "test 'YES'",
			args: args{v: "YES"},
			want: false, //case sensitive
		},
		{
			name: "test 'N'",
			args: args{v: "N"},
			want: false,
		},
		{
			name: "test 'NO'",
			args: args{v: "NO"},
			want: false,
		},
		{
			name: "test empty",
			args: args{v: ""},
			want: false,
		},
		{
			name: "test foobar",
			args: args{v: "foo"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNo(tt.args.v); got != tt.want {
				t.Errorf("isNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkBoolean(t *testing.T) {
	type args struct {
		answer string
	}
	tests := []struct {
		name               string
		args               args
		wantDefaultYES     bool
		wantDefaultNO      bool
		wantDefaultNONE    bool
		wantErrDefaultYES  bool
		wantErrDefaultNO   bool
		wantErrDefaultNONE bool
	}{
		{
			name:               "y",
			args:               args{answer: "y"},
			wantDefaultYES:     true,
			wantDefaultNO:      true,
			wantDefaultNONE:    true,
			wantErrDefaultYES:  false,
			wantErrDefaultNO:   false,
			wantErrDefaultNONE: false,
		},
		{
			name:               "yes",
			args:               args{answer: "yes"},
			wantDefaultYES:     true,
			wantDefaultNO:      true,
			wantDefaultNONE:    true,
			wantErrDefaultYES:  false,
			wantErrDefaultNO:   false,
			wantErrDefaultNONE: false,
		},
		{
			name:               "n",
			args:               args{answer: "n"},
			wantDefaultYES:     false,
			wantDefaultNO:      false,
			wantDefaultNONE:    false,
			wantErrDefaultYES:  false,
			wantErrDefaultNO:   false,
			wantErrDefaultNONE: false,
		},
		{
			name:               "no",
			args:               args{answer: "no"},
			wantDefaultYES:     false,
			wantDefaultNO:      false,
			wantDefaultNONE:    false,
			wantErrDefaultYES:  false,
			wantErrDefaultNO:   false,
			wantErrDefaultNONE: false,
		},
		{
			name:              "empty",
			args:              args{answer: ""},
			wantDefaultYES:    true,
			wantErrDefaultYES: false,

			wantDefaultNO:    false,
			wantErrDefaultNO: false,

			wantDefaultNONE:    false,
			wantErrDefaultNONE: true,
		},
		{
			name:              "foobar",
			args:              args{answer: "foobar"},
			wantDefaultYES:    true,
			wantErrDefaultYES: false,

			wantDefaultNO:    false,
			wantErrDefaultNO: false,

			wantDefaultNONE:    false,
			wantErrDefaultNONE: true,
		},
		{
			name:              "YES", // case sensitive
			args:              args{answer: "YES"},
			wantDefaultYES:    true,
			wantErrDefaultYES: false,

			wantDefaultNO:    false,
			wantErrDefaultNO: false,

			wantDefaultNONE:    false,
			wantErrDefaultNONE: true,
		},
		{
			name:              "NO", // case sensitive
			args:              args{answer: "NO"},
			wantDefaultYES:    true,
			wantErrDefaultYES: false,

			wantDefaultNO:    false,
			wantErrDefaultNO: false,

			wantDefaultNONE:    false,
			wantErrDefaultNONE: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// GetBooleanDefault_No
			got, err := checkBoolean(tt.args.answer, GetBooleanDefault_No)
			if (err != nil) != tt.wantErrDefaultNO {
				t.Errorf("checkBoolean(%s, GetBooleanDefault_No) error = %v, wantErr %v", tt.args.answer, err, tt.wantErrDefaultNO)
				return
			}
			if got != tt.wantDefaultNO {
				t.Errorf("checkBoolean(%s, GetBooleanDefault_No) = %v, want %v", tt.args.answer, got, tt.wantDefaultNO)
			}

			// GetBooleanDefault_Yes
			got, err = checkBoolean(tt.args.answer, GetBooleanDefault_Yes)
			if (err != nil) != tt.wantErrDefaultYES {
				t.Errorf("checkBoolean(%s, GetBooleanDefault_Yes) error = %v, wantErr %v", tt.args.answer, err, tt.wantErrDefaultYES)
				return
			}
			if got != tt.wantDefaultYES {
				t.Errorf("checkBoolean(%s, GetBooleanDefault_Yes) = %v, want %v", tt.args.answer, got, tt.wantDefaultYES)
			}

			// GetBooleanDefault_None
			got, err = checkBoolean(tt.args.answer, GetBooleanDefault_None)
			if (err != nil) != tt.wantErrDefaultNONE {
				t.Errorf("checkBoolean(%s, GetBooleanDefault_None) error = %v, wantErr %v", tt.args.answer, err, tt.wantErrDefaultNONE)
				return
			}
			if got != tt.wantDefaultNONE {
				t.Errorf("checkBoolean(%s, GetBooleanDefault_None) = %v, want %v", tt.args.answer, got, tt.wantDefaultNONE)
			}
		})
	}
}
