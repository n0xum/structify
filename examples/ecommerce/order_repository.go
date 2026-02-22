package ecommerce

import "context"

// OrderRepository defines all database operations for the Order entity.
//
// Demonstrates:
//   - FindBy across FK customer column
//   - SmartQuery Count / Exists patterns
//   - CustomSQL for revenue queries and bulk status updates
type OrderRepository interface {
	// --- CRUD ---

	Create(ctx context.Context, order *Order) (*Order, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Order, error)

	// --- FindBy ---

	// FindByCustomerID returns all orders placed by a customer.
	FindByCustomerID(ctx context.Context, customerID int64) ([]*Order, error)

	// FindByStatus returns all orders in a given status.
	FindByStatus(ctx context.Context, status string) ([]*Order, error)

	// FindByCustomerIDAndStatus filters by customer and status together.
	FindByCustomerIDAndStatus(ctx context.Context, customerID int64, status string) ([]*Order, error)

	// --- SmartQuery ---

	// CountOrdersByCustomerID emits SELECT COUNT(*) FROM "order" WHERE customer_id = $1.
	CountOrdersByCustomerID(ctx context.Context, customerID int64) (int64, error)

	// CountOrdersByStatus emits SELECT COUNT(*) FROM "order" WHERE status = $1.
	CountOrdersByStatus(ctx context.Context, status string) (int64, error)

	// ExistsOrderByCustomerID emits SELECT EXISTS(SELECT 1 FROM "order" WHERE customer_id = $1).
	ExistsOrderByCustomerID(ctx context.Context, customerID int64) (bool, error)

	// ListOrdersByStatusOrderByCreatedAtDesc lists orders newest-first.
	ListOrdersByStatusOrderByCreatedAtDesc(ctx context.Context, status string) ([]*Order, error)

	// --- CustomSQL ---

	//sql:"SELECT * FROM \"order\" WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2"
	FindRecentByCustomer(ctx context.Context, customerID int64, limit int) ([]*Order, error)

	// COALESCE handles the NULL that SUM returns when no rows match.
	//sql:"SELECT COALESCE(SUM(total_amount), 0.0) FROM \"order\" WHERE customer_id = $1 AND status = 'delivered'"
	GetTotalSpentByCustomer(ctx context.Context, customerID int64) (float64, error)

	//sql:"UPDATE \"order\" SET status = 'cancelled' WHERE status = 'pending' AND created_at < $1"
	CancelOldPendingOrders(ctx context.Context, olderThan int64) error
}
