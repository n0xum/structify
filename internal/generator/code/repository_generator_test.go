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

	if !strings.Contains(result, "func GetOrderItem(ctx context.Context, db *sql.DB, order_id int64, item_id int64") {
		t.Logf("Generated output:\n%s", result)
		t.Error("Generate() result missing GetOrderItem with composite PK params")
	}
	if !strings.Contains(result, "func DeleteOrderItem(ctx context.Context, db *sql.DB, order_id int64, item_id int64") {
		t.Error("Generate() result missing DeleteOrderItem with composite PK params")
	}
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

	if !strings.Contains(result, "func GetOrderWithUser") {
		t.Error("Generate() result missing GetOrderWithUser JOIN method")
	}
	if !strings.Contains(result, "type OrderWithUser struct") {
		t.Error("Generate() result missing OrderWithUser struct")
	}
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

	userEnt := gen.findEntity("user", entities)
	if userEnt == nil {
		t.Fatal("findEntity() returned nil for user")
	}
	if userEnt.Name != "User" {
		t.Errorf("findEntity() returned %s, want User", userEnt.Name)
	}

	orderEnt := gen.findEntity("order", entities)
	if orderEnt == nil {
		t.Fatal("findEntity() returned nil for order")
	}
	if orderEnt.Name != "Order" {
		t.Errorf("findEntity() returned entity with Name %s, want Order", orderEnt.Name)
	}

	productEnt := gen.findEntity("products", entities)
	if productEnt == nil {
		t.Fatal("findEntity() returned nil for products")
	}
	if productEnt.TableName != "products" {
		t.Errorf("findEntity() returned entity with TableName %s, want products", productEnt.TableName)
	}

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

// --- Multi-JOIN generator tests ---

func TestRepositoryGeneratorMultiJoinMethod(t *testing.T) {
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
			Name: "Category",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Title", Type: "string"},
			},
		},
		{
			Name: "Product",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{
					Name:        "UserID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "user", Column: "id"},
				},
				{
					Name:        "CategoryID",
					Type:        "int64",
					FKReference: &entity.FKReference{Table: "category", Column: "id"},
				},
				{Name: "Price", Type: "float64"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should generate the multi-join method
	if !strings.Contains(result, "func GetProductWithRelations") {
		t.Error("Generate() result missing GetProductWithRelations multi-join method")
	}

	// Should generate the multi-join result struct
	if !strings.Contains(result, "type ProductWithRelations struct") {
		t.Error("Generate() result missing ProductWithRelations struct")
	}

	// Should have multiple JOIN clauses
	joinCount := strings.Count(result, "JOIN")
	// At least 2 individual joins + the multi-join method's joins
	if joinCount < 3 {
		t.Errorf("Generate() result has %d JOIN clauses, want at least 3", joinCount)
	}

	// The multi-join struct should embed both related entities
	if !strings.Contains(result, "User\n") || !strings.Contains(result, "Category\n") {
		t.Logf("Generated output:\n%s", result)
		t.Error("ProductWithRelations struct should embed both User and Category")
	}
}

// --- Interface-driven generator tests ---

func TestGenerateFromInterfaceImplStruct(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{Name: "Create", Kind: entity.MethodCreate, EntityName: "User", Params: []entity.MethodParam{{Name: "item", Type: "*User"}}, ReturnsSingle: true, ReturnsError: true},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "type UserRepositoryImpl struct") {
		t.Error("missing impl struct")
	}
	if !strings.Contains(result, "db *sql.DB") {
		t.Error("impl struct missing db field")
	}
}

func TestGenerateFromInterfaceConstructor(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods:    []entity.RepositoryMethod{},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "func NewUserRepository(db *sql.DB) UserRepository") {
		t.Error("missing constructor")
	}
	if !strings.Contains(result, "return &UserRepositoryImpl{db: db}") {
		t.Error("constructor missing return")
	}
}

func TestGenerateFromInterfaceFindBy(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Email", Type: "string"},
			{Name: "Active", Type: "bool"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:          "FindByEmail",
				Kind:          entity.MethodFindBy,
				EntityName:    "User",
				Params:        []entity.MethodParam{{Name: "email", Type: "string"}},
				FindByFields:  []string{"Email"},
				ReturnsSingle: true,
				ReturnsError:  true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*User, error)") {
		t.Errorf("missing FindByEmail method signature, got:\n%s", result)
	}
	if !strings.Contains(result, "WHERE email = $1") {
		t.Error("FindByEmail missing WHERE clause")
	}
	if !strings.Contains(result, "QueryRowContext") {
		t.Error("FindByEmail should use QueryRowContext for single result")
	}
}

