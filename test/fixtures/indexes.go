package fixtures

type IndexedUser struct {
	ID       int64  `db:"pk"`
	Email    string `db:"unique_index:idx_email"`
	Username string `db:"index"`
	Active   bool   `db:"default:true"`
}

type IndexedProduct struct {
	ID        int64   `db:"pk"`
	Name      string  `db:"index:idx_name_category"`
	Category  string  `db:"index:idx_name_category"`
	Price     float64 `db:"check:price >= 0"`
	Available bool    `db:"default:true"`
}

type OrderStatus struct {
	ID        int64      `db:"pk"`
	Status    string     `db:"enum:pending,processing,shipped,delivered,cancelled"`
	Priority  string     `db:"enum:low,medium,high"`
	CreatedAt int64      `db:"default:extract(epoch from now())"`
}
