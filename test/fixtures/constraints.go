package fixtures

type Person struct {
	ID       int64  `db:"pk"`
	Name     string `db:"check:length(name) > 0"`
	Age      int    `db:"check:age >= 18"`
	Email    string `db:"check:email ~* '^[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,}$'"`
	Active   bool   `db:"default:true"`
	Role     string `db:"default:'user'"`
	Created  int64  `db:"default:extract(epoch from now())"`
}
