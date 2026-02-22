//go:build integration

package integration

// Test structs covering the full range of supported db tags and types.
// These serve as inputs for the integration tests below.

// BlogPost tests: pk, unique, ignored field, custom table name
type BlogPost struct {
	ID        int64  `db:"pk"`
	Slug      string `db:"unique"`
	Title     string
	Body      string
	Published bool
	Internal  string `db:"-"`
}

// OrderItem tests: composite fields, float types, pointer-like naming
type OrderItem struct {
	ID        int64 `db:"pk"`
	OrderID   int64
	ProductID int64
	Quantity  int32
	UnitPrice float64
	Discount  float32
}

// Account tests: bool NOT NULL, multiple unique fields
type Account struct {
	ID       int64  `db:"pk"`
	Email    string `db:"unique"`
	Username string `db:"unique"`
	Active   bool
	Balance  float64
}

// AuditLog tests: no primary key, all string fields
type AuditLog struct {
	Action    string
	TableName string
	RecordID  string
	ChangedBy string
}
