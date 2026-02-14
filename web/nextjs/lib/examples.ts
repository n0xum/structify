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
];
