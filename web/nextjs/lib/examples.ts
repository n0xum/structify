export type Example = {
  label: string;
  description: string;
  source: string;
};

/** Examples for SQL schema generation (struct-only, no interface required). */
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

/** Examples for repository code generation (struct + interface). */
export const REPO_EXAMPLES: Example[] = [
  {
    label: "User — CRUD + FindBy",
    description: "Create, GetByID, Update, Delete, List, FindByEmail",
    source: `package models

import "context"

type User struct {
\tID       int64  \`db:"pk"\`
\tUsername string \`db:"unique"\`
\tEmail    string \`db:"unique"\`
\tActive   bool   \`db:"default:true"\`
\tCreatedAt int64
}

type UserRepository interface {
\tCreate(ctx context.Context, u *User) (*User, error)
\tGetByID(ctx context.Context, id int64) (*User, error)
\tUpdate(ctx context.Context, u *User) error
\tDelete(ctx context.Context, id int64) error
\tList(ctx context.Context) ([]*User, error)
\tFindByEmail(ctx context.Context, email string) (*User, error)
\tFindByUsername(ctx context.Context, username string) (*User, error)
\tFindByActiveAndUsername(ctx context.Context, active bool, username string) (*User, error)
}`,
  },
  {
    label: "Article — SmartQuery",
    description: "Count, Exists, ListByField, OrderBy",
    source: `package models

import "context"

type Article struct {
\tID          int64  \`db:"pk"\`
\tSlug        string \`db:"unique"\`
\tTitle       string
\tAuthorID    int64
\tPublished   bool   \`db:"default:false"\`
\tViewCount   int64  \`db:"default:0"\`
\tCreatedAt   int64  \`db:"default:now()"\`
}

type ArticleRepository interface {
\tCreate(ctx context.Context, a *Article) (*Article, error)
\tGetByID(ctx context.Context, id int64) (*Article, error)
\tList(ctx context.Context) ([]*Article, error)

\t// SmartQuery: method name encodes the full query
\tCountArticlesByPublished(ctx context.Context, published bool) (int64, error)
\tExistsArticleBySlug(ctx context.Context, slug string) (bool, error)
\tListArticlesByAuthorID(ctx context.Context, authorID int64) ([]*Article, error)
\tListArticlesByPublishedOrderByCreatedAtDesc(ctx context.Context, published bool) ([]*Article, error)
\tListArticlesByViewCountGreaterThan(ctx context.Context, viewCount int64) ([]*Article, error)
}`,
  },
  {
    label: "Order — CustomSQL",
    description: "//sql: comments for arbitrary queries",
    source: `package models

import "context"

type Order struct {
\tID          int64   \`db:"pk"\`
\tCustomerID  int64
\tStatus      string  \`db:"default:'pending',enum:pending,confirmed,shipped,delivered,cancelled"\`
\tTotalAmount float64
\tCreatedAt   int64   \`db:"default:now()"\`
}

type OrderRepository interface {
\tCreate(ctx context.Context, o *Order) (*Order, error)
\tGetByID(ctx context.Context, id int64) (*Order, error)
\tList(ctx context.Context) ([]*Order, error)
\tFindByCustomerID(ctx context.Context, customerID int64) ([]*Order, error)

\t// SmartQuery
\tCountOrdersByStatus(ctx context.Context, status string) (int64, error)
\tListOrdersByCustomerIDOrderByCreatedAtDesc(ctx context.Context, customerID int64) ([]*Order, error)

\t// CustomSQL: provide exact SQL via //sql: comment
\t//sql:"SELECT * FROM \\"order\\" WHERE customer_id = $1 AND status = 'delivered' ORDER BY created_at DESC LIMIT $2"
\tFindRecentDelivered(ctx context.Context, customerID int64, limit int) ([]*Order, error)

\t//sql:"SELECT COALESCE(SUM(total_amount), 0.0) FROM \\"order\\" WHERE customer_id = $1"
\tGetTotalSpent(ctx context.Context, customerID int64) (float64, error)
}`,
  },
  {
    label: "Product — FindBy operators",
    description: "GreaterThan, Like, IsNull, In operators",
    source: `package models

import "context"

type Product struct {
\tID          int64   \`db:"pk"\`
\tSKU         string  \`db:"unique"\`
\tName        string
\tPrice       float64
\tStockQty    int32   \`db:"default:0"\`
\tDeletedAt   int64
}

type ProductRepository interface {
\tCreate(ctx context.Context, p *Product) (*Product, error)
\tGetByID(ctx context.Context, id int64) (*Product, error)
\tList(ctx context.Context) ([]*Product, error)
\tFindBySKU(ctx context.Context, sku string) (*Product, error)

\t// Operator suffixes in method names
\tListProductsByPriceGreaterThan(ctx context.Context, price float64) ([]*Product, error)
\tListProductsByNameLike(ctx context.Context, name string) ([]*Product, error)
\tListProductsByDeletedAtIsNull(ctx context.Context) ([]*Product, error)
\tCountProductsByStockQtyGreaterThan(ctx context.Context, stockQty int32) (int64, error)
}`,
  },
];
