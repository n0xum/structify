package adapter

import (
	"testing"

	"github.com/n0xum/structify/internal/parser"
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

func TestParserAdapterCheckConstraint(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name      string
		tag       string
		wantCheck string
	}{
		{
			name:      "simple check constraint",
			tag:       "check:age >= 18",
			wantCheck: "age >= 18",
		},
		{
			name:      "check with function",
			tag:       "check:length(name) > 0",
			wantCheck: "length(name) > 0",
		},
		{
			name:      "check with regex",
			tag:       "check:email ~* '^[a-z]+'",
			wantCheck: "email ~* '^[a-z]+'",
		},
		{
			name:      "check with other tags",
			tag:       "pk, check:id > 0",
			wantCheck: "id > 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pField := parser.Field{
				Name:        "TestField",
				Type:        "string",
				DatabaseTag: tt.tag,
			}

			domainField := adapter.toDomainField(pField)

			if domainField.CheckExpr != tt.wantCheck {
				t.Errorf("CheckExpr = %v, want %v", domainField.CheckExpr, tt.wantCheck)
			}
		})
	}
}

func TestParserAdapterDefaultConstraint(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name         string
		tag          string
		wantDefault  string
	}{
		{
			name:        "boolean default",
			tag:         "default:true",
			wantDefault: "true",
		},
		{
			name:        "string default",
			tag:         "default:'user'",
			wantDefault: "'user'",
		},
		{
			name:        "function default",
			tag:         "default:now()",
			wantDefault: "now()",
		},
		{
			name:        "complex default",
			tag:         "default:extract(epoch from now())",
			wantDefault: "extract(epoch from now())",
		},
		{
			name:        "default with other tags",
			tag:         "unique, default:0",
			wantDefault: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pField := parser.Field{
				Name:        "TestField",
				Type:        "string",
				DatabaseTag: tt.tag,
			}

			domainField := adapter.toDomainField(pField)

			if domainField.DefaultVal != tt.wantDefault {
				t.Errorf("DefaultVal = %v, want %v", domainField.DefaultVal, tt.wantDefault)
			}
		})
	}
}

func TestParserAdapterCombinedConstraints(t *testing.T) {
	adapter := NewParserAdapter()

	pField := parser.Field{
		Name:        "Age",
		Type:        "int",
		DatabaseTag: "check:age >= 18, default:18",
	}

	domainField := adapter.toDomainField(pField)

	if domainField.CheckExpr != "age >= 18" {
		t.Errorf("CheckExpr = %v, want age >= 18", domainField.CheckExpr)
	}
	if domainField.DefaultVal != "18" {
		t.Errorf("DefaultVal = %v, want 18", domainField.DefaultVal)
	}
}

func TestParserAdapterIndexConstraint(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name         string
		tag          string
		wantIndex    string
		wantUnique   bool
	}{
		{
			name:       "auto-named index",
			tag:        "index",
			wantIndex:  "test_field_idx",
			wantUnique: false,
		},
		{
			name:       "named index",
			tag:        "index:idx_email",
			wantIndex:  "idx_email",
			wantUnique: false,
		},
		{
			name:       "auto-named unique index",
			tag:        "unique_index",
			wantIndex:  "test_field_idx",
			wantUnique: true,
		},
		{
			name:       "named unique index",
			tag:        "unique_index:idx_email",
			wantIndex:  "idx_email",
			wantUnique: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pField := parser.Field{
				Name:        "TestField",
				Type:        "string",
				DatabaseTag: tt.tag,
			}

			domainField := adapter.toDomainField(pField)

			if domainField.IndexName != tt.wantIndex {
				t.Errorf("IndexName = %v, want %v", domainField.IndexName, tt.wantIndex)
			}
			if domainField.IsIndexUnique != tt.wantUnique {
				t.Errorf("IsIndexUnique = %v, want %v", domainField.IsIndexUnique, tt.wantUnique)
			}
		})
	}
}

func TestParserAdapterEnumConstraint(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name         string
		tag          string
		wantEnumVals []string
	}{
		{
			name:         "simple enum",
			tag:          "enum:active,inactive",
			wantEnumVals: []string{"active", "inactive"},
		},
		{
			name:         "enum with multiple values",
			tag:          "enum:pending,processing,shipped,delivered",
			wantEnumVals: []string{"pending", "processing", "shipped", "delivered"},
		},
		{
			name:         "enum with spaces",
			tag:          "enum:low, medium, high",
			wantEnumVals: []string{"low", "medium", "high"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pField := parser.Field{
				Name:        "TestField",
				Type:        "string",
				DatabaseTag: tt.tag,
			}

			domainField := adapter.toDomainField(pField)

			if len(domainField.EnumValues) != len(tt.wantEnumVals) {
				t.Errorf("EnumValues length = %d, want %d", len(domainField.EnumValues), len(tt.wantEnumVals))
				return
			}
			for i, val := range domainField.EnumValues {
				if val != tt.wantEnumVals[i] {
					t.Errorf("EnumValues[%d] = %v, want %v", i, val, tt.wantEnumVals[i])
				}
			}
		})
	}
}

