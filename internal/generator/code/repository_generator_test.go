package code

import (
	"context"
	"strings"
	"testing"

	"github.com/n0xum/structify/internal/domain/entity"
)

func TestRepositoryGeneratorGenerate(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name:    "User",
			Package: "models",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string", IsUnique: true},
				{Name: "Email", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "package models") {
		t.Error("Generate() result missing package declaration")
	}
	if !strings.Contains(result, "import") {
		t.Error("Generate() result missing import block")
	}
	if !strings.Contains(result, "database/sql") {
		t.Error("Generate() result missing database/sql import")
	}
	if !strings.Contains(result, "type User struct") {
		t.Error("Generate() result missing User struct")
	}
	if !strings.Contains(result, "func CreateUser") {
		t.Error("Generate() result missing CreateUser function")
	}
	if !strings.Contains(result, "func GetUserByID") {
		t.Error("Generate() result missing GetUserByID function")
	}
	if !strings.Contains(result, "func UpdateUser") {
		t.Error("Generate() result missing UpdateUser function")
	}
	if !strings.Contains(result, "func DeleteUser") {
		t.Error("Generate() result missing DeleteUser function")
	}
	if !strings.Contains(result, "func ListUser") {
		t.Error("Generate() result missing ListUser function")
	}
}

func TestRepositoryGeneratorGenerateWithIgnoredField(t *testing.T) {
	gen := NewRepositoryGenerator()

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

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if strings.Contains(result, "Password string") {
		t.Error("Generate() result contains ignored Password field in struct")
	}
}

func TestRepositoryGeneratorCreateMethodUsesCorrectAPI(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, ".Scan(&id)") {
		t.Error("Generate() Create method missing .Scan() call")
	}
	if !strings.Contains(result, "QueryRowContext(ctx, query,") {
		t.Error("Generate() Create method using wrong API")
	}
}

func TestRepositoryGeneratorGetByIDMethodUsesCorrectAPI(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, ".Scan(") {
		t.Error("Generate() GetByID method missing .Scan() call")
	}
	if !strings.Contains(result, "QueryRowContext(ctx, query, id)") {
		t.Error("Generate() GetByID method using wrong API")
	}
}

func TestRepositoryGeneratorListMethodUsesCorrectAPI(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "rows.Scan(") {
		t.Error("Generate() List method missing .Scan() call")
	}
	if !strings.Contains(result, "defer rows.Close()") {
		t.Error("Generate() List method missing defer rows.Close()")
	}
}

func TestRepositoryGeneratorGenerateWithCompositePrimaryKey(t *testing.T) {
	gen := NewRepositoryGenerator()

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

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should have GetOrderItem method instead of GetOrderItemByID
	if !strings.Contains(result, "func GetOrderItem(ctx context.Context, db *sql.DB, order_id int64, item_id int64") {
		t.Logf("Generated output:\n%s", result)
		t.Error("Generate() result missing GetOrderItem with composite PK params")
	}
	// Should have DeleteOrderItem with composite PK params
	if !strings.Contains(result, "func DeleteOrderItem(ctx context.Context, db *sql.DB, order_id int64, item_id int64") {
		t.Error("Generate() result missing DeleteOrderItem with composite PK params")
	}
	// Should not have GetOrderItemByID
	if strings.Contains(result, "func GetOrderItemByID") {
		t.Error("Generate() result should not have GetOrderItemByID for composite PK")
	}
}

func TestRepositoryGeneratorGenerateWithForeignKey(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Name", Type: "string"},
			},
		},
		{
			Name: "Order",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "UserID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "user", Column: "id"},
				},
				{Name: "Total", Type: "float64"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should generate JOIN method
	if !strings.Contains(result, "func GetOrderWithUser") {
		t.Error("Generate() result missing GetOrderWithUser JOIN method")
	}
	// Should generate result struct
	if !strings.Contains(result, "type OrderWithUser struct") {
		t.Error("Generate() result missing OrderWithUser struct")
	}
	// Should have JOIN clause in query
	if !strings.Contains(result, "JOIN") {
		t.Error("Generate() result missing JOIN clause")
	}
}

func TestRepositoryGeneratorGroupFKsByTable(t *testing.T) {
	gen := NewRepositoryGenerator()

	fields := []entity.Field{
		{Name: "ID", Type: "int64", IsPrimary: true},
		{
			Name:        "UserID",
			Type:        "int64",
			FKReference: &entity.FKReference{Table: "users", Column: "id"},
		},
		{
			Name:        "ProductID",
			Type:        "int64",
			FKReference: &entity.FKReference{Table: "products", Column: "id"},
		},
		{
			Name:        "CategoryID",
			Type:        "int64",
			FKReference: &entity.FKReference{Table: "products", Column: "category_id"},
		},
	}

	groups := gen.groupFKsByTable(fields)

	if len(groups) != 2 {
		t.Errorf("groupFKsByTable() returned %d groups, want 2", len(groups))
	}

	userFKs, ok := groups["users"]
	if !ok {
		t.Fatal("groupFKsByTable() missing users group")
	}
	if len(userFKs) != 1 {
		t.Errorf("users group has %d fields, want 1", len(userFKs))
	}

	productFKs, ok := groups["products"]
	if !ok {
		t.Fatal("groupFKsByTable() missing products group")
	}
	if len(productFKs) != 2 {
		t.Errorf("products group has %d fields, want 2", len(productFKs))
	}
}

func TestRepositoryGeneratorFindEntity(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
		{Name: "Order", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
		{
			Name:      "Product",
			TableName: "products",
			Fields:    []entity.Field{{Name: "ID", Type: "int64"}},
		},
	}

	// Find by table name (uses GetTableName which converts to snake_case)
	userEnt := gen.findEntity("user", entities)
	if userEnt == nil {
		t.Fatal("findEntity() returned nil for user")
	}
	if userEnt.Name != "User" {
		t.Errorf("findEntity() returned %s, want User", userEnt.Name)
	}

	// Find by table name (default snake_case from Order)
	orderEnt := gen.findEntity("order", entities)
	if orderEnt == nil {
		t.Fatal("findEntity() returned nil for order")
	}
	if orderEnt.Name != "Order" {
		t.Errorf("findEntity() returned entity with Name %s, want Order", orderEnt.Name)
	}

	// Find by explicit table name
	productEnt := gen.findEntity("products", entities)
	if productEnt == nil {
		t.Fatal("findEntity() returned nil for products")
	}
	if productEnt.TableName != "products" {
		t.Errorf("findEntity() returned entity with TableName %s, want products", productEnt.TableName)
	}

	// Not found
	notFound := gen.findEntity("category", entities)
	if notFound != nil {
		t.Error("findEntity() returned non-nil for non-existent category")
	}
}

func TestRepositoryGeneratorGetEntityColumns(t *testing.T) {
	gen := NewRepositoryGenerator()

	entity := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Secret", Type: "string", IsIgnored: true},
		},
	}

	columns := gen.getEntityColumns(entity)

	if len(columns) != 2 {
		t.Errorf("getEntityColumns() returned %d columns, want 2", len(columns))
	}

	expectedColumns := []string{"id", "username"}
	for i, col := range columns {
		if col != expectedColumns[i] {
			t.Errorf("getEntityColumns()[%d] = %s, want %s", i, col, expectedColumns[i])
		}
	}
}
