package mapper

import (
	"testing"
)

func TestMapperMapType(t *testing.T) {
	mapper := NewMapper()

	tests := []struct {
		name        string
		goType      string
		wantType    string
		wantNotNull bool
	}{
		{
			name:        "int64",
			goType:      "int64",
			wantType:    "BIGINT",
			wantNotNull: true,
		},
		{
			name:        "string",
			goType:      "string",
			wantType:    "VARCHAR(255)",
			wantNotNull: false,
		},
		{
			name:        "bool",
			goType:      "bool",
			wantType:    "BOOLEAN",
			wantNotNull: true,
		},
		{
			name:        "float64",
			goType:      "float64",
			wantType:    "DOUBLE PRECISION",
			wantNotNull: true,
		},
		{
			name:        "time.Time maps to TIMESTAMP",
			goType:      "time.Time",
			wantType:    "TIMESTAMP",
			wantNotNull: true,
		},
		{
			name:        "json.RawMessage maps to JSONB",
			goType:      "json.RawMessage",
			wantType:    "JSONB",
			wantNotNull: false,
		},
		{
			name:        "pointer to time.Time",
			goType:      "*time.Time",
			wantType:    "TIMESTAMP",
			wantNotNull: true,
		},
		{
			name:        "unknown type",
			goType:      "CustomType",
			wantType:    "TEXT",
			wantNotNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := mapper.MapType(tt.goType)
			if mapping.PostgresType != tt.wantType {
				t.Errorf("PostgresType = %v, want %v", mapping.PostgresType, tt.wantType)
			}
			if mapping.IsNotNull != tt.wantNotNull {
				t.Errorf("IsNotNull = %v, want %v", mapping.IsNotNull, tt.wantNotNull)
			}
		})
	}
}

func TestMapperFormatColumnDefinition(t *testing.T) {
	mapper := NewMapper()

	tests := []struct {
		name    string
		field   string
		mapping TypeMapping
		tags    []string
		want    string
	}{
		{
			name:    "primary key",
			field:   "id",
			mapping: TypeMapping{PostgresType: "BIGINT", IsNotNull: true},
			tags:    []string{"pk"},
			want:    "BIGINT PRIMARY KEY",
		},
		{
			name:    "unique field",
			field:   "email",
			mapping: TypeMapping{PostgresType: "VARCHAR(255)", IsNotNull: false},
			tags:    []string{"unique"},
			want:    "VARCHAR(255) UNIQUE",
		},
		{
			name:    "ignored field",
			field:   "password",
			mapping: TypeMapping{PostgresType: "VARCHAR(255)", IsNotNull: true},
			tags:    []string{"-"},
			want:    "",
		},
		{
			name:    "not null without pk",
			field:   "name",
			mapping: TypeMapping{PostgresType: "VARCHAR(255)", IsNotNull: true},
			tags:    []string{},
			want:    "VARCHAR(255) NOT NULL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.FormatColumnDefinition(tt.field, tt.mapping, tt.tags)
			if got != tt.want {
				t.Errorf("FormatColumnDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapperHasTag(t *testing.T) {
	mapper := NewMapper()

	tests := []struct {
		name string
		tags []string
		tag  string
		want bool
	}{
		{
			name: "tag exists",
			tags: []string{"pk", "unique"},
			tag:  "pk",
			want: true,
		},
		{
			name: "tag does not exist",
			tags: []string{"pk", "unique"},
			tag:  "index",
			want: false,
		},
		{
			name: "empty tags",
			tags: []string{},
			tag:  "pk",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.HasTag(tt.tags, tt.tag)
			if got != tt.want {
				t.Errorf("HasTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapperGetBaseType(t *testing.T) {
	mapper := NewMapper()

	tests := []struct {
		name     string
		input    string
		wantType string
	}{
		{
			name:     "pointer type",
			input:    "*string",
			wantType: "string",
		},
		{
			name:     "slice type returns int",
			input:    "[]int",
			wantType: "int",
		},
		{
			name:     "map type returns map before split",
			input:    "map[string]int",
			wantType: "map",
		},
		{
			name:     "qualified type with mapping returns full name",
			input:    "time.Time",
			wantType: "time.Time",
		},
		{
			name:     "qualified type without mapping strips package",
			input:    "some.CustomType",
			wantType: "CustomType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.getBaseType(tt.input)
			if got != tt.wantType {
				t.Errorf("getBaseType() = %v, want %v", got, tt.wantType)
			}
		})
	}
}
