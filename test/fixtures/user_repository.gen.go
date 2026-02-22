package fixtures

import (
	"context"
	"database/sql"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, item *User) (*User, error) {
	query := `INSERT INTO user (username, email, active, created) VALUES ($1, $2, $3, $4) RETURNING i_d, username, email, active, created`
	var result User
	err := r.db.QueryRowContext(ctx, query, item.Username, item.Email, item.Active, item.Created).Scan(&result.ID, &result.Username, &result.Email, &result.Active, &result.Created)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT i_d, username, email, active, created FROM user WHERE i_d = $1`
	var item User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Username, &item.Email, &item.Active, &item.Created)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, item *User) error {
	query := `UPDATE user SET username = $1, email = $2, active = $3, created = $4 WHERE i_d = $5`
	_, err := r.db.ExecContext(ctx, query, item.Username, item.Email, item.Active, item.Created, item.ID)
	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM user WHERE i_d = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT i_d, username, email, active, created FROM user WHERE email = $1`
	var item User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&item.ID, &item.Username, &item.Email, &item.Active, &item.Created)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *UserRepositoryImpl) FindByActive(ctx context.Context, active bool) ([]*User, error) {
	query := `SELECT i_d, username, email, active, created FROM user WHERE active = $1`
	rows, err := r.db.QueryContext(ctx, query, active)
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

func (r *UserRepositoryImpl) FindRecentUsers(ctx context.Context, since int64) ([]*User, error) {
	query := `SELECT * FROM users WHERE created > $1 ORDER BY username`
	rows, err := r.db.QueryContext(ctx, query, since)
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