func TestGenerateFromInterfaceFindBySlice(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Active", Type: "bool"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:          "FindByActive",
				Kind:          entity.MethodFindBy,
				EntityName:    "User",
				Params:        []entity.MethodParam{{Name: "active", Type: "bool"}},
				FindByFields:  []string{"Active"},
				ReturnsSingle: false,
				ReturnsError:  true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "[]*User") {
		t.Error("FindByActive should return []*User")
	}
	if !strings.Contains(result, "QueryContext") {
		t.Error("FindByActive should use QueryContext for slice result")
	}
}

func TestGenerateFromInterfaceCustomSQL(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Active", Type: "bool"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:            "FindActiveUsers",
				Kind:            entity.MethodCustomSQL,
				EntityName:      "User",
				Params:          []entity.MethodParam{{Name: "active", Type: "bool"}},
				CustomSQL:       "SELECT * FROM users WHERE active = $1 ORDER BY username",
				ReturnsSingle:   false,
				ReturnsError:    true,
				HasEntityReturn: true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "SELECT * FROM users WHERE active = $1 ORDER BY username") {
		t.Error("CustomSQL should contain verbatim SQL")
	}
}

func TestGenerateFromInterfaceAllMethods(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Email", Type: "string"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{Name: "Create", Kind: entity.MethodCreate, EntityName: "User", Params: []entity.MethodParam{{Name: "item", Type: "*User"}}, ReturnsSingle: true, ReturnsError: true},
			{Name: "GetByID", Kind: entity.MethodGetByID, EntityName: "User", Params: []entity.MethodParam{{Name: "id", Type: "int64"}}, ReturnsSingle: true, ReturnsError: true},
			{Name: "Update", Kind: entity.MethodUpdate, EntityName: "User", Params: []entity.MethodParam{{Name: "item", Type: "*User"}}, ReturnsError: true},
			{Name: "Delete", Kind: entity.MethodDelete, EntityName: "User", Params: []entity.MethodParam{{Name: "id", Type: "int64"}}, ReturnsError: true},
			{Name: "List", Kind: entity.MethodList, EntityName: "User", ReturnsError: true},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	// All methods should use receiver
	for _, methodName := range []string{"Create", "GetByID", "Update", "Delete", "List"} {
		if !strings.Contains(result, "(r *UserRepositoryImpl) "+methodName) {
			t.Errorf("missing receiver-based method: %s", methodName)
		}
	}

	// Create should use INSERT
	if !strings.Contains(result, "INSERT INTO") {
		t.Error("Create missing INSERT INTO")
	}
	// GetByID should use SELECT...WHERE
	if !strings.Contains(result, "SELECT") && !strings.Contains(result, "WHERE") {
		t.Error("GetByID missing SELECT...WHERE")
	}
	// Update should use UPDATE...SET
	if !strings.Contains(result, "UPDATE") || !strings.Contains(result, "SET") {
		t.Error("Update missing UPDATE...SET")
	}
	// Delete should use DELETE FROM
	if !strings.Contains(result, "DELETE FROM") {
		t.Error("Delete missing DELETE FROM")
	}
	// All db operations should use r.db
	if !strings.Contains(result, "r.db.") {
		t.Error("methods should use r.db receiver")
	}
}

func TestGenerateFromInterfaceCustomSQLSingleReturn(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
			{Name: "Email", Type: "string"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:            "FindLatestUser",
				Kind:            entity.MethodCustomSQL,
				EntityName:      "User",
				Params:          []entity.MethodParam{},
				CustomSQL:       "SELECT * FROM users ORDER BY id DESC LIMIT 1",
				ReturnsSingle:   true,
				ReturnsError:    true,
				HasEntityReturn: true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "(*User, error)") {
		t.Error("CustomSQL single return should have (*User, error) return type")
	}
	if !strings.Contains(result, "QueryRowContext") {
		t.Error("CustomSQL single return should use QueryRowContext")
	}
	if !strings.Contains(result, "SELECT * FROM users ORDER BY id DESC LIMIT 1") {
		t.Error("CustomSQL should contain verbatim SQL")
	}
	if strings.Contains(result, "QueryContext") && !strings.Contains(result, "QueryRowContext") {
		t.Error("CustomSQL single return should NOT use QueryContext (multi-row)")
	}
}

