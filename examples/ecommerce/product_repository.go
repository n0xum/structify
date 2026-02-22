package ecommerce

import "context"

// ProductRepository defines all database operations for the Product entity.
//
// Demonstrates:
//   - FindBy across a FK column (category)
//   - FindBy with AND of status + category
//   - SmartQuery: ListXByY, CountXByY, ExistsXByY
//   - CustomSQL for search, bulk update, and aggregation
type ProductRepository interface {
	// --- CRUD ---

	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Product, error)

	// --- FindBy ---

	// FindBySKU returns the product with the given stock-keeping unit.
	FindBySKU(ctx context.Context, sku string) (*Product, error)

	// FindByCategoryID returns all products in a category.
	FindByCategoryID(ctx context.Context, categoryID int64) ([]*Product, error)

	// FindByStatus returns all products in the given lifecycle status.
	FindByStatus(ctx context.Context, status string) ([]*Product, error)

	// FindByCategoryIDAndStatus filters by category and lifecycle status.
	FindByCategoryIDAndStatus(ctx context.Context, categoryID int64, status string) ([]*Product, error)

	// --- SmartQuery ---

	// ListProductsByStatus emits SELECT … FROM product WHERE status = $1.
	ListProductsByStatus(ctx context.Context, status string) ([]*Product, error)

	// ListProductsByCategoryID emits SELECT … FROM product WHERE category_id = $1.
	ListProductsByCategoryID(ctx context.Context, categoryID int64) ([]*Product, error)

	// CountProductsByCategoryID emits SELECT COUNT(*) FROM product WHERE category_id = $1.
	CountProductsByCategoryID(ctx context.Context, categoryID int64) (int64, error)

	// CountProductsByStatus emits SELECT COUNT(*) FROM product WHERE status = $1.
	CountProductsByStatus(ctx context.Context, status string) (int64, error)

	// ExistsProductBySKU emits SELECT EXISTS(SELECT 1 FROM product WHERE sku = $1).
	ExistsProductBySKU(ctx context.Context, sku string) (bool, error)

	// --- CustomSQL ---

	//sql:"SELECT * FROM product WHERE stock_qty > 0 AND status = 'active' ORDER BY created_at DESC LIMIT $1"
	FindInStock(ctx context.Context, limit int) ([]*Product, error)

	//sql:"SELECT * FROM product WHERE price BETWEEN $1 AND $2 AND status = 'active' ORDER BY price ASC"
	FindByPriceRange(ctx context.Context, minPrice float64, maxPrice float64) ([]*Product, error)

	//sql:"UPDATE product SET stock_qty = stock_qty - $2 WHERE id = $1 AND stock_qty >= $2"
	DecrementStock(ctx context.Context, id int64, qty int) error
}
