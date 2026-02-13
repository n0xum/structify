package parser

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	p := New()

	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestParserParseFiles(t *testing.T) {
	p := New()

	err := p.ParseFiles([]string{"../../test/fixtures/user.go"})
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	structs := p.GetStructs()

	if len(structs) == 0 {
		t.Error("ParseFiles() returned no structs")
	}
}

func TestParserParseFilesNonExistent(t *testing.T) {
	p := New()

	err := p.ParseFiles([]string{"nonexistent.go"})
	if err == nil {
		t.Error("ParseFiles() should return error for non-existent file")
	}
}

func TestParserGetStructs(t *testing.T) {
	p := New()

	_ = p.ParseFiles([]string{"../../test/fixtures/user.go"})

	structs := p.GetStructs()

	if structs == nil {
		t.Error("GetStructs() returned nil")
	}
}

func TestParseDBTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		wantValue string
	}{
		{
			name:     "pk tag",
			tag:      `db:"pk"`,
			wantValue: "pk",
		},
		{
			name:     "unique tag",
			tag:      `db:"unique"`,
			wantValue: "unique",
		},
		{
			name:     "ignore tag",
			tag:      `db:"-"`,
			wantValue: "-",
		},
		{
			name:     "table tag",
			tag:      `db:"table:custom_name"`,
			wantValue: "table:custom_name",
		},
		{
			name:     "multiple tags",
			tag:      `db:"pk,unique"`,
			wantValue: "pk,unique",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDBTag(tt.tag)
			if result != tt.wantValue {
				t.Errorf("parseDBTag() = %v, want %v", result, tt.wantValue)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single word",
			input: "User",
			want:  "user",
		},
		{
			name:  "camel case",
			input: "UserProfile",
			want:  "user_profile",
		},
		{
			name:  "pascal case - UserID",
			input: "UserID",
			want:  "user_i_d",
		},
		{
			name:  "consecutive capitals - APIKey",
			input: "APIKey",
			want:  "a_p_i_key",
		},
		{
			name:  "single capital - ID",
			input: "ID",
			want:  "i_d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.want {
				t.Errorf("toSnakeCase() = %v, want %v", result, tt.want)
			}
		})
	}
}
