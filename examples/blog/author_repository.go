package blog

import "context"

// AuthorRepository defines all database operations for the Author entity.
//
// Demonstrates:
//   - full CRUD (Create / GetByID / Update / Delete / List)
//   - FindBy single field (generates parameterized SELECT … WHERE)
//   - FindBy multiple fields joined with AND
//   - SmartQuery: method name encodes the whole query (CountXByY, ExistsXByY)
//   - CustomSQL: arbitrary SQL supplied via //sql: comment
type AuthorRepository interface {
	// --- CRUD ---

	Create(ctx context.Context, author *Author) (*Author, error)
	GetByID(ctx context.Context, id int64) (*Author, error)
	Update(ctx context.Context, author *Author) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Author, error)

	// --- FindBy (field-based WHERE clauses) ---

	// FindByEmail returns the author with the given e-mail address.
	FindByEmail(ctx context.Context, email string) (*Author, error)

	// FindByUsername looks up an author by their unique username.
	FindByUsername(ctx context.Context, username string) (*Author, error)

	// FindByActive returns all authors matching the given active flag.
	FindByActive(ctx context.Context, active bool) ([]*Author, error)

	// FindByRole returns all authors that hold the given role.
	FindByRole(ctx context.Context, role string) ([]*Author, error)

	// FindByEmailAndActive filters by both email and active in one query.
	FindByEmailAndActive(ctx context.Context, email string, active bool) (*Author, error)

	// --- SmartQuery (pattern-matched SQL generation) ---

	// CountAuthorsByRole uses the Count…By… pattern to emit SELECT COUNT(*).
	CountAuthorsByRole(ctx context.Context, role string) (int64, error)

	// ExistsAuthorByEmail uses the Exists…By… pattern to emit SELECT EXISTS(…).
	ExistsAuthorByEmail(ctx context.Context, email string) (bool, error)

	// ListAuthorsByActiveOrderByCreatedAtDesc orders active authors by sign-up date.
	ListAuthorsByActiveOrderByCreatedAtDesc(ctx context.Context, active bool) ([]*Author, error)

	// --- CustomSQL (exact SQL in //sql: comment) ---

	//sql:"SELECT * FROM author WHERE active = true AND role != 'reader' ORDER BY created_at DESC LIMIT $1"
	FindActiveWriters(ctx context.Context, limit int) ([]*Author, error)

	//sql:"UPDATE author SET active = false WHERE created_at < $1 AND active = true"
	DeactivateInactiveAuthors(ctx context.Context, cutoff int64) error
}
