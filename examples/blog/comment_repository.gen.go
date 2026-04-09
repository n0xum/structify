package blog

import (
	"context"
	"database/sql"
)

type CommentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, item *Comment) (*Comment, error) {
	query := `INSERT INTO "comment" (post_id, author_id, parent_id, content, approved, likes, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, post_id, author_id, parent_id, content, approved, likes, created_at`
	var result Comment
	err := r.db.QueryRowContext(ctx, query, item.PostID, item.AuthorID, item.ParentID, item.Content, item.Approved, item.Likes, item.CreatedAt).Scan(&result.ID, &result.PostID, &result.AuthorID, &result.ParentID, &result.Content, &result.Approved, &result.Likes, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *CommentRepositoryImpl) GetByID(ctx context.Context, id int64) (*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" WHERE id = $1`
	var item Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, item *Comment) error {
	query := `UPDATE "comment" SET post_id = $1, author_id = $2, parent_id = $3, content = $4, approved = $5, likes = $6, created_at = $7 WHERE id = $8`
	_, err := r.db.ExecContext(ctx, query, item.PostID, item.AuthorID, item.ParentID, item.Content, item.Approved, item.Likes, item.CreatedAt, item.ID)
	return err
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM "comment" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CommentRepositoryImpl) List(ctx context.Context) ([]*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) FindByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" WHERE post_id = $1`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) FindByAuthorID(ctx context.Context, authorID int64) ([]*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" WHERE author_id = $1`
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) FindByParentID(ctx context.Context, parentID int64) ([]*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" WHERE parent_id = $1`
	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) FindByPostIDAndApproved(ctx context.Context, postID int64, approved bool) ([]*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" WHERE post_id = $1 AND approved = $2`
	rows, err := r.db.QueryContext(ctx, query, postID, approved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) CountCommentsByPostID(ctx context.Context, postID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM "comment" WHERE post_id = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	return count, err
}

func (r *CommentRepositoryImpl) CountCommentsByApproved(ctx context.Context, approved bool) (int64, error) {
	query := `SELECT COUNT(*) FROM "comment" WHERE approved = $1`
	var count int64
	err := r.db.QueryRowContext(ctx, query, approved).Scan(&count)
	return count, err
}

func (r *CommentRepositoryImpl) ListCommentsByApprovedOrderByCreatedAtDesc(ctx context.Context, approved bool) ([]*Comment, error) {
	query := `SELECT id, post_id, author_id, parent_id, content, approved, likes, created_at FROM "comment" WHERE approved = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, approved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) FindTopLikedByPost(ctx context.Context, postID int64, limit int) ([]*Comment, error) {
	query := `SELECT * FROM comment WHERE post_id = $1 AND approved = true ORDER BY likes DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, postID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Comment
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.PostID, &item.AuthorID, &item.ParentID, &item.Content, &item.Approved, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CommentRepositoryImpl) ApproveAllByPost(ctx context.Context, postID int64) error {
	query := `UPDATE comment SET approved = true WHERE post_id = $1 AND approved = false`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

