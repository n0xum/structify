package adapter

import (
	"testing"

	"github.com/ak/structify/internal/parser"
)

func TestParserAdapterToDomain(t *testing.T) {
	adapter := NewParserAdapter()

	pStruct := &parser.Struct{
		Name:       "User",
		TableName:  "users",
		PackageName: "models",
		Fields: []parser.Field{
			{Name: "ID", Type: "int64", DatabaseTag: "pk"},
			{Name: "Name", Type: "string", DatabaseTag: "unique"},
			{Name: "Email", Type: "string", DatabaseTag: ""},
			{Name: "Password", Type: "string", DatabaseTag: "-"},
		},
	}

	domainEntity := adapter.ToDomain(pStruct)

	if domainEntity.Name != "User" {
		t.Errorf("Name = %v, want User", domainEntity.Name)
	}
	if domainEntity.TableName != "users" {
		t.Errorf("TableName = %v, want users", domainEntity.TableName)
	}
	if domainEntity.Package != "models" {
		t.Errorf("Package = %v, want models", domainEntity.Package)
	}

	if len(domainEntity.Fields) != 4 {
		t.Fatalf("Fields length = %d, want 4", len(domainEntity.Fields))
	}

	idField := domainEntity.Fields[0]
	if idField.Name != "ID" {
		t.Errorf("Field[0].Name = %v, want ID", idField.Name)
	}
	if !idField.IsPrimary {
		t.Error("Field[0].IsPrimary = false, want true")
	}

	passwordField := domainEntity.Fields[3]
	if !passwordField.IsIgnored {
		t.Error("PasswordField.IsIgnored = false, want true")
	}

	generateableFields := domainEntity.GetGenerateableFields()
	if len(generateableFields) != 3 {
		t.Errorf("GetGenerateableFields() length = %d, want 3", len(generateableFields))
	}
}

func TestParserAdapterToDomainSlice(t *testing.T) {
	adapter := NewParserAdapter()

	pStructs := []*parser.Struct{
		{
			Name: "User",
			Fields: []parser.Field{
				{Name: "ID", Type: "int64", DatabaseTag: "pk"},
			},
		},
		{
			Name: "Product",
			Fields: []parser.Field{
				{Name: "ID", Type: "int64", DatabaseTag: "pk"},
			},
		},
	}

	domainEntities := adapter.ToDomainSlice(pStructs)

	if len(domainEntities) != 2 {
		t.Errorf("ToDomainSlice() returned %d entities, want 2", len(domainEntities))
	}
}

func TestParserAdapterParseTags(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name string
		tag  string
		want []string
	}{
		{
			name: "empty tag",
			tag:  "",
			want: nil,
		},
		{
			name: "single tag",
			tag:  "pk",
			want: []string{"pk"},
		},
		{
			name: "multiple tags",
			tag:  "pk, unique",
			want: []string{"pk", "unique"},
		},
		{
			name: "tags with spaces",
			tag:  "pk, unique, check:value > 0",
			want: []string{"pk", "unique", "check:value > 0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.parseTags(tt.tag)
			if len(got) != len(tt.want) {
				t.Errorf("parseTags() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseTags()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParserAdapterExtractCustomTableName(t *testing.T) {
	adapter := NewParserAdapter()

	fields := []parser.Field{
		{Name: "ID", Type: "int64", DatabaseTag: "table:custom_users"},
	}

	tableName := adapter.extractCustomTableName(fields)
	if tableName != "custom_users" {
		t.Errorf("extractCustomTableName() = %v, want custom_users", tableName)
	}
}

func TestParserAdapterToMap(t *testing.T) {
	adapter := NewParserAdapter()

	structs := map[string][]*parser.Struct{
		"models": {
			{Name: "User", Fields: []parser.Field{{Name: "ID", Type: "int64"}}},
		},
	}

	result := adapter.ToMap(structs)

	if len(result) != 1 {
		t.Errorf("ToMap() returned %d packages, want 1", len(result))
	}

	models, ok := result["models"]
	if !ok {
		t.Error("ToMap() missing models package")
		return
	}

	if len(models) != 1 {
		t.Errorf("ToMap()[models] length = %d, want 1", len(models))
	}
}
