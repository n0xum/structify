package fixtures

import (
	"database/sql"
	"context"
)

type User struct {
	ID int64
	Username string
	Email string
	Active bool
	Created int64
}

func CreateUser(ctx context.Context, db *sql.DB, item *User) (*User, error) {
	query := `INSERT INTO users (username, email, active, created) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := db.QueryRowContext(ctx, query, item.Username, item.Email, item.Active, item.Created).Scan(&id)
	if err != nil {
		return nil, err
	}
	return GetUserByID(ctx, db, id)
}

func GetUserByID(ctx context.Context, db *sql.DB, id int64) (*User, error) {
	query := `SELECT id, username, email, active, created FROM users WHERE id = $1`
	var item User
	err := db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Username, &item.Email, &item.Active, &item.Created)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, item *User) error {
	query := `UPDATE users SET username = $1, email = $2, active = $3, created = $4 WHERE id = $5`
	_, err := db.ExecContext(ctx, query, item.Username, item.Email, item.Active, item.Created, item.ID)
	return err
}

func DeleteUser(ctx context.Context, db *sql.DB, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func ListUser(ctx context.Context, db *sql.DB) ([]*User, error) {
	query := `SELECT id, username, email, active, created FROM users ORDER BY id`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*User
	for rows.Next() {
		var item User
		if err := rows.Scan(&item.ID, &item.Username, &item.Email, &item.Active, &item.Created); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

type Product struct {
	ID int64
	Name string
	Price float64
	Description string
}

func CreateProduct(ctx context.Context, db *sql.DB, item *Product) (*Product, error) {
	query := `INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := db.QueryRowContext(ctx, query, item.Name, item.Price, item.Description).Scan(&id)
	if err != nil {
		return nil, err
	}
	return GetProductByID(ctx, db, id)
}

func GetProductByID(ctx context.Context, db *sql.DB, id int64) (*Product, error) {
	query := `SELECT id, name, price, description FROM products WHERE id = $1`
	var item Product
	err := db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name, &item.Price, &item.Description)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func UpdateProduct(ctx context.Context, db *sql.DB, item *Product) error {
	query := `UPDATE products SET name = $1, price = $2, description = $3 WHERE id = $4`
	_, err := db.ExecContext(ctx, query, item.Name, item.Price, item.Description, item.ID)
	return err
}

func DeleteProduct(ctx context.Context, db *sql.DB, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func ListProduct(ctx context.Context, db *sql.DB) ([]*Product, error) {
	query := `SELECT id, name, price, description FROM products ORDER BY id`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
