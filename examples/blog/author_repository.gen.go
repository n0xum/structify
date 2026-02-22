package blog

import (
	"context"
	"database/sql"
)

type AuthorRepositoryImpl struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) AuthorRepository {
	return &AuthorRepositoryImpl{db: db}
}

func (r *AuthorRepositoryImpl) Create(ctx context.Context, item *Author) (*Author, error) {
	query := `INSERT INTO "author" (email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at`
	var result Author
	err := r.db.QueryRowContext(ctx, query, item.Email, item.Username, item.FullName, item.Bio, item.AvatarURL, item.Role, item.Metadata, item.Active, item.CreatedAt, item.UpdatedAt).Scan(&result.ID, &result.Email, &result.Username, &result.FullName, &result.Bio, &result.AvatarURL, &result.Role, &result.Metadata, &result.Active, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AuthorRepositoryImpl) GetByID(ctx context.Context, id int64) (*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE id = $1`
	var item Author
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AuthorRepositoryImpl) Update(ctx context.Context, item *Author) error {
	query := `UPDATE "author" SET email = $1, username = $2, full_name = $3, bio = $4, avatar_url = $5, role = $6, metadata = $7, active = $8, created_at = $9, updated_at = $10 WHERE id = $11`
	_, err := r.db.ExecContext(ctx, query, item.Email, item.Username, item.FullName, item.Bio, item.AvatarURL, item.Role, item.Metadata, item.Active, item.CreatedAt, item.UpdatedAt, item.ID)
	return err
}

func (r *AuthorRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM "author" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *AuthorRepositoryImpl) List(ctx context.Context) ([]*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuthorRepositoryImpl) FindByEmail(ctx context.Context, email string) (*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE email = $1`
	var item Author
	err := r.db.QueryRowContext(ctx, query, email).Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AuthorRepositoryImpl) FindByUsername(ctx context.Context, username string) (*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE username = $1`
	var item Author
	err := r.db.QueryRowContext(ctx, query, username).Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AuthorRepositoryImpl) FindByActive(ctx context.Context, active bool) ([]*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE active = $1`
	rows, err := r.db.QueryContext(ctx, query, active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuthorRepositoryImpl) FindByRole(ctx context.Context, role string) ([]*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE role = $1`
	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuthorRepositoryImpl) FindByEmailAndActive(ctx context.Context, email string, active bool) (*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE email = $1 AND active = $2`
	var item Author
	err := r.db.QueryRowContext(ctx, query, email, active).Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AuthorRepositoryImpl) CountAuthorsByRole(ctx context.Context, role string) (int64, error) {
	query := `SELECT COUNT(*) FROM "author" WHERE role = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, role).Scan(&count)
	return count, err
}

func (r *AuthorRepositoryImpl) ExistsAuthorByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "author" WHERE email = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *AuthorRepositoryImpl) ListAuthorsByActiveOrderByCreatedAtDesc(ctx context.Context, active bool) ([]*Author, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, role, metadata, active, created_at, updated_at FROM "author" WHERE active = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuthorRepositoryImpl) FindActiveWriters(ctx context.Context, limit int) ([]*Author, error) {
	query := `SELECT * FROM author WHERE active = true AND role != 'reader' ORDER BY created_at DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Author
	for rows.Next() {
		var item Author
		if err := rows.Scan(&item.ID, &item.Email, &item.Username, &item.FullName, &item.Bio, &item.AvatarURL, &item.Role, &item.Metadata, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuthorRepositoryImpl) DeactivateInactiveAuthors(ctx context.Context, cutoff int64) error {
	query := `UPDATE author SET active = false WHERE created_at < $1 AND active = true`
	_, err := r.db.ExecContext(ctx, query, cutoff)
	return err
}

