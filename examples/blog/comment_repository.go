package blog

import "context"

// CommentRepository defines all database operations for the Comment entity.
//
// Demonstrates:
//   - FindBy across multiple FK columns
//   - SmartQuery for approval workflow queries
//   - CustomSQL for threaded/recursive-style lookups
type CommentRepository interface {
	// --- CRUD ---

	Create(ctx context.Context, comment *Comment) (*Comment, error)
	GetByID(ctx context.Context, id int64) (*Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Comment, error)

	// --- FindBy ---

	// FindByPostID returns all comments on a specific post.
	FindByPostID(ctx context.Context, postID int64) ([]*Comment, error)

	// FindByAuthorID returns all comments made by an author.
	FindByAuthorID(ctx context.Context, authorID int64) ([]*Comment, error)

	// FindByParentID returns direct replies to a comment.
	FindByParentID(ctx context.Context, parentID int64) ([]*Comment, error)

	// FindByPostIDAndApproved returns approved (or pending) comments for a post.
	FindByPostIDAndApproved(ctx context.Context, postID int64, approved bool) ([]*Comment, error)

	// --- SmartQuery ---

	// CountCommentsByPostID emits SELECT COUNT(*) FROM comment WHERE post_id = $1.
	CountCommentsByPostID(ctx context.Context, postID int64) (int64, error)

	// CountCommentsByApproved counts pending or approved comments globally.
	CountCommentsByApproved(ctx context.Context, approved bool) (int64, error)

	// ListCommentsByApprovedOrderByCreatedAtDesc lists comments newest-first.
	ListCommentsByApprovedOrderByCreatedAtDesc(ctx context.Context, approved bool) ([]*Comment, error)

	// --- CustomSQL ---

	//sql:"SELECT * FROM comment WHERE post_id = $1 AND approved = true ORDER BY likes DESC LIMIT $2"
	FindTopLikedByPost(ctx context.Context, postID int64, limit int) ([]*Comment, error)

	//sql:"UPDATE comment SET approved = true WHERE post_id = $1 AND approved = false"
	ApproveAllByPost(ctx context.Context, postID int64) error
}
