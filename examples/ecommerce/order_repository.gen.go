package ecommerce

import (
	"context"
	"database/sql"
)

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) Create(ctx context.Context, item *Order) (*Order, error) {
	query := `INSERT INTO "order" (customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at`
	var result Order
	err := r.db.QueryRowContext(ctx, query, item.CustomerID, item.Status, item.TotalAmount, item.Currency, item.CreatedAt, item.ShippedAt, item.DeliveredAt).Scan(&result.ID, &result.CustomerID, &result.Status, &result.TotalAmount, &result.Currency, &result.CreatedAt, &result.ShippedAt, &result.DeliveredAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id int64) (*Order, error) {
	query := `SELECT id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at FROM "order" WHERE id = $1`
	var item Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrderRepositoryImpl) Update(ctx context.Context, item *Order) error {
	query := `UPDATE "order" SET customer_id = $1, status = $2, total_amount = $3, currency = $4, created_at = $5, shipped_at = $6, delivered_at = $7 WHERE id = $8`
	_, err := r.db.ExecContext(ctx, query, item.CustomerID, item.Status, item.TotalAmount, item.Currency, item.CreatedAt, item.ShippedAt, item.DeliveredAt, item.ID)
	return err
}

func (r *OrderRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM "order" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *OrderRepositoryImpl) List(ctx context.Context) ([]*Order, error) {
	query := `SELECT id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at FROM "order" ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepositoryImpl) FindByCustomerID(ctx context.Context, customerID int64) ([]*Order, error) {
	query := `SELECT id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at FROM "order" WHERE customer_id = $1`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepositoryImpl) FindByStatus(ctx context.Context, status string) ([]*Order, error) {
	query := `SELECT id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at FROM "order" WHERE status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepositoryImpl) FindByCustomerIDAndStatus(ctx context.Context, customerID int64, status string) ([]*Order, error) {
	query := `SELECT id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at FROM "order" WHERE customer_id = $1 AND status = $2`
	rows, err := r.db.QueryContext(ctx, query, customerID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepositoryImpl) CountOrdersByCustomerID(ctx context.Context, customerID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM "order" WHERE customer_id = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&count)
	return count, err
}

func (r *OrderRepositoryImpl) CountOrdersByStatus(ctx context.Context, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM "order" WHERE status = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}

func (r *OrderRepositoryImpl) ExistsOrderByCustomerID(ctx context.Context, customerID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "order" WHERE customer_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&exists)
	return exists, err
}

func (r *OrderRepositoryImpl) ListOrdersByStatusOrderByCreatedAtDesc(ctx context.Context, status string) ([]*Order, error) {
	query := `SELECT id, customer_id, status, total_amount, currency, created_at, shipped_at, delivered_at FROM "order" WHERE status = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepositoryImpl) FindRecentByCustomer(ctx context.Context, customerID int64, limit int) ([]*Order, error) {
	query := `SELECT * FROM "order" WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, customerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Order
	for rows.Next() {
		var item Order
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Status, &item.TotalAmount, &item.Currency, &item.CreatedAt, &item.ShippedAt, &item.DeliveredAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderRepositoryImpl) GetTotalSpentByCustomer(ctx context.Context, customerID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(total_amount), 0.0) FROM "order" WHERE customer_id = $1 AND status = 'delivered'`
	var result float64
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&result)
	return result, err
}

func (r *OrderRepositoryImpl) CancelOldPendingOrders(ctx context.Context, olderThan int64) error {
	query := `UPDATE "order" SET status = 'cancelled' WHERE status = 'pending' AND created_at < $1`
	_, err := r.db.ExecContext(ctx, query, olderThan)
	return err
}

