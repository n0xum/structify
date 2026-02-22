export type Example = {
  label: string;
  description: string;
  source: string;
};

export const EXAMPLES: Example[] = [
  {
    label: "User",
    description: "pk, unique, basic repo",
    source: `package models

import "context"

type User struct {
	ID       int64  \`db:"pk"\`
	Username string \`db:"unique"\`
	Email    string
	Active   bool
	Created  int64
}

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}`,
  },
  {
    label: "Product",
    description: "pk, ignored field, basic repo",
    source: `package models

import "context"

type Product struct {
	ID          int64   \`db:"pk"\`
	Name        string
	Price       float64
	Description string
	InStock     bool    \`db:"-"\`
}

type ProductRepository interface {
	FindByID(ctx context.Context, id int64) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Create(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
}`,
  },
  {
    label: "OrderItem",
    description: "multiple fields, float types",
    source: `package models

type OrderItem struct {
\tID        int64   \`db:"pk"\`
\tOrderID   int64
\tProductID int64
\tQuantity  int32
\tUnitPrice float64
\tDiscount  float32
}`,
  },
  {
    label: "Constraints",
    description: "check, default, enum",
    source: `package models

type Person struct {
\tID        int64   \`db:"pk"\`
\tName      string  \`db:"check:length(name) > 0"\`
\tAge       int     \`db:"check:age >= 18,default:18"\`
\tStatus    string  \`db:"enum:active,inactive,banned"\`
\tActive    bool    \`db:"default:true"\`
\tCreatedAt int64   \`db:"default:now()"\`
}`,
  },
  {
    label: "Indexes",
    description: "index, unique_index, composite index, search repo",
    source: `package models

import "context"

type Article struct {
	ID       int64  \`db:"pk"\`
	Slug     string \`db:"unique_index:uq_slug"\`
	Title    string \`db:"index:idx_title_lang"\`
	Language string \`db:"index:idx_title_lang"\`
	Views    int64  \`db:"index"\`
}

type ArticleRepository interface {
	FindBySlug(ctx context.Context, slug string) (*Article, error)
	FindByTitleAndLanguage(ctx context.Context, title string, language string) ([]*Article, error)
	FindByViewsGreaterThan(ctx context.Context, views int64) ([]*Article, error)
}`,
  },
  {
    label: "Foreign Keys",
    description: "fk, on_delete:CASCADE",
    source: `package models

type User struct {
\tID    int64  \`db:"pk"\`
\tEmail string \`db:"unique"\`
}

type Post struct {
\tID      int64  \`db:"pk"\`
\tUserID  int64  \`db:"fk:users,id,on_delete:CASCADE"\`
\tTitle   string
\tContent string
}`,
  },
  {
    label: "Composite PK & FK",
    description: "composite primary key, composite foreign key",
    source: `package models

type OrderItem struct {
\tOrderID   int64 \`db:"pk"\`
\tProductID int64 \`db:"pk"\`
\tQuantity  int
\tPrice     float64
}

type OrderItemNote struct {
\tNoteID    int64  \`db:"pk"\`
\tOrderID   int64  \`db:"fk:fk_oi,order_items,order_id,on_delete:CASCADE"\`
\tProductID int64  \`db:"fk:fk_oi,order_items,product_id"\`
\tNote      string
}`,
  },
];
