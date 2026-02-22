package blog

import "time"

// Comment is a reader's response to a Post, optionally replying to another Comment.
//
// Demonstrates:
//   - fk with ON DELETE CASCADE (comments die with the post)
//   - fk with ON DELETE SET NULL (keep comments when author account deleted)
//   - self-referential fk (parent_id → comment.id) for threaded discussions
//   - index on a FK column for fast lookup by post
//   - check constraint with function call
//   - default false / default 0
type Comment struct {
	ID        int64     `db:"pk"`
	PostID    int64     `db:"fk:post,id,on_delete:CASCADE"`
	AuthorID  int64     `db:"fk:author,id,on_delete:SET_NULL"`
	ParentID  int64     `db:"fk:comment,id,on_delete:CASCADE"`
	Content   string    `db:"check:length(content) > 0"`
	Approved  bool      `db:"default:false"`
	Likes     int       `db:"default:0,check:likes >= 0"`
	CreatedAt time.Time `db:"default:now()"`
}
