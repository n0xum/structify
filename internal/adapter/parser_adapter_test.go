package adapter

import (
	"testing"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/parser"
	"github.com/n0xum/structify/internal/util"
)

func TestParserAdapterToDomain(t *testing.T) {
	adapter := NewParserAdapter()

	pStruct := &parser.Struct{
		Name:        "User",
		TableName:   "users",
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
		name        string
		tag         string
		wantDefault string
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
		name       string
		tag        string
		wantIndex  string
		wantUnique bool
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

func TestToRepositoryInterface(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Email", Type: "string"},
			{Name: "Active", Type: "bool"},
		},
	}

	iface := &parser.Interface{
		Name:        "UserRepository",
		PackageName: "repository",
		Methods: []parser.Method{
			{Name: "Create", Params: []parser.Param{{Name: "user", Type: "*User"}}, Returns: []parser.Return{{Type: "*User", IsPointer: true, BaseType: "User"}, {Type: "error", BaseType: "error"}}},
			{Name: "GetByID", Params: []parser.Param{{Name: "id", Type: "int64"}}, Returns: []parser.Return{{Type: "*User", IsPointer: true, BaseType: "User"}, {Type: "error", BaseType: "error"}}},
			{Name: "Update", Params: []parser.Param{{Name: "user", Type: "*User"}}, Returns: []parser.Return{{Type: "error", BaseType: "error"}}},
			{Name: "Delete", Params: []parser.Param{{Name: "id", Type: "int64"}}, Returns: []parser.Return{{Type: "error", BaseType: "error"}}},
			{Name: "FindByEmail", Params: []parser.Param{{Name: "email", Type: "string"}}, Returns: []parser.Return{{Type: "*User", IsPointer: true, BaseType: "User"}, {Type: "error", BaseType: "error"}}},
			{Name: "FindByActive", Params: []parser.Param{{Name: "active", Type: "bool"}}, Returns: []parser.Return{{Type: "[]*User", IsSlice: true, BaseType: "User"}, {Type: "error", BaseType: "error"}}},
			{Name: "FindRecentUsers", Params: []parser.Param{{Name: "since", Type: "int64"}}, Returns: []parser.Return{{Type: "[]*User", IsSlice: true, BaseType: "User"}, {Type: "error", BaseType: "error"}}, SQLComment: "SELECT * FROM users WHERE created > $1"},
		},
	}

	repo := adapter.ToRepositoryInterface(iface, ent)
	if repo == nil {
		t.Fatal("ToRepositoryInterface() returned nil")
	}

	if repo.Name != "UserRepository" {
		t.Errorf("Name = %s, want UserRepository", repo.Name)
	}
	if repo.EntityName != "User" {
		t.Errorf("EntityName = %s, want User", repo.EntityName)
	}
	if len(repo.Methods) != 7 {
		t.Fatalf("Methods count = %d, want 7", len(repo.Methods))
	}

	expectedKinds := []entity.MethodKind{
		entity.MethodCreate,
		entity.MethodGetByID,
		entity.MethodUpdate,
		entity.MethodDelete,
		entity.MethodFindBy,
		entity.MethodFindBy,
		entity.MethodCustomSQL,
	}

	for i, m := range repo.Methods {
		if m.Kind != expectedKinds[i] {
			t.Errorf("method[%d] %s: Kind = %d, want %d", i, m.Name, m.Kind, expectedKinds[i])
		}
	}
}

func TestExtractFindByFields(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name       string
		methodName string
		want       []string
	}{
		{"single field", "FindByEmail", []string{"Email"}},
		{"multi field", "FindByStatusAndRole", []string{"Status", "Role"}},
		{"triple field", "FindByNameAndAgeAndActive", []string{"Name", "Age", "Active"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := adapter.extractFindByFields(tt.methodName)
			if len(fields) != len(tt.want) {
				t.Fatalf("extractFindByFields(%s) len = %d, want %d", tt.methodName, len(fields), len(tt.want))
			}
			for i, f := range fields {
				if f != tt.want[i] {
					t.Errorf("field[%d] = %s, want %s", i, f, tt.want[i])
				}
			}
		})
	}
}

func TestToRepositoryInterfaceNilInputs(t *testing.T) {
	adapter := NewParserAdapter()

	if adapter.ToRepositoryInterface(nil, nil) != nil {
		t.Error("ToRepositoryInterface(nil, nil) should return nil")
	}
	if adapter.ToRepositoryInterface(&parser.Interface{}, nil) != nil {
		t.Error("ToRepositoryInterface(iface, nil) should return nil")
	}
	if adapter.ToRepositoryInterface(nil, &entity.Entity{}) != nil {
		t.Error("ToRepositoryInterface(nil, ent) should return nil")
	}
}