func TestGenerateFromInterfaceSmartQueryCount(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Active", Type: "bool"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:         "CountActiveUsers",
				Kind:         entity.MethodSmartQuery,
				EntityName:   "User",
				Params:       []entity.MethodParam{{Name: "active", Type: "bool"}},
				GeneratedSQL: "SELECT COUNT(*) FROM user WHERE active = $1",
				ReturnsError: true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "(int64, error)") {
		t.Errorf("SmartQuery COUNT should return (int64, error), got:\n%s", result)
	}
	if !strings.Contains(result, "var count int64") {
		t.Error("SmartQuery COUNT should declare count variable")
	}
	if !strings.Contains(result, "Scan(&count)") {
		t.Error("SmartQuery COUNT should scan into count")
	}
	if !strings.Contains(result, "return count, err") {
		t.Error("SmartQuery COUNT should return count, err")
	}
}

func TestGenerateFromInterfaceSmartQueryExists(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Email", Type: "string"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:         "ExistsByEmail",
				Kind:         entity.MethodSmartQuery,
				EntityName:   "User",
				Params:       []entity.MethodParam{{Name: "email", Type: "string"}},
				GeneratedSQL: "SELECT EXISTS(SELECT 1 FROM user WHERE email = $1)",
				ReturnsError: true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "(bool, error)") {
		t.Errorf("SmartQuery EXISTS should return (bool, error), got:\n%s", result)
	}
	if !strings.Contains(result, "var exists bool") {
		t.Error("SmartQuery EXISTS should declare exists variable")
	}
	if !strings.Contains(result, "Scan(&exists)") {
		t.Error("SmartQuery EXISTS should scan into exists")
	}
	if !strings.Contains(result, "return exists, err") {
		t.Error("SmartQuery EXISTS should return exists, err")
	}
}

func TestGenerateFromInterfaceSmartQuerySliceReturn(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Active", Type: "bool"},
			{Name: "Username", Type: "string"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:          "FindRecentActive",
				Kind:          entity.MethodSmartQuery,
				EntityName:    "User",
				Params:        []entity.MethodParam{{Name: "active", Type: "bool"}},
				GeneratedSQL:  "SELECT i_d, active, username FROM user WHERE active = $1 ORDER BY i_d DESC",
				ReturnsSingle: false,
				ReturnsError:  true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "[]*User, error)") {
		t.Errorf("SmartQuery slice return should have []*User return type, got:\n%s", result)
	}
	if !strings.Contains(result, "QueryContext") {
		t.Error("SmartQuery slice return should use QueryContext")
	}
	if !strings.Contains(result, "defer rows.Close()") {
		t.Error("SmartQuery slice return should close rows")
	}
}

func TestGenerateFromInterfaceSmartQuerySingleReturn(t *testing.T) {
	gen := NewRepositoryGenerator()

	ent := &entity.Entity{
		Name: "User",
		Fields: []entity.Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Username", Type: "string"},
		},
	}

	repo := &entity.RepositoryInterface{
		Name:       "UserRepository",
		EntityName: "User",
		Methods: []entity.RepositoryMethod{
			{
				Name:          "FindOldestUser",
				Kind:          entity.MethodSmartQuery,
				EntityName:    "User",
				Params:        []entity.MethodParam{},
				GeneratedSQL:  "SELECT i_d, username FROM user ORDER BY i_d ASC LIMIT 1",
				ReturnsSingle: true,
				ReturnsError:  true,
			},
		},
	}

	result, err := gen.GenerateFromInterface(context.Background(), "repository", ent, repo)
	if err != nil {
		t.Fatalf("GenerateFromInterface() error = %v", err)
	}

	if !strings.Contains(result, "(*User, error)") {
		t.Errorf("SmartQuery single return should have (*User, error), got:\n%s", result)
	}
	if !strings.Contains(result, "QueryRowContext") {
		t.Error("SmartQuery single return should use QueryRowContext")
	}
	if !strings.Contains(result, "return &item, nil") {
		t.Error("SmartQuery single return should return &item")
	}
}
