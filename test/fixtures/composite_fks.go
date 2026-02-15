package fixtures

// OrderItem with composite primary key (order_id, item_id)
type CompositeFKOrderItem struct {
	OrderID int64 `db:"pk"`
	ItemID  int64 `db:"pk"`
	Quantity int    `db:"check:quantity > 0"`
	Price    float64
}

// OrderItemNote references OrderItem with composite FK
// Multiple fields (order_id, item_id) together reference the composite PK of order_items
type CompositeFKOrderItemNote struct {
	NoteID  int64  `db:"pk"`
	OrderID int64  `db:"fk:fk_order_item,composite_f_k_order_item,order_id"`
	ItemID  int64  `db:"fk:fk_order_item,composite_f_k_order_item,item_id"`
	Note    string
	Created int64 `db:"default:0"`
}

// Shipment with composite FK referencing OrderItem with CASCADE
type CompositeFKShipment struct {
	ShipmentID int64 `db:"pk"`
	OrderID    int64 `db:"fk:fk_shipment_order,composite_f_k_order_item,order_id,on_delete:CASCADE,on_update:CASCADE"`
	ItemID     int64 `db:"fk:fk_shipment_order,composite_f_k_order_item,item_id"`
	Tracking   string
}

// ProductCategory with composite primary key
type CompositeFKProductCategory struct {
	CategoryID int64 `db:"pk"`
	RegionID   int64 `db:"pk"`
	Name       string
}

// ProductCategoryAssignment references ProductCategory with composite FK
type CompositeFKProductCategoryAssignment struct {
	AssignmentID int64 `db:"pk"`
	ProductID    int64 `db:"fk:products,id"`
	CategoryID   int64 `db:"fk:fk_product_category,composite_f_k_product_category,category_id"`
	RegionID     int64 `db:"fk:fk_product_category,composite_f_k_product_category,region_id"`
}