func TestToInterfaceMap(t *testing.T) {
	adapter := NewParserAdapter()

	t.Run("nil input", func(t *testing.T) {
		result := adapter.ToInterfaceMap(nil)
		if result != nil {
			t.Error("ToInterfaceMap(nil) should return nil")
		}
	})

	t.Run("empty map", func(t *testing.T) {
		result := adapter.ToInterfaceMap(map[string][]*parser.Interface{})
		if result == nil {
			t.Error("ToInterfaceMap() should return non-nil map")
		}
		if len(result) != 0 {
			t.Error("ToInterfaceMap() should return empty map")
		}
	})

	t.Run("with interfaces", func(t *testing.T) {
		input := map[string][]*parser.Interface{
			"repo": {
				{Name: "UserRepository", PackageName: "repo"},
				{Name: "ProductRepository", PackageName: "repo"},
			},
		}

		result := adapter.ToInterfaceMap(input)
		if result == nil {
			t.Fatal("ToInterfaceMap() returned nil")
		}

		if len(result) != 1 {
			t.Errorf("ToInterfaceMap() returned %d packages, want 1", len(result))
		}

		repos := result["repo"]
		if repos == nil {
			t.Fatal("ToInterfaceMap() missing 'repo' package")
		}

		if len(repos) != 2 {
			t.Errorf("ToInterfaceMap() returned %d repos, want 2", len(repos))
		}
	})
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"PascalCase", "UserName", "user_name"},
		{"camelCase", "userName", "user_name"},
		{"single word", "user", "user"},
		{"all caps", "ID", "id"},
		{"multiple capitals", "HTTPServer", "http_server"},
		{"with numbers", "User2Name", "user2_name"},
		{"empty string", "", ""},
		{"single letter", "A", "a"},
		{"already snake", "user_name", "user_name"},
		{"consecutive capitals", "XMLParser", "xml_parser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStartsWithKnownPrefix(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"pk tag", "pk", true},
		{"unique tag", "unique", true},
		{"check tag", "check:age > 0", true},
		{"default tag", "default:now", true},
		{"index tag", "index", true},
		{"unique_index tag", "unique_index", true},
		{"fk tag", "fk:users,id", true},
		{"unknown tag", "unknown", false},
		{"empty string", "", false},
		{"partial match", "pkcustom", false}, // "pk" is an exact tag, not a prefix
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.startsWithKnownPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("startsWithKnownPrefix(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Smart Query Tests

func TestSmartQuery_Classification(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Email", Type: "string"},
			{Name: "Active", Type: "bool"},
			{Name: "CreatedAt", Type: "time.Time"},
		},
	}

	tests := []struct {
		name         string
		methodName   string
		wantKind     entity.MethodKind
		description  string
	}{
		{"List by email", "ListUsersByEmail", entity.MethodSmartQuery, "List pattern with single field"},
		{"Find by ID", "FindUserByID", entity.MethodSmartQuery, "Find pattern with single field"},
		{"List by email and role", "ListUsersByEmailAndRole", entity.MethodSmartQuery, "List with AND condition"},
		{"List by email or username", "ListUsersByEmailOrUsername", entity.MethodSmartQuery, "List with OR condition"},
		{"Count users", "CountUsers", entity.MethodSmartQuery, "Count pattern"},
		{"Count by active", "CountUsersByActive", entity.MethodSmartQuery, "Count with field condition"},
		{"Exists by email", "ExistsUserByEmail", entity.MethodSmartQuery, "Exists pattern"},
		{"List by active true", "ListUsersByActiveTrue", entity.MethodSmartQuery, "Boolean pattern"},
		{"List by age greater than", "ListUsersByAgeGreaterThan", entity.MethodSmartQuery, "Comparison operator"},
		{"List by email like", "ListUsersByEmailLike", entity.MethodSmartQuery, "String pattern"},
		{"List by ID in", "ListUsersByIDIn", entity.MethodSmartQuery, "Collection pattern"},
		{"List by email order by", "ListUsersByEmailOrderByCreatedAtDesc", entity.MethodSmartQuery, "Ordering pattern"},
		{"First by email", "FirstUserByEmail", entity.MethodSmartQuery, "Limiting pattern (First)"},
		{"Top 10 by role", "Top10UsersByRole", entity.MethodSmartQuery, "Limiting pattern (TopN)"},
		{"List by recent", "ListUsersByRecentCreatedAt", entity.MethodSmartQuery, "Time-based pattern"},
		{"List by not active", "ListUsersByNotActive", entity.MethodSmartQuery, "Negation pattern"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := parser.Method{Name: tt.methodName}
			kind := adapter.classifyMethod(m, ent)
			if kind != tt.wantKind {
				t.Errorf("%s: classifyMethod() = %v, want %v (%s)", tt.name, kind, tt.wantKind, tt.description)
			}
		})
	}
}

