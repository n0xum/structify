package ecommerce

import "time"

// Order is a purchase placed by a Customer.
//
// Demonstrates:
//   - fk with ON DELETE CASCADE (orders removed with the customer)
//   - enum for order lifecycle
//   - default for epoch-style timestamp
//   - nullable fields (ShippedAt, DeliveredAt have no default → NULL)
type Order struct {
	ID          int64   `db:"pk"`
	CustomerID  int64   `db:"fk:customer,id,on_delete:CASCADE"`
	Status      string  `db:"default:'pending',enum:pending,confirmed,shipped,delivered,cancelled,refunded"`
	TotalAmount float64 `db:"check:total_amount >= 0"`
	Currency    string  `db:"default:'USD',check:length(currency) = 3"`
	CreatedAt   int64   `db:"default:extract(epoch from now())"`
	ShippedAt   int64
	DeliveredAt int64
}

// Customer is a registered buyer.
//
// Demonstrates:
//   - unique with check combined in one db tag
//   - unique_index (index-backed uniqueness, vs constraint-backed)
//   - enum with default
//   - default on bool and timestamp
type Customer struct {
	ID        int64     `db:"pk"`
	Email     string    `db:"unique,check:email ~* '^[^@]+@[^@]+\\.[^@]+$'"`
	Username  string    `db:"unique_index:idx_customer_username"`
	FullName  string    `db:"check:length(full_name) > 0"`
	Phone     string
	Tier      string    `db:"default:'free',enum:free,silver,gold,platinum"`
	Active    bool      `db:"default:true"`
	CreatedAt time.Time `db:"default:now()"`
}

// OrderItem is a single line in an Order for a specific Product.
//
// Demonstrates:
//   - composite primary key (order_id + product_id)
//   - pk and fk combined on each PK column
//   - ON DELETE CASCADE from both parent tables
//   - ON UPDATE CASCADE on the product FK
//   - check on quantity and unit price
//   - default for discount
type OrderItem struct {
	OrderID   int64   `db:"pk,fk:order,id,on_delete:CASCADE"`
	ProductID int64   `db:"pk,fk:product,id,on_delete:RESTRICT,on_update:CASCADE"`
	Quantity  int     `db:"check:quantity > 0"`
	UnitPrice float64 `db:"check:unit_price >= 0"`
	Discount  float64 `db:"default:0,check:discount >= 0"`
}
