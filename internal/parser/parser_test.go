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

func TestParserResetsBetweenCalls(t *testing.T) {
	p := New()

	err := p.ParseFiles([]string{"../../test/fixtures/user.go"})
	if err != nil {
		t.Fatalf("first ParseFiles() error = %v", err)
	}

	err = p.ParseFiles([]string{"../../test/fixtures/user.go"})
	if err != nil {
		t.Fatalf("second ParseFiles() error = %v", err)
	}

	structs := p.GetStructs()
	total := 0
	for _, s := range structs {
		total += len(s)
	}

	// Must not accumulate — same file parsed twice should yield same count
	err = New().ParseFiles([]string{"../../test/fixtures/user.go"})
	if err != nil {
		t.Fatalf("fresh ParseFiles() error = %v", err)
	}
	fresh := New()
	_ = fresh.ParseFiles([]string{"../../test/fixtures/user.go"})
	expected := 0
	for _, s := range fresh.GetStructs() {
		expected += len(s)
	}

	if total != expected {
		t.Errorf("parser accumulated state: got %d structs after 2 calls, expected %d", total, expected)
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

func TestParseDBTagComplex(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		wantValue string
	}{
		{
			name:     "check constraint with spaces",
			tag:      `db:"check:age >= 18"`,
			wantValue: "check:age >= 18",
		},
		{
			name:     "default with function",
			tag:      `db:"default:now()"`,
			wantValue: "default:now()",
		},
		{
			name:     "enum with multiple values",
			tag:      `db:"enum:pending,processing,shipped"`,
			wantValue: "enum:pending,processing,shipped",
		},
		{
			name:     "foreign key with cascade",
			tag:      `db:"fk:users,id,on_delete:CASCADE"`,
			wantValue: "fk:users,id,on_delete:CASCADE",
		},
		{
			name:     "index with name",
			tag:      `db:"index:idx_email"`,
			wantValue: "index:idx_email",
		},
		{
			name:     "unique index",
			tag:      `db:"unique_index:idx_email"`,
			wantValue: "unique_index:idx_email",
		},
		{
			name:     "unique constraint with group",
			tag:      `db:"unique:uq_user_email"`,
			wantValue: "unique:uq_user_email",
		},
		{
			name:     "composite FK",
			tag:      `db:"fk:fk_order,order_items,order_id"`,
			wantValue: "fk:fk_order,order_items,order_id",
		},
		{
			name:     "quoted string in check",
			tag:      `db:"check:name ~* '^[a-z]+'"`,
			wantValue: "check:name ~* '^[a-z]+'",
		},
		{
			name:     "check with regex",
			tag:      `db:"check:email LIKE '%@%'"`,
			wantValue: "check:email LIKE '%@%'",
		},
		{
			name:     "multiple tags with json",
			tag:      `db:"pk" json:"id"`,
			wantValue: "pk",
		},
		{
			name:     "combined constraints space separated",
			tag:      `db:"check:age >= 18" db:"default:18"`,
			wantValue: "check:age >= 18", // Returns first db: tag found
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

func TestParseTagParts(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		wantLen  int
		wantFirst string
		wantLast string
	}{
		{
			name:     "simple tag",
			tag:      "pk",
			wantLen:  1,
			wantFirst: "pk",
			wantLast: "pk",
		},
		{
			name:     "two tags space separated",
			tag:      "pk unique",
			wantLen:  2,
			wantFirst: "pk",
			wantLast: "unique",
		},
		{
			name:     "tag with value containing comma",
			tag:      "enum:a,b,c",
			wantLen:  1,
			wantFirst: "enum:a,b,c",
			wantLast: "enum:a,b,c",
		},
		{
			name:     "tag with function call",
			tag:      "default:now()",
			wantLen:  1,
			wantFirst: "default:now()",
			wantLast: "default:now()",
		},
		{
			name:     "multiple simple tags",
			tag:      "pk unique index",
			wantLen:  3,
			wantFirst: "pk",
			wantLast: "index",
		},
		{
			name:     "tag with quotes preserves content",
			tag:      `db:"pk"`,
			wantLen:  1,
			wantFirst: `db:"pk"`,
			wantLast: `db:"pk"`,
		},
		{
			name:     "multiple tags with quotes",
			tag:      `db:"pk" json:"id"`,
			wantLen:  2,
			wantFirst: `db:"pk"`,
			wantLast: `json:"id"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTagParts(tt.tag)
			if len(result) != tt.wantLen {
				t.Errorf("parseTagParts() length = %d, want %d", len(result), tt.wantLen)
			}
			if len(result) > 0 {
				if result[0] != tt.wantFirst {
					t.Errorf("parseTagParts()[0] = %v, want %v", result[0], tt.wantFirst)
				}
				if result[len(result)-1] != tt.wantLast {
					t.Errorf("parseTagParts()[-1] = %v, want %v", result[len(result)-1], tt.wantLast)
				}
			}
		})
	}
}

func TestParseFilesMultiple(t *testing.T) {
	p := New()

	files := []string{
		"../../test/fixtures/user.go",
		"../../test/fixtures/constraints.go",
	}

	err := p.ParseFiles(files)
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	structs := p.GetStructs()
	if len(structs) == 0 {
		t.Error("ParseFiles() returned no structs")
	}

	totalStructs := 0
	for _, fileStructs := range structs {
		totalStructs += len(fileStructs)
	}

	if totalStructs == 0 {
		t.Error("ParseFiles() returned 0 total structs")
	}
}

func TestParseFilesEmptySlice(t *testing.T) {
	p := New()

	err := p.ParseFiles([]string{})
	if err != nil {
		t.Fatalf("ParseFiles() with empty slice should not error, got %v", err)
	}

	structs := p.GetStructs()
	if len(structs) != 0 {
		t.Errorf("ParseFiles() with empty slice should return no structs, got %d", len(structs))
	}
}
