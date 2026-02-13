package cli

import (
	"testing"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd == nil {
		t.Fatal("NewCommand() returned nil")
	}

	if cmd.FS == nil {
		t.Error("FS is nil")
	}
}

func TestCommandParse(t *testing.T) {
	cmd := NewCommand()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "valid args",
			args: []string{"cmd", "--to-sql", "file.go"},
		},
		{
			name: "no args",
			args: []string{"cmd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = cmd.Parse(tt.args)
		})
	}
}

func TestCommandValidate(t *testing.T) {
	tests := []struct {
		name   string
		toSQL  bool
		toDBCode bool
		fromJSON bool
		jsonFile string
		wantErr bool
	}{
		{
			name:   "all flags false",
			wantErr: true,
		},
		{
			name:   "to SQL only",
			toSQL:  true,
			wantErr: false,
		},
		{
			name:   "to DB code only",
			toDBCode: true,
			wantErr: false,
		},
		{
			name:   "from JSON without file",
			fromJSON: true,
			jsonFile:  "",
			wantErr: true,
		},
		{
			name:   "from JSON with file",
			fromJSON: true,
			jsonFile:  "test.json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			cmd.ToSQL = tt.toSQL
			cmd.ToDBCode = tt.toDBCode
			cmd.FromJSON = tt.fromJSON
			cmd.JSONFile = tt.jsonFile

			err := cmd.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}
