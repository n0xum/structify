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
		name          string
		toSQL         bool
		toRepo        bool
		modelFile     string
		interfaceFile string
		wantErr       bool
	}{
		{
			name:    "all flags false",
			wantErr: true,
		},
		{
			name:    "to SQL only",
			toSQL:   true,
			wantErr: false,
		},
		{
			name:    "to-repo without model",
			toRepo:  true,
			wantErr: true,
		},
		{
			name:      "to-repo without interface",
			toRepo:    true,
			modelFile: "model.go",
			wantErr:   true,
		},
		{
			name:          "to-repo with both",
			toRepo:        true,
			modelFile:     "model.go",
			interfaceFile: "repo.go",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			cmd.ToSQL = tt.toSQL
			cmd.ToRepo = tt.toRepo
			cmd.ModelFile = tt.modelFile
			cmd.InterfaceFile = tt.interfaceFile

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
