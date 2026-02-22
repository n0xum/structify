package entity

import (
	"testing"

	"github.com/n0xum/structify/internal/util"
)

func TestValidateFieldName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "UserID", false},
		{"empty name", "", true},
		{"name with space", "user name", true},
		{"name with tab", "user\tname", true},
		{"name with newline", "user\nname", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid type", "string", false},
		{"pointer type", "*string", false},
		{"slice type", "[]int", true},
		{"qualified type", "time.Time", false},
		{"empty type", "", true},
		{"only pointer", "*", true},
		{"only slice", "[]", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidationToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"User", "user"},
		{"UserProfile", "user_profile"},
		{"HTTPSServer", "https_server"},
		{"ID", "id"},
		{"APIKey", "api_key"},
		{"simpleword", "simpleword"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := util.ToSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