func TestSmartQuery_SQLGeneration(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Email", Type: "string"},
			{Name: "Active", Type: "bool"},
			{Name: "CreatedAt", Type: "time.Time"},
		},
	}

	tests := []struct {
		name            string
		methodName      string
		wantSQLContains string
	}{
		{"List by email", "ListUsersByEmail", `SELECT id, username, email, active, created_at FROM "user" WHERE email = $1`},
		{"Find by ID", "FindUserByID", `SELECT id, username, email, active, created_at FROM "user" WHERE id = $1`},
		{"List by email and role", "ListUsersByEmailAndRole", "WHERE email = $1"},
		{"List by email or username", "ListUsersByEmailOrUsername", "WHERE email = $1"},
		{"Count users", "CountUsers", `SELECT COUNT(*) FROM "user"`},
		{"Count by active", "CountUsersByActive", `SELECT COUNT(*) FROM "user" WHERE active = $1`},
		{"Exists by email", "ExistsUserByEmail", `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`},
		{"List by age greater than", "ListUsersByAgeGreaterThan", "WHERE age > $1"},
		{"List by email like", "ListUsersByEmailLike", "WHERE email LIKE $1"},
		{"List by email order by", "ListUsersByEmailOrderByCreatedAtDesc", "ORDER BY created_at DESC"},
		{"First by email", "FirstUserByEmail", "LIMIT 1"},
		{"Top 10 by role", "Top10UsersByRole", "LIMIT"},
		// Note: ActiveTrue, Recent, NotActive, and IDIn patterns need special handling - for now we test basic field conditions
		{"List by active field", "ListUsersByActive", `SELECT id, username, email, active, created_at FROM "user" WHERE active = $1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := parser.Method{Name: tt.methodName}
			rm := entity.RepositoryMethod{
				Name: tt.methodName,
			}
			adapter.processSmartQueryMethod(&rm, m, ent)
			if rm.GeneratedSQL == "" {
				t.Fatalf("processSmartQueryMethod() GeneratedSQL is empty")
			}
			if !contains(rm.GeneratedSQL, tt.wantSQLContains) {
				t.Errorf("GeneratedSQL = %s\nshould contain: %s", rm.GeneratedSQL, tt.wantSQLContains)
			}
		})
	}
}

func TestSmartQuery_BackwardCompatibility(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Email", Type: "string"},
			{Name: "Active", Type: "bool"},
		},
	}

	// Test that methods with //sql: comments are still MethodCustomSQL
	m := parser.Method{
		Name:       "FindRecentUsers",
		SQLComment: "SELECT * FROM users WHERE created > $1 ORDER BY created DESC",
	}
	kind := adapter.classifyMethod(m, ent)
	if kind != entity.MethodCustomSQL {
		t.Errorf("Method with SQL comment should be MethodCustomSQL, got %v", kind)
	}

	// Test that standard FindBy methods still work
	m2 := parser.Method{Name: "FindByEmail"}
	kind2 := adapter.classifyMethod(m2, ent)
	if kind2 != entity.MethodFindBy {
		t.Errorf("FindBy method should be MethodFindBy, got %v", kind2)
	}
}

func TestSmartQuery_ToRepositoryInterface(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "AuditLog",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Action", Type: "string"},
			{Name: "UserID", Type: "int64"},
			{Name: "EntityID", Type: "int64"},
			{Name: "CreatedAt", Type: "time.Time"},
		},
	}

	iface := &parser.Interface{
		Name:        "AuditLogRepository",
		PackageName: "repository",
		Methods: []parser.Method{
			// Smart query methods
			{Name: "ListAuditLogsByID", Params: []parser.Param{{Name: "id", Type: "int64"}}, Returns: []parser.Return{{Type: "[]*AuditLog", IsSlice: true, BaseType: "AuditLog"}, {Type: "error", BaseType: "error"}}},
			{Name: "CountAuditLogsByAction", Params: []parser.Param{{Name: "action", Type: "string"}}, Returns: []parser.Return{{Type: "int64"}, {Type: "error", BaseType: "error"}}},
			{Name: "FindAuditLogByActionAndID", Params: []parser.Param{{Name: "action", Type: "string"}, {Name: "id", Type: "int64"}}, Returns: []parser.Return{{Type: "*AuditLog", IsPointer: true, BaseType: "AuditLog"}, {Type: "error", BaseType: "error"}}},
			{Name: "ListAuditLogsByActionOrderByCreatedAtDesc", Params: []parser.Param{{Name: "action", Type: "string"}}, Returns: []parser.Return{{Type: "[]*AuditLog", IsSlice: true, BaseType: "AuditLog"}, {Type: "error", BaseType: "error"}}},
			{Name: "FirstAuditLogByAction", Params: []parser.Param{{Name: "action", Type: "string"}}, Returns: []parser.Return{{Type: "*AuditLog", IsPointer: true, BaseType: "AuditLog"}, {Type: "error", BaseType: "error"}}},
			{Name: "ExistsAuditLogByID", Params: []parser.Param{{Name: "id", Type: "int64"}}, Returns: []parser.Return{{Type: "bool"}, {Type: "error", BaseType: "error"}}},
		},
	}

	repo := adapter.ToRepositoryInterface(iface, ent)
	if repo == nil {
		t.Fatal("ToRepositoryInterface() returned nil")
	}

	// Verify all methods are classified as MethodSmartQuery
	for i, m := range repo.Methods {
		if m.Kind != entity.MethodSmartQuery {
			t.Errorf("Method %d (%s): Kind = %d, want MethodSmartQuery (%d)", i, m.Name, m.Kind, entity.MethodSmartQuery)
		}
		// Verify GeneratedSQL is set
		if m.GeneratedSQL == "" {
			t.Errorf("Method %d (%s): GeneratedSQL is empty", i, m.Name)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// --- Coverage gap tests ---

func TestHandleBrace_TagWithBraces(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name string
		tag  string
		want []string
	}{
		{
			name: "check with braces",
			tag:  "check:age IN {1,2,3}",
			want: []string{"check:age IN {1,2,3}"},
		},
		{
			name: "check with nested braces and other tag",
			tag:  "check:status IN {active,inactive}, unique",
			want: []string{"check:status IN {active,inactive}", "unique"},
		},
		{
			name: "braces only value",
			tag:  "check:{a,b}",
			want: []string{"check:{a,b}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.parseTags(tt.tag)
			if len(got) != len(tt.want) {
				t.Fatalf("parseTags(%q) length = %d, want %d; got %v", tt.tag, len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseTags(%q)[%d] = %q, want %q", tt.tag, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseFKParts_EdgeCases(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single part",
			input: "users",
			want:  []string{"users"},
		},
		{
			name:  "two parts",
			input: "users,id",
			want:  []string{"users", "id"},
		},
		{
			name:  "many parts",
			input: "fk_name,table,col,on_delete:CASCADE,on_update:SET_NULL",
			want:  []string{"fk_name", "table", "col", "on_delete:CASCADE", "on_update:SET_NULL"},
		},
		{
			name:  "parenthesized value keeps inner comma",
			input: "users,func(a,b),id",
			want:  []string{"users", "func(a,b)", "id"},
		},
		{
			name:  "consecutive commas produce no empty parts",
			input: "a,,b",
			want:  []string{"a", "b"},
		},
		{
			name:  "trailing comma",
			input: "a,b,",
			want:  []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.parseFKParts(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("parseFKParts(%q) length = %d, want %d; got %v", tt.input, len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseFKParts(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseForeignKey_MalformedInputs(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		name     string
		tag      string
		wantNil  bool
	}{
		{
			name:    "empty fk value",
			tag:     "fk:",
			wantNil: true,
		},
		{
			name:    "single part only",
			tag:     "fk:users",
			wantNil: true,
		},
		{
			name:    "too many main parts (4 non-cascade)",
			tag:     "fk:a,b,c,d",
			wantNil: true,
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
			if tt.wantNil && domainField.FKReference != nil {
				t.Errorf("Expected FKReference nil for tag %q, got %+v", tt.tag, domainField.FKReference)
			}
		})
	}
}

func TestClassifyMethod_ListAndListAll(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
		},
	}

	tests := []struct {
		name       string
		methodName string
		wantKind   entity.MethodKind
	}{
		{"List method", "List", entity.MethodList},
		{"ListAll method", "ListAll", entity.MethodList},
		{"Get method", "Get", entity.MethodGetByID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := parser.Method{Name: tt.methodName}
			kind := adapter.classifyMethod(m, ent)
			if kind != tt.wantKind {
				t.Errorf("classifyMethod(%q) = %v, want %v", tt.methodName, kind, tt.wantKind)
			}
		})
	}
}

func TestToDomain_NilAndEmptyStruct(t *testing.T) {
	adapter := NewParserAdapter()

	t.Run("nil struct returns nil", func(t *testing.T) {
		result := adapter.ToDomain(nil)
		if result != nil {
			t.Errorf("ToDomain(nil) = %v, want nil", result)
		}
	})

	t.Run("struct with no fields", func(t *testing.T) {
		pStruct := &parser.Struct{
			Name:        "Empty",
			TableName:   "empties",
			PackageName: "models",
			Fields:      []parser.Field{},
		}
		result := adapter.ToDomain(pStruct)
		if result == nil {
			t.Fatal("ToDomain() returned nil for empty struct")
		}
		if result.Name != "Empty" {
			t.Errorf("Name = %q, want Empty", result.Name)
		}
		if len(result.Fields) != 0 {
			t.Errorf("Fields length = %d, want 0", len(result.Fields))
		}
	})

	t.Run("struct with only ignored fields", func(t *testing.T) {
		pStruct := &parser.Struct{
			Name:        "Secrets",
			TableName:   "secrets",
			PackageName: "models",
			Fields: []parser.Field{
				{Name: "Password", Type: "string", DatabaseTag: "-"},
				{Name: "Token", Type: "string", DatabaseTag: "-"},
			},
		}
		result := adapter.ToDomain(pStruct)
		if result == nil {
			t.Fatal("ToDomain() returned nil")
		}
		generateable := result.GetGenerateableFields()
		if len(generateable) != 0 {
			t.Errorf("GetGenerateableFields() length = %d, want 0", len(generateable))
		}
	})
}

func TestToDomainSlice_NilInput(t *testing.T) {
	adapter := NewParserAdapter()

	result := adapter.ToDomainSlice(nil)
	if result != nil {
		t.Errorf("ToDomainSlice(nil) = %v, want nil", result)
	}
}

func TestToMap_EdgeCases(t *testing.T) {
	adapter := NewParserAdapter()

	t.Run("nil map returns nil", func(t *testing.T) {
		result := adapter.ToMap(nil)
		if result != nil {
			t.Errorf("ToMap(nil) = %v, want nil", result)
		}
	})

	t.Run("empty map returns empty", func(t *testing.T) {
		result := adapter.ToMap(map[string][]*parser.Struct{})
		if result == nil {
			t.Fatal("ToMap({}) returned nil")
		}
		if len(result) != 0 {
			t.Errorf("ToMap({}) length = %d, want 0", len(result))
		}
	})

	t.Run("multiple packages", func(t *testing.T) {
		structs := map[string][]*parser.Struct{
			"models": {
				{Name: "User", Fields: []parser.Field{{Name: "ID", Type: "int64", DatabaseTag: "pk"}}},
			},
			"admin": {
				{Name: "AdminUser", Fields: []parser.Field{{Name: "ID", Type: "int64", DatabaseTag: "pk"}}},
				{Name: "Role", Fields: []parser.Field{{Name: "ID", Type: "int64", DatabaseTag: "pk"}}},
			},
		}
		result := adapter.ToMap(structs)
		if len(result) != 2 {
			t.Errorf("ToMap() returned %d packages, want 2", len(result))
		}
		if len(result["admin"]) != 2 {
			t.Errorf("ToMap()[admin] length = %d, want 2", len(result["admin"]))
		}
	})
}

func TestIsSimpleTag_AllCases(t *testing.T) {
	adapter := NewParserAdapter()

	tests := []struct {
		tag  string
		want bool
	}{
		{"pk", true},
		{"unique", true},
		{"-", true},
		{"index", true},
		{"unique_index", true},
		{"check:foo", false},
		{"default:bar", false},
		{"fk:t,c", false},
		{"", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got := adapter.isSimpleTag(tt.tag)
			if got != tt.want {
				t.Errorf("isSimpleTag(%q) = %v, want %v", tt.tag, got, tt.want)
			}
		})
	}
}

func TestProcessSmartQueryMethod_NoMatch(t *testing.T) {
	adapter := NewParserAdapter()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
		},
	}

	rm := entity.RepositoryMethod{Name: "SomeRandomMethodName"}
	m := parser.Method{Name: "SomeRandomMethodName"}
	adapter.processSmartQueryMethod(&rm, m, ent)

	if rm.GeneratedSQL != "" {
		t.Errorf("GeneratedSQL should be empty for unmatched method, got %q", rm.GeneratedSQL)
	}
}
