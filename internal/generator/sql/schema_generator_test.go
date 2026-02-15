package sql

import (
	"context"
	"strings"
	"testing"

	"github.com/n0xum/structify/internal/domain/entity"
)

func TestSchemaGeneratorGenerate(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string", IsUnique: true},
				{Name: "Email", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	t.Logf("Generated SQL:\n%s", result)

	if !strings.Contains(result, "CREATE TABLE") {
		t.Error("Generate() result missing CREATE TABLE")
	}
	if !strings.Contains(result, "user") {
		t.Error("Generate() result missing table name")
	}
	if !strings.Contains(result, "id") {
		t.Error("Generate() result missing id column")
	}
	if !strings.Contains(result, "PRIMARY KEY") {
		t.Error("Generate() result missing PRIMARY KEY")
	}
}

func TestSchemaGeneratorGenerateWithIgnoredField(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
				{Name: "Password", Type: "string", IsIgnored: true},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if strings.Contains(result, "password") {
		t.Error("Generate() result contains ignored field password")
	}
}

func TestSchemaGeneratorGenerateWithCheckConstraint(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "Person",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Age", Type: "int", CheckExpr: "age >= 18"},
				{Name: "Email", Type: "string", CheckExpr: "email ~* '^[a-z]+'"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "CHECK (age >= 18)") {
		t.Error("Generate() result missing CHECK constraint for age")
	}
	if !strings.Contains(result, "CHECK (email ~* '^[a-z]+')") {
		t.Error("Generate() result missing CHECK constraint for email")
	}
}

func TestSchemaGeneratorGenerateWithDefaultConstraint(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Active", Type: "bool", DefaultVal: "true"},
				{Name: "Role", Type: "string", DefaultVal: "'user'"},
				{Name: "CreatedAt", Type: "time.Time", DefaultVal: "now()"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "DEFAULT true") {
		t.Error("Generate() result missing DEFAULT for Active")
	}
	if !strings.Contains(result, "DEFAULT 'user'") {
		t.Error("Generate() result missing DEFAULT for Role")
	}
	if !strings.Contains(result, "DEFAULT now()") {
		t.Error("Generate() result missing DEFAULT for CreatedAt")
	}
}

func TestSchemaGeneratorGenerateWithCombinedConstraints(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "Person",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Age", Type: "int", CheckExpr: "age >= 18", DefaultVal: "18"},
				{Name: "Active", Type: "bool", DefaultVal: "true"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "CHECK (age >= 18)") {
		t.Error("Generate() result missing CHECK constraint")
	}
	if !strings.Contains(result, "DEFAULT 18") {
		t.Error("Generate() result missing DEFAULT value")
	}
}

func TestSchemaGeneratorGenerateWithIndex(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name:      "User",
			TableName: "users",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Email", Type: "string", IndexName: "idx_email", IsIndexUnique: true},
				{Name: "Username", Type: "string", IndexName: "username_idx"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "CREATE UNIQUE INDEX") {
		t.Error("Generate() result missing CREATE UNIQUE INDEX")
	}
	if !strings.Contains(result, "CREATE INDEX") {
		t.Error("Generate() result missing CREATE INDEX")
	}
	if !strings.Contains(result, "idx_email") {
		t.Error("Generate() result missing idx_email index name")
	}
	if !strings.Contains(result, "username_idx") {
		t.Error("Generate() result missing username_idx index name")
	}
}

func TestSchemaGeneratorGenerateWithEnum(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "Order",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Status", Type: "string", EnumValues: []string{"pending", "processing", "shipped"}},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "CHECK (\"status\" IN ('pending', 'processing', 'shipped'))") {
		t.Error("Generate() result missing CHECK constraint for enum")
	}
}

func TestSchemaGeneratorGenerateWithCompositeIndex(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name:      "Product",
			TableName: "products",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Name", Type: "string", IndexName: "idx_name_category"},
				{Name: "Category", Type: "string", IndexName: "idx_name_category"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "CREATE INDEX") {
		t.Error("Generate() result missing CREATE INDEX")
	}
	if !strings.Contains(result, "idx_name_category") {
		t.Error("Generate() result missing idx_name_category index name")
	}
	if !strings.Contains(result, "(\"name\", \"category\")") {
		t.Error("Generate() result missing composite columns")
	}
}

func TestSchemaGeneratorGenerateWithForeignKey(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name:      "Order",
			TableName: "orders",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "UserID", Type: "int64", FKReference: &entity.FKReference{Table: "users", Column: "id"}},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "REFERENCES") {
		t.Error("Generate() result missing REFERENCES clause")
	}
	if !strings.Contains(result, "users") {
		t.Error("Generate() result missing referenced table name")
	}
}

func TestSchemaGeneratorGenerateWithForeignKeyCascade(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name:      "Order",
			TableName: "orders",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "UserID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "users", Column: "id"},
					FKOnDelete:  "CASCADE",
					FKOnUpdate:  "CASCADE",
				},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "ON DELETE CASCADE") {
		t.Error("Generate() result missing ON DELETE CASCADE")
	}
	if !strings.Contains(result, "ON UPDATE CASCADE") {
		t.Error("Generate() result missing ON UPDATE CASCADE")
	}
}

func TestSchemaGeneratorGenerateWithForeignKeySetNull(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name:      "Order",
			TableName: "orders",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "UserID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "users", Column: "id"},
					FKOnDelete:  "SET_NULL",
				},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "ON DELETE SET NULL") {
		t.Error("Generate() result missing ON DELETE SET NULL")
	}
}