func TestParserAdapterForeignKeyConstraint(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name         string
		tag          string
		wantTable    string
		wantColumn   string
		wantOnDelete string
		wantOnUpdate string
	}{
		{
			name:       "simple FK",
			tag:        "fk:users,id",
			wantTable:  "users",
			wantColumn: "id",
		},
		{
			name:         "FK with CASCADE",
			tag:          "fk:users,id,on_delete:CASCADE",
			wantTable:    "users",
			wantColumn:   "id",
			wantOnDelete: "CASCADE",
		},
		{
			name:         "FK with SET_NULL",
			tag:          "fk:products,id,on_delete:SET_NULL",
			wantTable:    "products",
			wantColumn:   "id",
			wantOnDelete: "SET_NULL",
		},
		{
			name:         "FK with both actions",
			tag:          "fk:orders,id,on_delete:CASCADE,on_update:CASCADE",
			wantTable:    "orders",
			wantColumn:   "id",
			wantOnDelete: "CASCADE",
			wantOnUpdate: "CASCADE",
		},
		{
			name:         "FK with NO_ACTION",
			tag:          "fk:users,id,on_delete:NO_ACTION",
			wantTable:    "users",
			wantColumn:   "id",
			wantOnDelete: "NO_ACTION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pField := parser.Field{
				Name:        "TestField",
				Type:        "int64",
				DatabaseTag: tt.tag,
			}

			domainField := adapter.toDomainField(pField)

			if domainField.FKReference == nil {
				t.Fatalf("FKReference is nil")
			}
			if domainField.FKReference.Table != tt.wantTable {
				t.Errorf("FKReference.Table = %v, want %v", domainField.FKReference.Table, tt.wantTable)
			}
			if domainField.FKReference.Column != tt.wantColumn {
				t.Errorf("FKReference.Column = %v, want %v", domainField.FKReference.Column, tt.wantColumn)
			}
			if domainField.FKOnDelete != tt.wantOnDelete {
				t.Errorf("FKOnDelete = %v, want %v", domainField.FKOnDelete, tt.wantOnDelete)
			}
			if domainField.FKOnUpdate != tt.wantOnUpdate {
				t.Errorf("FKOnUpdate = %v, want %v", domainField.FKOnUpdate, tt.wantOnUpdate)
			}
		})
	}
}

func TestParserAdapterCompositeForeignKeyConstraint(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name         string
		tag          string
		wantFKGroup  string
		wantTable    string
		wantColumn   string
		wantOnDelete string
		wantOnUpdate string
	}{
		{
			name:        "Composite FK first field",
			tag:         "fk:fk_order_item,order_items,order_id",
			wantFKGroup: "fk_order_item",
			wantTable:   "order_items",
			wantColumn:  "order_id",
		},
		{
			name:        "Composite FK second field",
			tag:         "fk:fk_order_item,order_items,item_id",
			wantFKGroup: "fk_order_item",
			wantTable:   "order_items",
			wantColumn:  "item_id",
		},
		{
			name:         "Composite FK with CASCADE",
			tag:          "fk:fk_shipment,order_items,order_id,on_delete:CASCADE,on_update:CASCADE",
			wantFKGroup:  "fk_shipment",
			wantTable:    "order_items",
			wantColumn:   "order_id",
			wantOnDelete: "CASCADE",
			wantOnUpdate: "CASCADE",
		},
		{
			name:         "Composite FK with SET NULL",
			tag:          "fk:fk_product,products,product_id,on_delete:SET_NULL",
			wantFKGroup:  "fk_product",
			wantTable:    "products",
			wantColumn:   "product_id",
			wantOnDelete: "SET_NULL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pField := parser.Field{
				Name:        "TestField",
				Type:        "int64",
				DatabaseTag: tt.tag,
			}

			domainField := adapter.toDomainField(pField)

			if domainField.FKGroup != tt.wantFKGroup {
				t.Errorf("FKGroup = %v, want %v", domainField.FKGroup, tt.wantFKGroup)
			}
			if domainField.FKReference == nil {
				t.Fatalf("FKReference is nil")
			}
			if domainField.FKReference.Table != tt.wantTable {
				t.Errorf("FKReference.Table = %v, want %v", domainField.FKReference.Table, tt.wantTable)
			}
			if domainField.FKReference.Column != tt.wantColumn {
				t.Errorf("FKReference.Column = %v, want %v", domainField.FKReference.Column, tt.wantColumn)
			}
			if domainField.FKOnDelete != tt.wantOnDelete {
				t.Errorf("FKOnDelete = %v, want %v", domainField.FKOnDelete, tt.wantOnDelete)
			}
			if domainField.FKOnUpdate != tt.wantOnUpdate {
				t.Errorf("FKOnUpdate = %v, want %v", domainField.FKOnUpdate, tt.wantOnUpdate)
			}
		})
	}
}
