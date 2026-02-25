package blog

import "time"

// Post is an article written by an Author, belonging to a Category.
//
// Demonstrates:
//   - fk with ON DELETE CASCADE (author must exist)
//   - fk with ON DELETE SET NULL (category is optional)
//   - composite index across two columns
//   - enum with default
//   - integer column with both check and default
//   - bool column with default
//   - nullable time.Time (no default → nullable in SQL)
type Post struct {
	ID          int64     `db:"pk"`
	AuthorID    int64     `db:"fk:author,id,on_delete:CASCADE"`
	CategoryID  int64     `db:"fk:category,id,on_delete:SET_NULL"`
	Title       string    `db:"check:length(title) >= 5"`
	Slug        string    `db:"unique,check:length(slug) > 0"`
	Content     string    `db:"check:length(content) > 0"`
	Summary     string
	Status      string    `db:"default:'draft',enum:draft,review,published,archived"`
	ViewCount   int       `db:"default:0,check:view_count >= 0"`
	Featured    bool      `db:"default:false"`
	PublishedAt time.Time
	CreatedAt   time.Time `db:"default:now()"`
	UpdatedAt   time.Time
}

// Category organises posts into topics.
//
// Demonstrates:
//   - unique on multiple independent columns
//   - unique_index for index-backed uniqueness
//   - default literal and integer default
type Category struct {
	ID          int64  `db:"pk"`
	Name        string `db:"unique,check:length(name) > 0"`
	Slug        string `db:"unique_index:idx_category_slug,check:length(slug) > 0"`
	Description string
	Color       string `db:"default:'#6b7280'"`
	SortOrder   int    `db:"default:0"`
}

// Tag is a lightweight label that can be attached to many posts.
//
// Demonstrates:
//   - two independent unique columns on the same struct
type Tag struct {
	ID   int64  `db:"pk"`
	Name string `db:"unique,check:length(name) > 0"`
	Slug string `db:"unique"`
}

// PostTag is the many-to-many join table between Post and Tag.
//
// Demonstrates:
//   - composite primary key (post_id + tag_id)
//   - each PK field is also a FK with ON DELETE CASCADE
//   - pk and fk combined in one db tag
type PostTag struct {
	PostID    int64 `db:"pk,fk:post,id,on_delete:CASCADE"`
	TagID     int64 `db:"pk,fk:tag,id,on_delete:CASCADE"`
	SortOrder int   `db:"default:0,check:sort_order >= 0"`
}
