package fixtures

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByActive(ctx context.Context, active bool) ([]*User, error)

	//sql:"SELECT * FROM users WHERE created > $1 ORDER BY username"
	FindRecentUsers(ctx context.Context, since int64) ([]*User, error)
}
