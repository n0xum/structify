package mapper

import (
	"strings"
	"testing"
)

func TestMapperParseConstraints(t *testing.T) {
	m := NewMapper()

	tests := []struct {
		name   string
		tags   []string
		wantN  int
		wantIn string
	}{
		{
			name:  "no constraints",
			tags:  []string{"pk"},
			wantN: 0,
		},
		{
			name:   "check constraint",
			tags:   []string{"check:age > 0"},
			wantN:  1,
			wantIn: "CHECK (age > 0)",
		},
		{
			name:   "default constraint",
			tags:   []string{"default:now()"},
			wantN:  1,
			wantIn: "DEFAULT now()",
		},
		{
			name:   "enum constraint",
			tags:   []string{"enum:'a','b','c'"},
			wantN:  1,
			wantIn: "CHECK (column_name IN ('a','b','c'))",
		},
		{
			name:   "multiple constraints",
			tags:   []string{"check:x > 0", "default:0"},
			wantN:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.parseConstraints(tt.tags)
			if len(got) != tt.wantN {
				t.Errorf("parseConstraints() len = %d, want %d", len(got), tt.wantN)
			}
			if tt.wantIn != "" {
				found := false
				for _, c := range got {
					if strings.Contains(c, tt.wantIn) || c == tt.wantIn {
						found = true
					}
				}
				if !found {
					t.Errorf("parseConstraints() missing %q in %v", tt.wantIn, got)
				}
			}
		})
	}
}

func TestMapperToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"User", "user"},
		{"UserProfile", "user_profile"},
		{"ID", "i_d"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMapperMapTypeExtended(t *testing.T) {
	m := NewMapper()

	tests := []struct {
		goType   string
		wantType string
	}{
		{"int", "INTEGER"},
		{"int8", "SMALLINT"},
		{"int16", "SMALLINT"},
		{"int32", "INTEGER"},
		{"uint", "BIGINT"},
		{"uint8", "SMALLINT"},
		{"uint16", "SMALLINT"},
		{"uint32", "INTEGER"},
		{"uint64", "BIGINT"},
		{"float32", "REAL"},
		{"*string", "VARCHAR(255)"},
		{"*int64", "BIGINT"},
	}

	for _, tt := range tests {
		t.Run(tt.goType, func(t *testing.T) {
			got := m.MapType(tt.goType)
			if got.PostgresType != tt.wantType {
				t.Errorf("MapType(%q).PostgresType = %q, want %q", tt.goType, got.PostgresType, tt.wantType)
			}
		})
	}
}
