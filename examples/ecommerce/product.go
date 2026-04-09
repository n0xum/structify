package ecommerce

import (
	"encoding/json"
	"time"
)

// Product is an item available for purchase.
//
// Demonstrates:
//   - fk with ON DELETE RESTRICT (can't delete a category that has products)
//   - composite index on (category_id, status) for category listing pages
//   - check constraint on price and stock
//   - enum for lifecycle status
//   - json.RawMessage for flexible attribute storage (maps to JSONB)
//   - time.Time for timestamps
type Product struct {
	ID           int64           `db:"pk"`
	CategoryID   int64           `db:"fk:product_category,id,on_delete:RESTRICT"`
	SKU          string          `db:"unique,check:length(sku) > 0"`
	Name         string          `db:"index:idx_product_name,check:length(name) > 0"`
	Description  string
	Price        float64         `db:"check:price >= 0"`
	StockQty     int             `db:"default:0,check:stock_qty >= 0"`
	Status       string          `db:"default:'draft',enum:active,inactive,discontinued,draft"`
	Attributes   json.RawMessage `db:"default:'{}'"`
	FeaturedAt   time.Time
	CreatedAt    time.Time       `db:"default:now()"`
	UpdatedAt    time.Time
}

// ProductCategory groups products.
//
// Demonstrates:
//   - unique column constraint
//   - self-referential FK (parent category)
//   - default values
type ProductCategory struct {
	ID          int64  `db:"pk"`
	Name        string `db:"unique,check:length(name) > 0"`
	Slug        string `db:"unique_index:idx_category_slug"`
	Description string
	ParentID    int64  `db:"fk:product_category,id,on_delete:SET_NULL"`
	SortOrder   int    `db:"default:0"`
}
