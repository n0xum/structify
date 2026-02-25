package ecommerce

import (
	"context"
	"database/sql"
)

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, item *Product) (*Product, error) {
	query := `INSERT INTO "product" (category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at`
	var result Product
	err := r.db.QueryRowContext(ctx, query, item.CategoryID, item.SKU, item.Name, item.Description, item.Price, item.StockQty, item.Status, item.Attributes, item.FeaturedAt, item.CreatedAt, item.UpdatedAt).Scan(&result.ID, &result.CategoryID, &result.SKU, &result.Name, &result.Description, &result.Price, &result.StockQty, &result.Status, &result.Attributes, &result.FeaturedAt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id int64) (*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE id = $1`
	var item Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, item *Product) error {
	query := `UPDATE "product" SET category_id = $1, sku = $2, name = $3, description = $4, price = $5, stock_qty = $6, status = $7, attributes = $8, featured_at = $9, created_at = $10, updated_at = $11 WHERE id = $12`
	_, err := r.db.ExecContext(ctx, query, item.CategoryID, item.SKU, item.Name, item.Description, item.Price, item.StockQty, item.Status, item.Attributes, item.FeaturedAt, item.CreatedAt, item.UpdatedAt, item.ID)
	return err
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM "product" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ProductRepositoryImpl) List(ctx context.Context) ([]*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) FindBySKU(ctx context.Context, sku string) (*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE sku = $1`
	var item Product
	err := r.db.QueryRowContext(ctx, query, sku).Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ProductRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID int64) ([]*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) FindByStatus(ctx context.Context, status string) ([]*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) FindByCategoryIDAndStatus(ctx context.Context, categoryID int64, status string) ([]*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE category_id = $1 AND status = $2`
	rows, err := r.db.QueryContext(ctx, query, categoryID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) ListProductsByStatus(ctx context.Context, status string) ([]*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) ListProductsByCategoryID(ctx context.Context, categoryID int64) ([]*Product, error) {
	query := `SELECT id, category_id, sku, name, description, price, stock_qty, status, attributes, featured_at, created_at, updated_at FROM "product" WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) CountProductsByCategoryID(ctx context.Context, categoryID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM "product" WHERE category_id = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	return count, err
}

func (r *ProductRepositoryImpl) CountProductsByStatus(ctx context.Context, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM "product" WHERE status = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}

func (r *ProductRepositoryImpl) ExistsProductBySKU(ctx context.Context, sku string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "product" WHERE sku = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, sku).Scan(&exists)
	return exists, err
}

func (r *ProductRepositoryImpl) FindInStock(ctx context.Context, limit int) ([]*Product, error) {
	query := `SELECT * FROM product WHERE stock_qty > 0 AND status = 'active' ORDER BY created_at DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) FindByPriceRange(ctx context.Context, minPrice float64, maxPrice float64) ([]*Product, error) {
	query := `SELECT * FROM product WHERE price BETWEEN $1 AND $2 AND status = 'active' ORDER BY price ASC`
	rows, err := r.db.QueryContext(ctx, query, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.SKU, &item.Name, &item.Description, &item.Price, &item.StockQty, &item.Status, &item.Attributes, &item.FeaturedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepositoryImpl) DecrementStock(ctx context.Context, id int64, qty int) error {
	query := `UPDATE product SET stock_qty = stock_qty - $2 WHERE id = $1 AND stock_qty >= $2`
	_, err := r.db.ExecContext(ctx, query, id, qty)
	return err
}

