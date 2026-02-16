package fixtures

// OrderItem with composite primary key (order_id, item_id)
type CompositePKOrderItem struct {
	OrderID int64 `db:"pk"`
	ItemID  int64 `db:"pk"`
	Quantity int    `db:"check:quantity > 0"`
	Price    float64
}

// UserRole with composite unique constraint (tenant_id, role)
type CompositeUniqueUser struct {
	ID       int64  `db:"pk"`
	TenantID int64  `db:"unique:uq_tenant_role"`
	Username string `db:"unique:uq_tenant_role"`
	Role     string `db:"unique:uq_tenant_role"`
	Email    string `db:"unique"`
}

// Address with composite unique constraint
type CompositeUniqueAddress struct {
	ID       int64  `db:"pk"`
	UserID   int64  `db:"unique:uq_user_address"`
	AddressType string `db:"unique:uq_user_address"`
	Address  string
	City     string
	Zip      string
}

// Permission with composite primary key
type CompositePKPermission struct {
	RoleID    int64 `db:"pk"`
	ResourceID int64 `db:"pk"`
	Permission string `db:"unique"`
	CanWrite   bool   `db:"default:false"`
}
