export type Example = {
  label: string;
  description: string;
  source: string;
};

export const EXAMPLES: Example[] = [
  {
    label: "User",
    description: "pk, unique, basic types",
    source: `package models

type User struct {
\tID       int64  \`db:"pk"\`
\tUsername string \`db:"unique"\`
\tEmail    string
\tActive   bool
\tCreated  int64
}`,
  },
  {
    label: "Product",
    description: "pk, ignored field",
    source: `package models

type Product struct {
\tID          int64   \`db:"pk"\`
\tName        string
\tPrice       float64
\tDescription string
\tInStock     bool    \`db:"-"\`
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
    description: "index, unique_index, composite index",
    source: `package models

type Article struct {
\tID       int64  \`db:"pk"\`
\tSlug     string \`db:"unique_index:uq_slug"\`
\tTitle    string \`db:"index:idx_title_lang"\`
\tLanguage string \`db:"index:idx_title_lang"\`
\tViews    int64  \`db:"index"\`
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
