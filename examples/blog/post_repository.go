package blog

import "context"

// PostRepository defines all database operations for the Post entity.
//
// Demonstrates:
//   - GetByID and FindBy returning both single item and slice
//   - FindBy across a FK column (FindByAuthorID)
//   - FindBy with AND of two fields
//   - SmartQuery patterns: List…By…, Count…By…, Exists…By…
//   - CustomSQL with $1 / $2 placeholders
type PostRepository interface {
	// --- CRUD ---

	Create(ctx context.Context, post *Post) (*Post, error)
	GetByID(ctx context.Context, id int64) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Post, error)

	// --- FindBy ---

	// FindBySlug returns the single post with the given URL slug.
	FindBySlug(ctx context.Context, slug string) (*Post, error)

	// FindByAuthorID returns all posts written by a specific author.
	FindByAuthorID(ctx context.Context, authorID int64) ([]*Post, error)

	// FindByStatus returns all posts in the given status (draft/published/…).
	FindByStatus(ctx context.Context, status string) ([]*Post, error)

	// FindByCategoryID returns all posts in a category.
	FindByCategoryID(ctx context.Context, categoryID int64) ([]*Post, error)

	// FindByAuthorIDAndStatus combines author and status filters.
	FindByAuthorIDAndStatus(ctx context.Context, authorID int64, status string) ([]*Post, error)

	// FindByFeaturedAndStatus returns featured posts with a given status.
	FindByFeaturedAndStatus(ctx context.Context, featured bool, status string) ([]*Post, error)

	// --- SmartQuery ---

	// ListPostsByStatus emits SELECT … FROM post WHERE status = $1.
	ListPostsByStatus(ctx context.Context, status string) ([]*Post, error)

	// ListPostsByFeatured emits SELECT … FROM post WHERE featured = $1.
	ListPostsByFeatured(ctx context.Context, featured bool) ([]*Post, error)

	// CountPostsByAuthorID emits SELECT COUNT(*) FROM post WHERE author_id = $1.
	CountPostsByAuthorID(ctx context.Context, authorID int64) (int64, error)

	// CountPostsByStatus emits SELECT COUNT(*) FROM post WHERE status = $1.
	CountPostsByStatus(ctx context.Context, status string) (int64, error)

	// ExistsPostBySlug emits SELECT EXISTS(SELECT 1 FROM post WHERE slug = $1).
	ExistsPostBySlug(ctx context.Context, slug string) (bool, error)

	// --- CustomSQL ---

	//sql:"SELECT * FROM post WHERE status = 'published' ORDER BY created_at DESC LIMIT $1"
	FindLatestPublished(ctx context.Context, limit int) ([]*Post, error)

	//sql:"SELECT * FROM post WHERE category_id = $1 AND status = 'published' ORDER BY view_count DESC"
	FindPopularByCategory(ctx context.Context, categoryID int64) ([]*Post, error)

	//sql:"UPDATE post SET view_count = view_count + 1 WHERE id = $1"
	IncrementViewCount(ctx context.Context, id int64) error
}
