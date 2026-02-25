package blog

import (
	"context"
	"database/sql"
)

type PostRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, item *Post) (*Post, error) {
	query := `INSERT INTO "post" (author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at`
	var result Post
	err := r.db.QueryRowContext(ctx, query, item.AuthorID, item.CategoryID, item.Title, item.Slug, item.Content, item.Summary, item.Status, item.ViewCount, item.Featured, item.PublishedAt, item.CreatedAt, item.UpdatedAt).Scan(&result.ID, &result.AuthorID, &result.CategoryID, &result.Title, &result.Slug, &result.Content, &result.Summary, &result.Status, &result.ViewCount, &result.Featured, &result.PublishedAt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *PostRepositoryImpl) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE id = $1`
	var item Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PostRepositoryImpl) Update(ctx context.Context, item *Post) error {
	query := `UPDATE "post" SET author_id = $1, category_id = $2, title = $3, slug = $4, content = $5, summary = $6, status = $7, view_count = $8, featured = $9, published_at = $10, created_at = $11, updated_at = $12 WHERE id = $13`
	_, err := r.db.ExecContext(ctx, query, item.AuthorID, item.CategoryID, item.Title, item.Slug, item.Content, item.Summary, item.Status, item.ViewCount, item.Featured, item.PublishedAt, item.CreatedAt, item.UpdatedAt, item.ID)
	return err
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM "post" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostRepositoryImpl) List(ctx context.Context) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE slug = $1`
	var item Post
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PostRepositoryImpl) FindByAuthorID(ctx context.Context, authorID int64) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE author_id = $1`
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) FindByStatus(ctx context.Context, status string) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID int64) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) FindByAuthorIDAndStatus(ctx context.Context, authorID int64, status string) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE author_id = $1 AND status = $2`
	rows, err := r.db.QueryContext(ctx, query, authorID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) FindByFeaturedAndStatus(ctx context.Context, featured bool, status string) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE featured = $1 AND status = $2`
	rows, err := r.db.QueryContext(ctx, query, featured, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) ListPostsByStatus(ctx context.Context, status string) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) ListPostsByFeatured(ctx context.Context, featured bool) ([]*Post, error) {
	query := `SELECT id, author_id, category_id, title, slug, content, summary, status, view_count, featured, published_at, created_at, updated_at FROM "post" WHERE featured = $1`
	rows, err := r.db.QueryContext(ctx, query, featured)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) CountPostsByAuthorID(ctx context.Context, authorID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM "post" WHERE author_id = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	return count, err
}

func (r *PostRepositoryImpl) CountPostsByStatus(ctx context.Context, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM "post" WHERE status = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}

func (r *PostRepositoryImpl) ExistsPostBySlug(ctx context.Context, slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "post" WHERE slug = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	return exists, err
}

func (r *PostRepositoryImpl) FindLatestPublished(ctx context.Context, limit int) ([]*Post, error) {
	query := `SELECT * FROM post WHERE status = 'published' ORDER BY created_at DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) FindPopularByCategory(ctx context.Context, categoryID int64) ([]*Post, error) {
	query := `SELECT * FROM post WHERE category_id = $1 AND status = 'published' ORDER BY view_count DESC`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Post
	for rows.Next() {
		var item Post
		if err := rows.Scan(&item.ID, &item.AuthorID, &item.CategoryID, &item.Title, &item.Slug, &item.Content, &item.Summary, &item.Status, &item.ViewCount, &item.Featured, &item.PublishedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostRepositoryImpl) IncrementViewCount(ctx context.Context, id int64) error {
	query := `UPDATE post SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

