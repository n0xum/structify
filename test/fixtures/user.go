package fixtures

type User struct {
	ID        int64     `db:"pk"`
	Username  string    `db:"unique"`
	Email     string
	Active    bool
	Created   int64
}

type Product struct {
	ID          int64     `db:"pk"`
	Name        string
	Price       float64
	Description string
	InStock     bool      `db:"-"`
}
