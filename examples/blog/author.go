package blog

import (
	"encoding/json"
	"time"
)

// Author represents a registered writer on the platform.
//
// Demonstrates:
//   - pk: auto-generated primary key
//   - unique: single-column unique constraint
//   - unique_index: unique index (as opposed to UNIQUE constraint)
//   - check: arbitrary SQL expression
//   - enum: PostgreSQL CHECK-based allowed values
//   - default: literal and function defaults
//   - json.RawMessage: maps to JSONB
//   - time.Time: maps to TIMESTAMPTZ
type Author struct {
	ID        int64           `db:"pk"`
	Email     string          `db:"unique,check:email ~* '^[^@]+@[^@]+\\.[^@]+$'"`
	Username  string          `db:"unique_index:idx_author_username,check:length(username) >= 3"`
	FullName  string          `db:"check:length(full_name) > 0"`
	Bio       string
	AvatarURL string
	Role      string          `db:"default:'author',enum:reader,author,editor,admin"`
	Metadata  json.RawMessage `db:"default:'{}'"`
	Active    bool            `db:"default:true"`
	CreatedAt time.Time       `db:"default:now()"`
	UpdatedAt time.Time
}