func TestSchemaGeneratorGenerateWithCompositePrimaryKey(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "OrderItem",
			Fields: []entity.Field{
				{Name: "OrderID", Type: "int64", IsPrimary: true},
				{Name: "ItemID", Type: "int64", IsPrimary: true},
				{Name: "Quantity", Type: "int"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, `PRIMARY KEY ("order_id", "item_id")`) {
		t.Errorf("Generate() result missing composite PRIMARY KEY. Got:\n%s", result)
	}
}

func TestSchemaGeneratorGenerateWithCompositeUniqueConstraint(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "TenantID", Type: "int64", IsUnique: true, IndexGroup: "uq_tenant_role"},
				{Name: "UserID", Type: "int64", IsUnique: true, IndexGroup: "uq_tenant_role"},
				{Name: "Role", Type: "string", IsUnique: true, IndexGroup: "uq_tenant_role"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, `UNIQUE "uq_tenant_role"`) {
		t.Errorf("Generate() result missing composite UNIQUE constraint. Got:\n%s", result)
	}
	if !strings.Contains(result, `("tenant_id", "user_id", "role")`) {
		t.Errorf("Generate() result missing composite UNIQUE columns. Got:\n%s", result)
	}
}

func TestSchemaGeneratorGenerateWithCompositeForeignKey(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "OrderItem",
			Fields: []entity.Field{
				{Name: "OrderID", Type: "int64", IsPrimary: true},
				{Name: "ItemID", Type: "int64", IsPrimary: true},
				{Name: "Quantity", Type: "int"},
			},
		},
		{
			Name: "OrderItemNote",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "OrderID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "order_item", Column: "order_id"},
					FKGroup:     "fk_order_item",
				},
				{
					Name:        "ItemID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "order_item", Column: "item_id"},
					FKGroup:     "fk_order_item",
				},
				{Name: "Note", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, `FOREIGN KEY "fk_order_item" ("order_id", "item_id")`) {
		t.Errorf("Generate() result missing composite FOREIGN KEY. Got:\n%s", result)
	}
	if !strings.Contains(result, `REFERENCES "order_item" ("order_id", "item_id")`) {
		t.Errorf("Generate() result missing FK REFERENCES. Got:\n%s", result)
	}
}

func TestSchemaGeneratorGenerateWithCompositeForeignKeyCascade(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "OrderItem",
			Fields: []entity.Field{
				{Name: "OrderID", Type: "int64", IsPrimary: true},
				{Name: "ItemID", Type: "int64", IsPrimary: true},
				{Name: "Quantity", Type: "int"},
			},
		},
		{
			Name: "Shipment",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "OrderID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "order_item", Column: "order_id"},
					FKGroup:     "fk_shipment",
					FKOnDelete:  "CASCADE",
					FKOnUpdate:  "CASCADE",
				},
				{
					Name:        "ItemID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "order_item", Column: "item_id"},
					FKGroup:     "fk_shipment",
				},
				{Name: "Tracking", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, `ON DELETE CASCADE`) {
		t.Errorf("Generate() result missing ON DELETE CASCADE. Got:\n%s", result)
	}
	if !strings.Contains(result, `ON UPDATE CASCADE`) {
		t.Errorf("Generate() result missing ON UPDATE CASCADE. Got:\n%s", result)
	}
}

func TestSchemaGeneratorGroupFieldsByFK(t *testing.T) {
	gen := NewSchemaGenerator()

	fields := []entity.Field{
		{Name: "ID", Type: "int64", IsPrimary: true},
		{
			Name:        "OrderID",
			Type:        "int64",
			FKReference: &entity.FKReference{Table: "order_item", Column: "order_id"},
			FKGroup:     "fk_order",
		},
		{
			Name:        "ItemID",
			Type:        "int64",
			FKReference: &entity.FKReference{Table: "order_item", Column: "item_id"},
			FKGroup:     "fk_order",
		},
		{
			Name:        "UserID",
			Type:        "int64",
			FKReference: &entity.FKReference{Table: "users", Column: "id"},
		},
	}

	groups := gen.groupFieldsByFK(fields)

	if len(groups) != 1 {
		t.Errorf("groupFieldsByFK() returned %d groups, want 1", len(groups))
	}

	fkFields, ok := groups["fk_order"]
	if !ok {
		t.Fatal("groupFieldsByFK() missing fk_order group")
	}

	if len(fkFields) != 2 {
		t.Errorf("fk_order group has %d fields, want 2", len(fkFields))
	}
}

func TestSchemaGeneratorHasCompositeForeignKey(t *testing.T) {
	gen := NewSchemaGenerator()

	tests := []struct {
		name  string
		fields []entity.Field
		want  bool
	}{
		{
			name:  "no FKs",
			fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Name", Type: "string"},
			},
			want: false,
		},
		{
			name: "single FK",
			fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "UserID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "users", Column: "id"},
				},
			},
			want: false,
		},
		{
			name: "composite FK",
			fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "OrderID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "order_item", Column: "order_id"},
					FKGroup:     "fk_order",
				},
				{
					Name:        "ItemID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "order_item", Column: "item_id"},
					FKGroup:     "fk_order",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gen.hasCompositeForeignKey(tt.fields); got != tt.want {
				t.Errorf("hasCompositeForeignKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
